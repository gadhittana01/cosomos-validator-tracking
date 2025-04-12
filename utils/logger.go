package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var MapInfoLevel = map[string]zapcore.Level{
	"debug": zap.DebugLevel,
	"info":  zap.InfoLevel,
	"warn":  zap.WarnLevel,
	"error": zap.ErrorLevel,
}

type Logger struct {
	logger *zap.Logger
}

type LoggerSvc interface {
	Info(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
	Panic(message string, fields ...zap.Field)
	Debug(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Warn(message string, fields ...zap.Field)
}

func NewLogger(config *BaseConfig) LoggerSvc {
	var err error
	logCfg := zap.NewProductionConfig()
	logCfg.Level = zap.NewAtomicLevelAt(MapInfoLevel[config.LogLevel])

	client, err := logCfg.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}

	return &Logger{
		logger: client,
	}
}

func (l *Logger) Info(message string, fields ...zap.Field) {
	l.logger.Info(message, fields...)
}

func (l *Logger) Fatal(message string, fields ...zap.Field) {
	l.logger.Fatal(message, fields...)
}

func (l *Logger) Panic(message string, fields ...zap.Field) {
	l.logger.Panic(message, fields...)
}

func (l *Logger) Debug(message string, fields ...zap.Field) {
	l.logger.Debug(message, fields...)
}

func (l *Logger) Error(message string, fields ...zap.Field) {
	l.logger.Error(message, fields...)
}

func (l *Logger) Warn(message string, fields ...zap.Field) {
	l.logger.Warn(message, fields...)
}
