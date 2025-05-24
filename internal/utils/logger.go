package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger is the global logger instance
var Logger *zap.Logger

// InitLogger initializes the logger
func InitLogger(debug bool) error {
	config := zap.NewProductionConfig()
	if debug {
		config.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return err
	}

	Logger = logger
	return nil
}

// GetLogger returns the logger instance
func GetLogger() *zap.Logger {
	if Logger == nil {
		// Create a default logger if not initialized
		Logger, _ = zap.NewProduction()
	}
	return Logger
}
