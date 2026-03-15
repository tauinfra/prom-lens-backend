package config

import (
	"fmt"
	"strings"
	"time"
	"valyria-backend/internal/core/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// ConfigChange 表示配置变更
type ConfigChange struct {
	Key      string
	OldValue interface{}
	NewValue interface{}
}

type Config struct {
	Server           ServerConfig           `mapstructure:"server"`
	Database         DatabaseConfig         `mapstructure:"database"`
	Redis            RedisConfig            `mapstructure:"redis"`
	Auth             AuthConfig             `mapstructure:"auth"`
	Log              LogConfig              `mapstructure:"logging"`
	K8s              K8sConfig              `mapstructure:"kubernetes"`
	App              AppConfig              `mapstructure:"app"`
	Prometheus       PrometheusConfig       `mapstructure:"prometheus"`
	Dragon           DragonConfig           `mapstructure:"dragon"`
	PermissionWorker PermissionWorkerConfig `mapstructure:"permission_worker"`
}

// PermissionWorkerConfig K8s 权限同步 Worker 配置（config.yaml permission_worker）
// 所有字段均可选，未配置或无效时 Parse() 使用默认值
type PermissionWorkerConfig struct {
	ScanInterval string `mapstructure:"scan_interval"` // 扫描间隔，如 "10s", "2m"
	BatchSize    int    `mapstructure:"batch_size"`   // 每批拉取条数，建议 10～1000
	WorkerCount  int    `mapstructure:"worker_count"` // 预留：未来并行处理协程数，当前未使用
	MaxRetry     int    `mapstructure:"max_retry"`    // 最大重试次数，达此后置为 dead
	BaseDelay    string `mapstructure:"base_delay"`   // 退避基准，如 "10s"；退避公式 base_delay * 2^retry
	MaxDelay     string `mapstructure:"max_delay"`    // 退避上限，如 "15m"
}

// ParsedPermissionWorker 解析后的 Worker 配置（duration 已解析，零值已填默认）
type ParsedPermissionWorker struct {
	ScanInterval time.Duration
	BatchSize    int
	WorkerCount  int
	MaxRetry     int
	BaseDelay    time.Duration
	MaxDelay     time.Duration
}

const (
	permissionWorkerMaxBatchSize = 1000 // batch_size 上限，防止误配过大
)

// Parse 解析并返回可用配置，无效或零值使用默认；batch_size 限制在 1～permissionWorkerMaxBatchSize
func (c *PermissionWorkerConfig) Parse() ParsedPermissionWorker {
	const (
		defaultScanInterval = 2 * time.Minute
		defaultBatchSize   = 100
		defaultWorkerCount = 1
		defaultMaxRetry    = 5
		defaultBaseDelay   = 1 * time.Minute
		defaultMaxDelay    = 15 * time.Minute
	)
	out := ParsedPermissionWorker{
		BatchSize:    defaultBatchSize,
		WorkerCount:  defaultWorkerCount,
		MaxRetry:     defaultMaxRetry,
		ScanInterval: defaultScanInterval,
		BaseDelay:    defaultBaseDelay,
		MaxDelay:     defaultMaxDelay,
	}
	if c.BatchSize > 0 {
		out.BatchSize = c.BatchSize
		if out.BatchSize > permissionWorkerMaxBatchSize {
			out.BatchSize = permissionWorkerMaxBatchSize
		}
	}
	if c.WorkerCount > 0 {
		out.WorkerCount = c.WorkerCount
	}
	if c.MaxRetry > 0 {
		out.MaxRetry = c.MaxRetry
	}
	if c.ScanInterval != "" {
		if d, err := time.ParseDuration(c.ScanInterval); err == nil && d > 0 {
			out.ScanInterval = d
		}
	}
	if c.BaseDelay != "" {
		if d, err := time.ParseDuration(c.BaseDelay); err == nil && d >= 0 {
			out.BaseDelay = d
		}
	}
	if c.MaxDelay != "" {
		if d, err := time.ParseDuration(c.MaxDelay); err == nil && d > 0 {
			out.MaxDelay = d
		}
	}
	if out.MaxDelay > 0 && out.BaseDelay > out.MaxDelay {
		out.MaxDelay = out.BaseDelay // 保证 max_delay >= base_delay
	}
	return out
}

type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

type DatabaseConfig struct {
	Username        string        `mapstructure:"username"`
	Password        string        `mapstructure:"password"`
	Host            string        `mapstructure:"host"`
	DBName          string        `mapstructure:"dbname"`
	Port            int           `mapstructure:"port"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
}

type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

type AuthConfig struct {
	AccessTokenExpires  time.Duration `mapstructure:"access_token_expires"`
	RefreshTokenExpires time.Duration `mapstructure:"refresh_token_expires"`
	Whitelist           []string      `mapstructure:"whitelist"`
	Issuer              string        `mapstructure:"issuer"`
	Audience            string        `mapstructure:"audience"`
	LoginMaxRetries     int           `mapstructure:"login_max_retries"`
	LoginLockWindow     time.Duration `mapstructure:"login_lock_window"`
}

type LogConfig struct {
	LogDir        string `mapstructure:"log_dir"`
	InfoLogName   string `mapstructure:"info_log_name"`
	ErrorLogName  string `mapstructure:"error_log_name"`
	MaxSize       int    `mapstructure:"max_size"`
	MaxBackups    int    `mapstructure:"max_backups"`
	MaxAge        int    `mapstructure:"max_age"`
	Compress      bool   `mapstructure:"compress"`
	LogLevel      string `mapstructure:"log_level"`
	EnableConsole bool   `mapstructure:"enable_console"`
}

type AppConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpire     string `mapstructure:"jwt_expire"`
	EncryptionKey string `mapstructure:"encryption_key"`
}

type K8sConfig struct {
	Timeout int     `mapstructure:"timeout"`
	QPS     float32 `mapstructure:"qps"`
	Burst   int     `mapstructure:"burst"`
	ArgoCD  ArgoCDConfig `mapstructure:"argocd"`
	Tekton  TektonConfig `mapstructure:"tekton"`
	Actions []string `mapstructure:"actions"`
}

type DragonConfig struct {
	Kustomize KustomizeConfig `mapstructure:"kustomize"`
}

type ArgoCDConfig struct {
	Server   string `mapstructure:"server"`
	Insecure bool   `mapstructure:"insecure"`
}

type KustomizeConfig struct {
	RepoURL string `mapstructure:"repo_url"`
	Branch  string `mapstructure:"branch"`
}

type TektonConfig struct {
	Namespace string `mapstructure:"namespace"`
}

type PrometheusConfig struct {
	Target PrometheusTargetConfig `mapstructure:"target"`
	Rule   PrometheusRuleConfig   `mapstructure:"rule"`
}

type PrometheusTargetConfig struct {
	Namespace string `mapstructure:"namespace"`
	ConfigMap string `mapstructure:"configMap"`
}

type PrometheusRuleConfig struct {
	Namespace string `mapstructure:"namespace"`
	ConfigMap string `mapstructure:"configmap"`
}

var C *Config // Deprecated: 使用 InitConfig 返回的配置，全局变量仅用于向后兼容

// InitConfig 初始化配置
// 支持通过环境变量覆盖配置值，环境变量命名规则：
//   - 前缀：VALYRIA_
//   - 嵌套结构使用下划线分隔，例如：
//     VALYRIA_SERVER_PORT -> server.port
//     VALYRIA_DATABASE_HOST -> database.host
//     VALYRIA_AUTH_ACCESS_TOKEN_EXPIRES -> auth.access_token_expires
//   - 数组类型暂不支持通过环境变量设置
//
// 示例：
//
//	export VALYRIA_SERVER_PORT=9090
//	export VALYRIA_DATABASE_HOST=192.168.1.1
//	export VALYRIA_DATABASE_PASSWORD=mypassword
func InitConfig(cfgFile string) (*Config, error) {
	viper.SetConfigFile(cfgFile)

	// 设置环境变量前缀
	viper.SetEnvPrefix("VALYRIA_")
	// 设置环境变量键名替换器：将点号替换为下划线，支持嵌套结构
	// 例如：VALYRIA_SERVER_PORT -> server.port
	//      VALYRIA_DATABASE_HOST -> database.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// 自动读取环境变量
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal read the configuration file failed, %w", err)
	}

	cfg := &Config{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, fmt.Errorf("unmarshal configuration file config failed, %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// 设置全局变量（向后兼容）
	C = cfg

	// 监听配置文件变更
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.Infof("Config file changed: %s (op: %s)", e.Name, e.Op.String())

		oldC := C
		if oldC == nil {
			logger.Warn("Cannot reload config: current config is nil")
			return
		}

		newC := &Config{}
		if err := viper.Unmarshal(newC); err != nil {
			logger.Errorf("Config reload failed (unmarshal error): %v, keeping old config", err)
			C = oldC // 回滚到旧配置
			return
		}

		// 验证新配置
		if err := newC.Validate(); err != nil {
			logger.Errorf("Config reload failed (validation error): %v, keeping old config", err)
			C = oldC // 回滚到旧配置
			return
		}

		// 检测配置变更并记录
		changes := detectConfigChanges(oldC, newC)
		if len(changes) > 0 {
			logger.Infof("Config changes detected:")
			for _, change := range changes {
				logger.Infof("  - %s: %v -> %v", change.Key, change.OldValue, change.NewValue)
			}
		}

		// 应用新配置
		C = newC
		logger.Info("Config reloaded successfully")
	})

	return cfg, nil
}

// GetConfig 获取全局配置（向后兼容）
func GetConfig() *Config {
	return C
}

// detectConfigChanges 检测配置变更
func detectConfigChanges(old, new *Config) []ConfigChange {
	var changes []ConfigChange

	// 检测 Server 配置变更
	if old.Server.Port != new.Server.Port {
		changes = append(changes, ConfigChange{
			Key:      "server.port",
			OldValue: old.Server.Port,
			NewValue: new.Server.Port,
		})
	}
	if old.Server.Mode != new.Server.Mode {
		changes = append(changes, ConfigChange{
			Key:      "server.mode",
			OldValue: old.Server.Mode,
			NewValue: new.Server.Mode,
		})
	}

	// 检测 Database 配置变更
	if old.Database.Host != new.Database.Host {
		changes = append(changes, ConfigChange{
			Key:      "database.host",
			OldValue: old.Database.Host,
			NewValue: new.Database.Host,
		})
	}
	if old.Database.Port != new.Database.Port {
		changes = append(changes, ConfigChange{
			Key:      "database.port",
			OldValue: old.Database.Port,
			NewValue: new.Database.Port,
		})
	}
	if old.Database.DBName != new.Database.DBName {
		changes = append(changes, ConfigChange{
			Key:      "database.dbname",
			OldValue: old.Database.DBName,
			NewValue: new.Database.DBName,
		})
	}
	if old.Database.MaxIdleConns != new.Database.MaxIdleConns {
		changes = append(changes, ConfigChange{
			Key:      "database.max_idle_conns",
			OldValue: old.Database.MaxIdleConns,
			NewValue: new.Database.MaxIdleConns,
		})
	}
	if old.Database.MaxOpenConns != new.Database.MaxOpenConns {
		changes = append(changes, ConfigChange{
			Key:      "database.max_open_conns",
			OldValue: old.Database.MaxOpenConns,
			NewValue: new.Database.MaxOpenConns,
		})
	}

	// 检测 Redis 配置变更
	if old.Redis.Addr != new.Redis.Addr {
		changes = append(changes, ConfigChange{
			Key:      "redis.addr",
			OldValue: old.Redis.Addr,
			NewValue: new.Redis.Addr,
		})
	}
	if old.Redis.DB != new.Redis.DB {
		changes = append(changes, ConfigChange{
			Key:      "redis.db",
			OldValue: old.Redis.DB,
			NewValue: new.Redis.DB,
		})
	}

	// 检测 Auth 配置变更
	if old.Auth.AccessTokenExpires != new.Auth.AccessTokenExpires {
		changes = append(changes, ConfigChange{
			Key:      "auth.access_token_expires",
			OldValue: old.Auth.AccessTokenExpires,
			NewValue: new.Auth.AccessTokenExpires,
		})
	}
	if old.Auth.RefreshTokenExpires != new.Auth.RefreshTokenExpires {
		changes = append(changes, ConfigChange{
			Key:      "auth.refresh_token_expires",
			OldValue: old.Auth.RefreshTokenExpires,
			NewValue: new.Auth.RefreshTokenExpires,
		})
	}
	if old.Auth.Issuer != new.Auth.Issuer {
		changes = append(changes, ConfigChange{
			Key:      "auth.issuer",
			OldValue: old.Auth.Issuer,
			NewValue: new.Auth.Issuer,
		})
	}
	if old.Auth.Audience != new.Auth.Audience {
		changes = append(changes, ConfigChange{
			Key:      "auth.audience",
			OldValue: old.Auth.Audience,
			NewValue: new.Auth.Audience,
		})
	}
	if old.Auth.LoginMaxRetries != new.Auth.LoginMaxRetries {
		changes = append(changes, ConfigChange{
			Key:      "auth.login_max_retries",
			OldValue: old.Auth.LoginMaxRetries,
			NewValue: new.Auth.LoginMaxRetries,
		})
	}
	if old.Auth.LoginLockWindow != new.Auth.LoginLockWindow {
		changes = append(changes, ConfigChange{
			Key:      "auth.login_lock_window",
			OldValue: old.Auth.LoginLockWindow,
			NewValue: new.Auth.LoginLockWindow,
		})
	}

	// 检测 App 配置变更（敏感信息只显示是否变更，不显示值）
	if old.App.JWTSecret != new.App.JWTSecret {
		changes = append(changes, ConfigChange{
			Key:      "app.jwt_secret",
			OldValue: "[REDACTED]",
			NewValue: "[REDACTED]",
		})
	}
	if old.App.EncryptionKey != new.App.EncryptionKey {
		changes = append(changes, ConfigChange{
			Key:      "app.encryption_key",
			OldValue: "[REDACTED]",
			NewValue: "[REDACTED]",
		})
	}

	// 检测 Log 配置变更
	if old.Log.LogLevel != new.Log.LogLevel {
		changes = append(changes, ConfigChange{
			Key:      "log.log_level",
			OldValue: old.Log.LogLevel,
			NewValue: new.Log.LogLevel,
		})
	}
	if old.Log.LogDir != new.Log.LogDir {
		changes = append(changes, ConfigChange{
			Key:      "log.log_dir",
			OldValue: old.Log.LogDir,
			NewValue: new.Log.LogDir,
		})
	}
	if old.Log.InfoLogName != new.Log.InfoLogName {
		changes = append(changes, ConfigChange{
			Key:      "log.info_log_name",
			OldValue: old.Log.InfoLogName,
			NewValue: new.Log.InfoLogName,
		})
	}
	if old.Log.ErrorLogName != new.Log.ErrorLogName {
		changes = append(changes, ConfigChange{
			Key:      "log.error_log_name",
			OldValue: old.Log.ErrorLogName,
			NewValue: new.Log.ErrorLogName,
		})
	}
	if old.Log.EnableConsole != new.Log.EnableConsole {
		changes = append(changes, ConfigChange{
			Key:      "log.enable_console",
			OldValue: old.Log.EnableConsole,
			NewValue: new.Log.EnableConsole,
		})
	}

	return changes
}

// Validate 验证配置
func (c *Config) Validate() error {
	var errors []string

	// 验证 Server
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		errors = append(errors, "server.port must be between 1 and 65535")
	}

	// 验证 Database
	if c.Database.Username == "" {
		errors = append(errors, "database.username is required")
	}
	if c.Database.Password == "" {
		errors = append(errors, "database.password is required")
	}
	if c.Database.Host == "" {
		errors = append(errors, "database.host is required")
	}
	if c.Database.DBName == "" {
		errors = append(errors, "database.dbname is required")
	}
	if c.Database.Port < 1 || c.Database.Port > 65535 {
		errors = append(errors, "database.port must be between 1 and 65535")
	}
	if c.Database.MaxIdleConns < 0 {
		errors = append(errors, "database.max_idle_conns must be >= 0")
	}
	if c.Database.MaxOpenConns < 1 {
		errors = append(errors, "database.max_open_conns must be >= 1")
	}
	if c.Database.ConnMaxLifetime < 0 {
		errors = append(errors, "database.conn_max_lifetime must be >= 0")
	}

	// 验证 Redis
	if c.Redis.Addr == "" {
		errors = append(errors, "redis.addr is required")
	}
	if c.Redis.DB < 0 || c.Redis.DB > 15 {
		errors = append(errors, "redis.db must be between 0 and 15")
	}
	if c.Redis.PoolSize < 0 {
		errors = append(errors, "redis.pool_size must be >= 0")
	}

	// 验证 Auth
	if c.Auth.AccessTokenExpires <= 0 {
		errors = append(errors, "auth.access_token_expires must be > 0")
	}
	if c.Auth.RefreshTokenExpires <= 0 {
		errors = append(errors, "auth.refresh_token_expires must be > 0")
	}
	if c.Auth.RefreshTokenExpires <= c.Auth.AccessTokenExpires {
		errors = append(errors, "auth.refresh_token_expires must be > auth.access_token_expires")
	}
	if c.Auth.LoginMaxRetries < 0 {
		errors = append(errors, "auth.login_max_retries must be >= 0")
	}
	if c.Auth.LoginLockWindow < 0 {
		errors = append(errors, "auth.login_lock_window must be >= 0")
	}

	// 验证 App
	if c.App.JWTSecret == "" {
		errors = append(errors, "app.jwt_secret is required")
	}
	if len(c.App.JWTSecret) < 16 {
		errors = append(errors, "app.jwt_secret must be at least 16 characters")
	}
	if c.App.EncryptionKey == "" {
		errors = append(errors, "app.encryption_key is required")
	}
	if len(c.App.EncryptionKey) != 64 {
		errors = append(errors, "app.encryption_key must be 64 characters (32 bytes hex)")
	}

	// 验证 Tekton
	if c.K8s.Tekton.Namespace == "" {
		errors = append(errors, "kubernetes.tekton.namespace is required")
	}

	// 验证 K8s
	if c.K8s.Timeout < 0 {
		errors = append(errors, "kubernetes.timeout must be >= 0")
	}
	if c.K8s.QPS < 0 {
		errors = append(errors, "kubernetes.qps must be >= 0")
	}
	if c.K8s.Burst < 0 {
		errors = append(errors, "kubernetes.burst must be >= 0")
	}

	// 验证 Prometheus
	if c.Prometheus.Target.Namespace == "" {
		errors = append(errors, "prometheus.target.namespace is required")
	}
	if c.Prometheus.Target.ConfigMap == "" {
		errors = append(errors, "prometheus.target.configMap is required")
	}
	if c.Prometheus.Rule.Namespace == "" {
		errors = append(errors, "prometheus.rule.namespace is required")
	}
	if c.Prometheus.Rule.ConfigMap == "" {
		errors = append(errors, "prometheus.rule.configmap is required")
	}

	// 验证 permission_worker（可选段，若配置了则校验范围；0 表示使用默认值）
	if c.PermissionWorker.BatchSize < 0 || c.PermissionWorker.BatchSize > permissionWorkerMaxBatchSize*2 {
		errors = append(errors, fmt.Sprintf("permission_worker.batch_size must be 0 (default) or 1～%d", permissionWorkerMaxBatchSize*2))
	}
	if c.PermissionWorker.MaxRetry < 0 || c.PermissionWorker.MaxRetry > 100 {
		errors = append(errors, "permission_worker.max_retry must be 0 (default) or 1～100")
	}
	if c.PermissionWorker.ScanInterval != "" {
		if d, err := time.ParseDuration(c.PermissionWorker.ScanInterval); err != nil || d <= 0 {
			errors = append(errors, "permission_worker.scan_interval must be a positive duration (e.g. 10s, 2m)")
		}
	}
	if c.PermissionWorker.BaseDelay != "" {
		if d, err := time.ParseDuration(c.PermissionWorker.BaseDelay); err != nil || d < 0 {
			errors = append(errors, "permission_worker.base_delay must be a non-negative duration (e.g. 10s)")
		}
	}
	if c.PermissionWorker.MaxDelay != "" {
		if d, err := time.ParseDuration(c.PermissionWorker.MaxDelay); err != nil || d < 0 {
			errors = append(errors, "permission_worker.max_delay must be a non-negative duration (e.g. 15m)")
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors:\n  - %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// GetLoggerConfig 将 LogConfig 转换为 logger.Config
func (lc *LogConfig) GetLoggerConfig() logger.Config {
	return logger.Config{
		LogDir:        lc.LogDir,
		InfoLogName:   lc.InfoLogName,
		ErrorLogName:  lc.ErrorLogName,
		MaxSize:       lc.MaxSize,
		MaxBackups:    lc.MaxBackups,
		MaxAge:        lc.MaxAge,
		Compress:      lc.Compress,
		LogLevel:      lc.LogLevel,
		EnableConsole: lc.EnableConsole,
	}
}
