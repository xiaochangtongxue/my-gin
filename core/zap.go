package core

import (
	"os"
	"time"

	"github.com/xiaochangtongxue/my-gin/global"
	"github.com/xiaochangtongxue/my-gin/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	options []zap.Option
	level   zapcore.Level
)

func Zap() *zap.Logger {
	// 创建日志根目录
	creatRoot()
	// 设置日志输出等级
	getLevel()
	if global.MGIN_CONFIG.Log.ShowLine {
		options = append(options, zap.AddCaller())
	}
	// 打印error的堆栈信息
	options = append(options, zap.AddStacktrace(zap.ErrorLevel))

	return zap.New(getZapCore(), options...)
}

func creatRoot() {
	if ok, _ := utils.PathExits(global.MGIN_CONFIG.Log.Root); !ok {
		os.Mkdir(global.MGIN_CONFIG.Log.Root, os.ModePerm)
	}
}

func getLevel() {
	switch global.MGIN_CONFIG.Log.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zapcore.DPanicLevel
	case "panic":
		level = zapcore.PanicLevel
	case "fatal":
		level = zapcore.FatalLevel
	}
}

func getZapCore() zapcore.Core {
	var zapEncoder zapcore.Encoder
	var allZapCore []zapcore.Core
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("[" + "2006-01-02 15:04:05.000" + "]"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	if global.MGIN_CONFIG.Log.Format == "json" {
		zapEncoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		zapEncoder = zapcore.NewConsoleEncoder(encoderConfig)
	}
	file := zapcore.NewCore(zapEncoder, getLogWrite(), level)
	console := zapcore.NewCore(zapEncoder, os.Stdout, level)

	allZapCore = append(allZapCore, file)
	allZapCore = append(allZapCore, console)

	return zapcore.NewTee(allZapCore...)
}

func getLogWrite() zapcore.WriteSyncer {
	file := &lumberjack.Logger{
		Filename:   global.MGIN_CONFIG.Log.Root + "/" + global.MGIN_CONFIG.Log.FileName,
		MaxSize:    global.MGIN_CONFIG.Log.MaxSize,
		MaxBackups: global.MGIN_CONFIG.Log.MaxBackups,
		MaxAge:     global.MGIN_CONFIG.Log.MaxAge,
		Compress:   global.MGIN_CONFIG.Log.Compress,
	}
	return zapcore.AddSync(file)
}
