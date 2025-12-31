package main

import (
	"flag"
	"log"
	"valyria-backend/internal/core/configs"
	"valyria-backend/internal/core/initialize"
	"valyria-backend/internal/core/logger"
	"valyria-backend/internal/pkg/di"
	"valyria-backend/internal/routes"
)

var c string

func init() {
	flag.StringVar(&c, "c", "configs/config.yaml", "set configuration `file`")
}

func main() {
	flag.Parse() // 调用 flag
	// 初始化配置
	if err := configs.InitConfig(c); err != nil {
		log.Fatalf("Failed to init configs: %v", err)
	}
	cfg := configs.GetConfig()

	// 初始化组件
	db, encryptor, err := initialize.Components(cfg)
	if err != nil {
		logger.Fatalf("Failed to initialize components, err: %v", err)
	}
	// 初始化路由
	provider := di.NewProvider(cfg, db, encryptor)
	r := routes.SetupRouter(db, provider)

	// 启动服务
	if err := r.Run("127.0.0.1:8080"); err != nil {
		logger.Fatalf("failed to run server: %v\n", err)
	}
}
