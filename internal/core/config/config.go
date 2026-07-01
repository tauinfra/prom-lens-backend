package config

import (
	"fmt"
	"strings"
	"time"
	"prom-lens-backend/internal/core/logger"

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
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Auth       AuthConfig       `mapstructure:"auth"`
	Log        LogConfig        `mapstructure:"logging"`
	App        AppConfig        `mapstructure:"app"`
	BaseURL    string           `mapstructure:"base_url"`
	Prometheus PrometheusConfig `mapstructure:"prometheus"`
	Alerting   AlertingConfig   `mapstructure:"alerting"`
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

type AuthConfig struct {
	AccessTokenExpires  time.Duration `mapstructure:"access_token_expires"`
	RefreshTokenExpires time.Duration `mapstructure:"refresh_token_expires"`
	Whitelist           []string      `mapstructure:"whitelist"`
	Issuer              string        `mapstructure:"issuer"`
	Audience            string        `mapstructure:"audience"`
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
	JWTSecret string `mapstructure:"jwt_secret"`
	JWTExpire string `mapstructure:"jwt_expire"`
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

type AlertingConfig struct {
	LarkTimeout  time.Duration          `mapstructure:"lark_timeout"`
	Alertmanager AlertmanagerSyncConfig `mapstructure:"alertmanager"`
}

type AlertmanagerSyncConfig struct {
	Namespace       string `mapstructure:"namespace"`
	ConfigMap       string `mapstructure:"configmap"`
	ConfigKey         string `mapstructure:"config_key"`
	DefaultReceiver string `mapstructure:"default_receiver"`
}

var C *Config // Deprecated: 使用 InitConfig 返回的配置，全局变量仅用于向后兼容

// bindConfigEnvs 将配置项绑定到环境变量，Unmarshal 时才会覆盖 YAML 中的值。
// 变量名规则：PROM_LENS_ + 配置键（点号换为下划线），例如 database.password -> PROM_LENS_DATABASE_PASSWORD
func bindConfigEnvs() error {
	keys := []string{
		"server.host",
		"server.port",
		"server.mode",
		"database.username",
		"database.password",
		"database.host",
		"database.dbname",
		"database.port",
		"database.max_idle_conns",
		"database.max_open_conns",
		"database.conn_max_lifetime",
		"app.jwt_secret",
		"app.jwt_expire",
		"base_url",
		"auth.access_token_expires",
		"auth.refresh_token_expires",
		"auth.issuer",
		"auth.audience",
		"logging.log_dir",
		"logging.log_level",
		"logging.enable_console",
		"prometheus.target.namespace",
		"prometheus.target.configMap",
		"prometheus.rule.namespace",
		"prometheus.rule.configmap",
		"alerting.lark_timeout",
		"alerting.alertmanager.namespace",
		"alerting.alertmanager.configmap",
		"alerting.alertmanager.config_key",
		"alerting.alertmanager.default_receiver",
	}
	for _, key := range keys {
		if err := viper.BindEnv(key); err != nil {
			return fmt.Errorf("bind env for %s: %w", key, err)
		}
	}
	return nil
}

// InitConfig 初始化配置
// 支持通过环境变量覆盖配置值，环境变量命名规则：
//   - 前缀：PROM_LENS（自动拼接下划线，勿写成 PROM_LENS_）
//   - 嵌套结构使用下划线分隔，例如：
//     PROM_LENS_SERVER_PORT -> server.port
//     PROM_LENS_DATABASE_HOST -> database.host
//     PROM_LENS_AUTH_ACCESS_TOKEN_EXPIRES -> auth.access_token_expires
//   - 数组类型暂不支持通过环境变量设置
//
// 示例：
//
//	export PROM_LENS_SERVER_PORT=9090
//	export PROM_LENS_DATABASE_HOST=192.168.1.1
//	export PROM_LENS_DATABASE_PASSWORD=mypassword
func InitConfig(cfgFile string) (*Config, error) {
	viper.SetConfigFile(cfgFile)

	// 设置环境变量前缀
	viper.SetEnvPrefix("PROM_LENS")
	// 设置环境变量键名替换器：将点号替换为下划线，支持嵌套结构
	// 例如：PROM_LENS_SERVER_PORT -> server.port
	//      PROM_LENS_DATABASE_HOST -> database.host
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()
	if err := bindConfigEnvs(); err != nil {
		return nil, fmt.Errorf("bind environment variables failed: %w", err)
	}

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
	// 检测 App 配置变更（敏感信息只显示是否变更，不显示值）
	if old.App.JWTSecret != new.App.JWTSecret {
		changes = append(changes, ConfigChange{
			Key:      "app.jwt_secret",
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
	// 验证 App
	if c.App.JWTSecret == "" {
		errors = append(errors, "app.jwt_secret is required")
	}
	if len(c.App.JWTSecret) < 16 {
		errors = append(errors, "app.jwt_secret must be at least 16 characters")
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
