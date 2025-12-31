package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
)

var (
	log   *zap.Logger
	sugar *zap.SugaredLogger
	once  sync.Once
)

// Config 日志配置
type Config struct {
	Level      string // 日志级别: debug, info, warn, error, fatal
	FileName   string // 日志文件路径
	MaxSize    int    // 单文件最大尺寸(MB)
	MaxBackups int    // 保留旧文件最大数量
	MaxAge     int    // 保留旧文件最大天数
	Compress   bool   // 是否压缩
	Console    bool   // 是否输出到控制台
}

// Init 初始化日志
func Init(cfg *Config) error {
	var err error
	once.Do(func() {
		err = initLogger(cfg)
	})
	return err
}

// initLogger 初始化日志实例
func initLogger(cfg *Config) error {
	// 解析日志级别
	level := zapcore.InfoLevel
	if cfg != nil && cfg.Level != "" {
		_ = level.UnmarshalText([]byte(cfg.Level))
	}

	// 编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建多个输出目标
	var cores []zapcore.Core

	// 文件输出
	if cfg != nil && cfg.FileName != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   cfg.FileName,
			MaxSize:    cfg.MaxSize,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAge,
			Compress:   cfg.Compress,
		}
		fileCore := zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig),
			zapcore.AddSync(fileWriter),
			level,
		)
		cores = append(cores, fileCore)
	}

	// 控制台输出
	if cfg == nil || cfg.Console {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		consoleCore := zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		)
		cores = append(cores, consoleCore)
	}

	// 组合core
	core := zapcore.NewTee(cores...)

	// 创建logger
	log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1), zap.AddStacktrace(zapcore.ErrorLevel))
	sugar = log.Sugar()

	return nil
}

// Debug 日志
func Debug(msg string, fields ...zap.Field) {
	log.Debug(msg, fields...)
}

// Info 日志
func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

// Warn 日志
func Warn(msg string, fields ...zap.Field) {
	log.Warn(msg, fields...)
}

// Error 日志
func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

// Fatal 日志后退出
func Fatal(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

// Debugf 格式化日志
func Debugf(format string, args ...interface{}) {
	sugar.Debugf(format, args...)
}

// Infof 格式化日志
func Infof(format string, args ...interface{}) {
	sugar.Infof(format, args...)
}

// Warnf 格式化日志
func Warnf(format string, args ...interface{}) {
	sugar.Warnf(format, args...)
}

// Errorf 格式化日志
func Errorf(format string, args ...interface{}) {
	sugar.Errorf(format, args...)
}

// Fatalf 格式化日志后退出
func Fatalf(format string, args ...interface{}) {
	sugar.Fatalf(format, args...)
}

// With 创建子logger
func With(fields ...zap.Field) *zap.Logger {
	return log.With(fields...)
}

// Get 获取原生logger
func Get() *zap.Logger {
	if log == nil {
		_ = Init(nil)
	}
	return log
}

// Sync 刷新缓冲
func Sync() error {
	if log != nil {
		return log.Sync()
	}
	return nil
}
