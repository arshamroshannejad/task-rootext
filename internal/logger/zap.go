package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func New(debug bool) (*zap.Logger, error) {
	if debug {
		return development()
	}
	return production()
}

func development() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	cfg.EncoderConfig.TimeKey = zapcore.OmitKey
	return cfg.Build()
}

func production() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.TimeKey = zapcore.OmitKey
	return cfg.Build()
}
