package logger

import (
	"strings"

	"github.com/Delkira544/rakiduam/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var global *zap.Logger

func New(cfg config.LogConfig) error {
	var zcfg zap.Config
	if cfg.Format == "json" {
		zcfg = zap.NewProductionConfig()
		zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		zcfg = zap.NewDevelopmentConfig()
		zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	zcfg.Level = parseLevel(cfg.Level)
	zcfg.OutputPaths = []string{cfg.Output}
	l, err := zcfg.Build()
	if err != nil {
		return err
	}
	global = l
	return nil
}

func L() *zap.Logger {
	return global
}

func parseLevel(level string) zap.AtomicLevel {
	switch strings.ToLower(level) {
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn", "warning":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "fatal":
		return zap.NewAtomicLevelAt(zap.FatalLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel) // default
	}
}
