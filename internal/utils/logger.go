package utils

import (
	"github.com/gofiber/fiber/v2"
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

// GetLoggerFromContext returns a logger with the request ID from the Fiber context
// If no request ID is found in context, it returns the base logger with "unknown" as request ID
// This ensures all logs from a request have the same request ID for traceability
func GetLoggerFromContext(c *fiber.Ctx) *zap.Logger {
	requestID, ok := c.Locals("request_id").(string)
	if !ok || requestID == "" {
		requestID = "unknown"
	}
	return Logger.With(zap.String("request_id", requestID))
}
