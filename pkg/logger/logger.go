package logger

import (
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// Logger 全局日志实例
	Logger *zap.Logger
	// Sugar 全局 SugaredLogger 实例（提供更友好的 API）
	Sugar *zap.SugaredLogger
)

// Config 日志配置
type Config struct {
	Level      string // 日志级别: debug, info, warn, error
	FilePath   string // 日志文件路径
	MaxSize    int    // 单个日志文件最大大小（MB）
	MaxBackups int    // 保留旧文件的最大个数
	MaxAge     int    // 保留旧文件的最大天数
	Compress   bool   // 是否压缩旧文件
	Debug      bool   // 是否开启开发模式
}

// Init 初始化日志系统
func Init(config *Config) error {
	// 设置默认值
	if config.MaxSize == 0 {
		config.MaxSize = 100
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 7
	}
	if config.MaxAge == 0 {
		config.MaxAge = 30
	}

	// 解析日志级别
	level := parseLevel(config.Level)

	// 创建日志目录
	if config.FilePath != "" {
		dir := filepath.Dir(config.FilePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	}

	// 配置编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 开发模式使用更友好的编码器
	var encoder zapcore.Encoder
	if config.Debug {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 配置输出
	var cores []zapcore.Core

	// 控制台输出
	consoleOutput := zapcore.AddSync(os.Stdout)
	cores = append(cores, zapcore.NewCore(encoder, consoleOutput, level))

	// 文件输出
	if config.FilePath != "" {
		fileWriter := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
			LocalTime:  true,
		})

		// JSON 格式写入文件
		jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(jsonEncoder, fileWriter, level))
	}

	// 创建 logger
	core := zapcore.NewTee(cores...)

	var logger *zap.Logger
	if config.Debug {
		logger = zap.New(core, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 设置全局 logger
	Logger = logger
	Sugar = logger.Sugar()

	return nil
}

// parseLevel 解析日志级别
func parseLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Sync 同步日志缓冲区
func Sync() {
	if Logger != nil {
		_ = Logger.Sync()
	}
}

// 便捷方法

// Debug 调试日志
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info 信息日志
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn 警告日志
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 错误日志
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 致命错误日志
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Panic Panic日志
func Panic(msg string, fields ...zap.Field) {
	Logger.Panic(msg, fields...)
}

// With 添加字段
func With(fields ...zap.Field) *zap.Logger {
	return Logger.With(fields...)
}
