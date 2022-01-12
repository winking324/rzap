package rzap

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func encoder() zapcore.Encoder {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewJSONEncoder(config)
}

func NewCore(writer *lumberjack.Logger, level zapcore.LevelEnabler) zapcore.Core {
	syncer := zapcore.AddSync(writer)
	return zapcore.NewCore(encoder(), syncer, level)
}

func NewLogger(cores []zapcore.Core, opts ...zap.Option) *zap.Logger {
	return zap.New(zapcore.NewTee(cores...), opts...)
}

func NewGlobalLogger(cores []zapcore.Core, opts ...zap.Option) {
	zap.ReplaceGlobals(NewLogger(cores, opts...))
}
