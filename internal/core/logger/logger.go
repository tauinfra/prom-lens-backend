package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-lumberjack/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	once           sync.Once
	ZapLogger      *zap.Logger
	ZapSugar       *zap.SugaredLogger // 缓存的 Sugar logger
	fallbackLogger *zap.Logger        // 未初始化时的备用 logger
	fallbackSugar  *zap.SugaredLogger // 缓存的 fallback Sugar logger
)

func init() {
	// 初始化 fallback logger，使用标准输出
	fallbackLogger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(getEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	))
	// 初始化 fallback Sugar logger
	fallbackSugar = fallbackLogger.Sugar()
}

type Config struct {
	LogDir        string
	InfoLogName   string
	ErrorLogName  string
	MaxSize       int    // 日志文件最大大小(MB)
	MaxBackups    int    // 保留的旧日志文件最大数量
	MaxAge        int    // 保留旧日志文件的最大天数
	Compress      bool   // 是否压缩/归档旧日志文件
	LogLevel      string // 日志级别: debug, info, warn, error, dpanic, panic, fatal
	EnableConsole bool   // 是否输出到控制台（默认 true）
}

// ensureLogDir 确保日志目录存在，不存在则创建
func ensureLogDir(filename string) error {
	if filename == "" {
		return nil // 空路径表示不写文件，只输出到控制台
	}
	dir := filepath.Dir(filename)
	if dir == "." || dir == "" {
		return nil // 当前目录，无需创建
	}
	// 检查目录是否存在
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 目录不存在，创建它
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create log directory %s: %w", dir, err)
		}
	}
	return nil
}

// getLogWriter 创建带日志轮转的写入器
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress, enableConsole bool) zapcore.WriteSyncer {
	// 文件名为空时，根据控制台开关决定输出目的地
	if filename == "" {
		if enableConsole {
			return zapcore.AddSync(os.Stdout)
		}
		return zapcore.AddSync(io.Discard)
	}

	// 确保日志目录存在
	if err := ensureLogDir(filename); err != nil {
		// 如果创建目录失败，根据配置决定输出
		if enableConsole {
			fmt.Fprintf(os.Stderr, "Warning: %v, logging to stdout only\n", err)
			return zapcore.AddSync(os.Stdout)
		}
		return zapcore.AddSync(io.Discard)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   // MB
		MaxBackups: maxBackup, // 最大备份数量
		MaxAge:     maxAge,    // 最大保存天数
		Compress:   compress,  // 是否压缩
	}

	// 根据配置决定是否同时输出到控制台
	if enableConsole {
		return zapcore.NewMultiWriteSyncer(
			zapcore.AddSync(lumberJackLogger),
			zapcore.AddSync(os.Stdout),
		)
	}
	// 只输出到文件
	return zapcore.AddSync(lumberJackLogger)
}

// getEncoderConfig 自定义日志格式
func getEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,                        // 日志级别大写显示
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 自定义时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// parseLogLevel 解析日志级别字符串（大小写不敏感）
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// validateConfig 验证日志配置并设置默认值
func validateConfig(cfg *Config) error {
	var errors []string

	// 验证日志目录（如果提供）
	if cfg.LogDir != "" {
		normalized, err := normalizeLogDir(cfg.LogDir)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			cfg.LogDir = normalized
			if info, err := os.Stat(cfg.LogDir); err == nil {
				if !info.IsDir() {
					errors = append(errors, fmt.Sprintf("log dir is not a directory: %s", cfg.LogDir))
				}
			}
		}
	}
	// 默认日志文件名
	if cfg.InfoLogName == "" {
		cfg.InfoLogName = "info.log"
	}
	if cfg.ErrorLogName == "" {
		cfg.ErrorLogName = "error.log"
	}

	// 验证并设置 MaxSize 默认值（MB）
	if cfg.MaxSize < 0 {
		errors = append(errors, "max_size must be >= 0")
	} else if cfg.MaxSize == 0 {
		cfg.MaxSize = 100 // 默认 100MB
	}

	// 验证并设置 MaxBackups 默认值
	if cfg.MaxBackups < 0 {
		errors = append(errors, "max_backups must be >= 0")
	} else if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 10 // 默认 10 个备份
	}

	// 验证并设置 MaxAge 默认值（天）
	if cfg.MaxAge < 0 {
		errors = append(errors, "max_age must be >= 0")
	} else if cfg.MaxAge == 0 {
		cfg.MaxAge = 30 // 默认 30 天
	}

	// 至少启用一个输出目的地
	if !cfg.EnableConsole && cfg.LogDir == "" {
		errors = append(errors, "logger output is disabled: enable_console is false and log paths are empty")
	}

	if len(errors) > 0 {
		return fmt.Errorf("logger config validation failed. err: %s", strings.Join(errors, "\n  - "))
	}

	return nil
}

// normalizeLogDir 规范化日志目录
func normalizeLogDir(path string) (string, error) {
	cleaned := filepath.Clean(path)
	return cleaned, nil
}

// InitLogger 初始化日志记录器(线程安全)
func InitLogger(cfg Config) error {
	var initErr error
	once.Do(func() {
		// 验证配置并设置默认值
		if err := validateConfig(&cfg); err != nil {
			initErr = err
			return
		}

		// 设置日志级别
		logLevel := parseLogLevel(cfg.LogLevel)

		// 创建核心
		encoder := zapcore.NewJSONEncoder(getEncoderConfig())

		// 信息日志核心(包含info及以上级别)
		infoCore := zapcore.NewCore(
			encoder,
			getLogWriter(filepath.Join(cfg.LogDir, cfg.InfoLogName), cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress, cfg.EnableConsole),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= logLevel && level < zapcore.ErrorLevel
			}),
		)

		// 错误日志核心(只包含error及以上级别)
		errorCore := zapcore.NewCore(
			encoder,
			getLogWriter(filepath.Join(cfg.LogDir, cfg.ErrorLogName), cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress, cfg.EnableConsole),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= zapcore.ErrorLevel
			}),
		)

		// 合并核心
		core := zapcore.NewTee(infoCore, errorCore)

		// 创建logger
		ZapLogger = zap.New(core,
			zap.AddCaller(),
			zap.AddCallerSkip(1),                  // 跳过一层调用栈
			zap.AddStacktrace(zapcore.ErrorLevel), // 错误级别添加堆栈跟踪
		)
		// 缓存 Sugar logger 以提升性能
		ZapSugar = ZapLogger.Sugar()
	})

	return initErr
}

// Sync 刷新日志缓冲区
func Sync() error {
	if ZapLogger != nil {
		return ZapLogger.Sync()
	}
	return nil
}

// getLogger 获取 logger 实例，如果未初始化则返回 fallback
func getLogger() *zap.Logger {
	if ZapLogger != nil {
		return ZapLogger
	}
	return fallbackLogger
}

// getSugar 获取 Sugar logger 实例，如果未初始化则返回 fallback
func getSugar() *zap.SugaredLogger {
	if ZapSugar != nil {
		return ZapSugar
	}
	return fallbackSugar
}

// Debug 日志级别方法
func Debug(msg string, fields ...zap.Field) {
	getLogger().Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	getLogger().Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLogger().Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLogger().Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	getLogger().DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	getLogger().Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLogger().Fatal(msg, fields...)
}

// Debugf 格式化日志方法（使用缓存的 Sugar logger 提升性能）
func Debugf(template string, args ...interface{}) {
	getSugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	getSugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	getSugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	getSugar().Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	getSugar().DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	getSugar().Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	getSugar().Fatalf(template, args...)
}
