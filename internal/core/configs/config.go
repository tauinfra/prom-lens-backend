package configs

import (
	"fmt"
	"time"
	"valyria-backend/internal/core/logger"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Log      LogConfig      `mapstructure:"log"`
	Policy   PolicyConfig   `mapstructure:"k8s"`        // Casbin RBAC 策略
	K8s      K8sConfig      `mapstructure:"kubernetes"` // 注意：这里改为 kubernetes
	App      AppConfig      `mapstructure:"app"`
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
	Whitelist []string `mapstructure:"whitelist"`
}

type LogConfig struct {
	InfoLogPath  string `mapstructure:"info_log_path"`
	ErrorLogPath string `mapstructure:"error_log_path"`
	MaxSize      int    `mapstructure:"max_size"`
	MaxBackups   int    `mapstructure:"max_backups"`
	MaxAge       int    `mapstructure:"max_age"`
	Compress     bool   `mapstructure:"compress"`
	LogLevel     string `mapstructure:"log_level"`
}

type PolicyConfig struct {
	ModelPath   string `mapstructure:"model_path"`
	PolicyTable string `mapstructure:"policy_table"`
	AutoSave    bool   `mapstructure:"auto_save"`
}

type AppConfig struct {
	JWTSecret     string `mapstructure:"jwt_secret"`
	JWTExpire     string `mapstructure:"jwt_expire"`
	EncryptionKey string `mapstructure:"encryption_key"`
}

type K8sConfig struct {
	Timeout   int      `mapstructure:"timeout"`
	Resources []string `mapstructure:"resources"` // 新增：资源类型列表
	Actions   []string `mapstructure:"actions"`
}

var C *Config

func InitConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	viper.SetEnvPrefix("VALYRIA_")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("fatal read the configuration file failed, %w", err)
	}

	C = &Config{}
	if err := viper.Unmarshal(C); err != nil {
		return fmt.Errorf("unmarshal configuration file configs failed, %v", err)
	}

	// 监听配置文件变更
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		if err := viper.Unmarshal(C); err != nil {
			logger.Infof("Reload configs failed: %v", err)
		} else {
			logger.Infof("Reload configs successfully.")
		}
	})

	return nil
}

func GetConfig() *Config {
	return C
}

// GetLoggerConfig 将 LogConfig 转换为 logger.Config
func (lc *LogConfig) GetLoggerConfig() logger.Config {
	return logger.Config{
		InfoLogPath:  lc.InfoLogPath,
		ErrorLogPath: lc.ErrorLogPath,
		MaxSize:      lc.MaxSize,
		MaxBackups:   lc.MaxBackups,
		MaxAge:       lc.MaxAge,
		Compress:     lc.Compress,
		LogLevel:     lc.LogLevel,
	}
}
