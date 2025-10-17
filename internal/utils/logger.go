package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
}

// Sync flushes any buffered log entries (call this on shutdown)
func Sync() {
	_ = Logger.Sync()
}
