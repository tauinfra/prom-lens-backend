package logger

import (
	"github.com/go-lumberjack/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	once      sync.Once
	ZapLogger *zap.Logger
)

type Config struct {
	InfoLogPath  string
	ErrorLogPath string
	MaxSize      int    // 日志文件最大大小(MB)
	MaxBackups   int    // 保留的旧日志文件最大数量
	MaxAge       int    // 保留旧日志文件的最大天数
	Compress     bool   // 是否压缩/归档旧日志文件
	LogLevel     string // 日志级别: debug, info, warn, error, dpanic, panic, fatal
}

// getLogWriter 创建带日志轮转的写入器
func getLogWriter(filename string, maxSize, maxBackup, maxAge int, compress bool) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,   // MB
		MaxBackups: maxBackup, // 最大备份数量
		MaxAge:     maxAge,    // 最大保存天数
		Compress:   compress,  // 是否压缩
	}

	// 同时输出到文件和控制台
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(lumberJackLogger),
		zapcore.AddSync(os.Stdout),
	)
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

// parseLogLevel 解析日志级别字符串
func parseLogLevel(level string) zapcore.Level {
	switch level {
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

// InitLogger 初始化日志记录器(线程安全)
func InitLogger(cfg Config) {
	once.Do(func() {
		// 设置日志级别
		logLevel := parseLogLevel(cfg.LogLevel)

		// 创建核心
		encoder := zapcore.NewJSONEncoder(getEncoderConfig())

		// 信息日志核心(包含info及以上级别)
		infoCore := zapcore.NewCore(
			encoder,
			getLogWriter(cfg.InfoLogPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= logLevel && level < zapcore.ErrorLevel
			}),
		)

		// 错误日志核心(只包含error及以上级别)
		errorCore := zapcore.NewCore(
			encoder,
			getLogWriter(cfg.ErrorLogPath, cfg.MaxSize, cfg.MaxBackups, cfg.MaxAge, cfg.Compress),
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
	})

}

// Sync 刷新日志缓冲区
func Sync() error {
	if ZapLogger != nil {
		return ZapLogger.Sync()
	}
	return nil
}

// Debug 日志级别方法
func Debug(msg string, fields ...zap.Field) {
	ZapLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	ZapLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	ZapLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	ZapLogger.Error(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	ZapLogger.DPanic(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	ZapLogger.Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	ZapLogger.Fatal(msg, fields...)
}

// Debugf 格式化日志方法
func Debugf(template string, args ...interface{}) {
	ZapLogger.Sugar().Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	ZapLogger.Sugar().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	ZapLogger.Sugar().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	ZapLogger.Sugar().Errorf(template, args...)
}

func DPanicf(template string, args ...interface{}) {
	ZapLogger.Sugar().DPanicf(template, args...)
}

func Panicf(template string, args ...interface{}) {
	ZapLogger.Sugar().Panicf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	ZapLogger.Sugar().Fatalf(template, args...)
}
