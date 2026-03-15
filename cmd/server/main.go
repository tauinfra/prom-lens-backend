package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
	"valyria-backend/internal/apps/kubernetes/worker"
	"valyria-backend/internal/core/config"
	"valyria-backend/internal/core/database"
	"valyria-backend/internal/core/initialize"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/di"
	"valyria-backend/internal/routes"
)

var (
	c         string
	version   = "dev"     // 构建时通过 -ldflags 注入
	buildTime = "unknown" // 构建时通过 -ldflags 注入
)

func init() {
	flag.StringVar(&c, "c", "config/config.yaml", "set configuration `file`")
}

// printStartupInfo 打印启动信息
func printStartupInfo(cfg *config.Config, startTime time.Time) {
	logger.Info("========================================")
	logger.Info("  Valyria Backend Server")
	logger.Info("========================================")
	logger.Infof("Version:     %s", version)
	logger.Infof("Build Time:  %s", buildTime)
	logger.Infof("Go Version:  %s", runtime.Version())
	logger.Infof("OS/Arch:     %s/%s", runtime.GOOS, runtime.GOARCH)
	logger.Info("----------------------------------------")
	logger.Info("Configuration:")
	logger.Infof("  Server Port:     %d", cfg.Server.Port)
	logger.Infof("  Server Mode:     %s", cfg.Server.Mode)
	logger.Infof("  Database:        %s@%s:%d/%s", cfg.Database.Username, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName)
	logger.Infof("  Redis:           %s (DB: %d)", cfg.Redis.Addr, cfg.Redis.DB)
	logger.Infof("  Log Level:       %s", cfg.Log.LogLevel)
	logger.Infof("  Auth Issuer:     %s", cfg.Auth.Issuer)
	logger.Infof("  Auth Audience:   %s", cfg.Auth.Audience)
	logger.Infof("  Token Expires:   %v / %v", cfg.Auth.AccessTokenExpires, cfg.Auth.RefreshTokenExpires)
	logger.Info("----------------------------------------")
	logger.Infof("Server URL:  http://localhost:%d", cfg.Server.Port)
	logger.Infof("Start Time:  %s", startTime.Format("2006-01-02 15:04:05"))
	logger.Info("========================================")
}

func main() {
	startTime := time.Now()
	flag.Parse() // 调用 flag

	// 初始化配置
	cfg, err := config.InitConfig(c)
	if err != nil {
		log.Fatalf("Failed to init config: %v", err)
	}

	// 初始化组件
	db, encryptor, err := initialize.Components(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize components: %v", err)
	}
	
	// 打印启动信息
	printStartupInfo(cfg, startTime)

	// 初始化路由
	provider := di.NewProvider(cfg, db, encryptor)
	r := routes.SetupRouter(db, provider)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: r,
	}

	// 在 goroutine 中启动服务器
	go func() {
		logger.Infof("Server listening on :%d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("failed to start server: %v", err)
		}
	}()

	// 后台：K8s 权限同步 Worker（scan_interval / batch_size / max_retry / base_delay / max_delay 来自 config.permission_worker）
	runCtx, runCancel := context.WithCancel(context.Background())
	defer runCancel()
	permWorkerCfg := cfg.PermissionWorker.Parse()
	permSyncWorker := worker.NewPermissionSyncWorker(provider.Kubernetes.Permission.Service, permWorkerCfg.ScanInterval)
	go permSyncWorker.Run(runCtx)
	hpaHistoryWorker := worker.NewHpaHistorySyncWorker(provider.Kubernetes.HpaHistorySvc, 5*time.Minute)
	go hpaHistoryWorker.Run(runCtx)

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down server...")
	runCancel()

	// 优雅关闭超时：30秒
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 关闭 HTTP 服务器
	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
	}

	// 清理资源
	logger.Info("Cleaning up resources...")
	if err := database.CloseDatabase(db); err != nil {
		logger.Errorf("Failed to close database: %v", err)
	}
	if err := database.CloseRedis(); err != nil {
		logger.Errorf("Failed to close redis: %v", err)
	}
	if err := logger.Sync(); err != nil {
		logger.Errorf("Failed to sync logger: %v", err)
	}

	logger.Info("Server exited gracefully")
}
