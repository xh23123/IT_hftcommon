package common

import (
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var MaxSize int = 100
var MaxAge int = 30
var MaxBackups int = 30

func InitLogger(logpath string, loglevel string, isK8sLog bool) {
	hook := lumberjack.Logger{
		Filename:   logpath, //日志文件路径
		MaxSize:    MaxSize, //最大MB
		MaxAge:     MaxAge,
		MaxBackups: MaxBackups,
		Compress:   true,
	}
	// 设置日志级别,debug可以打印出info,debug,warn；info级别可以打印warn，info；warn只能打印warn
	// debug->info->warn->error
	var level zapcore.Level
	switch loglevel {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	default:
		level = zap.InfoLevel
	}

	var syncer zapcore.WriteSyncer

	if isK8sLog {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook))
	} else {
		syncer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(&hook))
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		syncer,
		level,
	)
	Logger = zap.New(core)
}
