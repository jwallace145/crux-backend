package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/utils"
)

// LoggerMiddleware logs all incoming requests and outgoing responses.
// It provides a comprehensive audit trail for all API calls including
// request/response details, duration, and status codes.
//
// The middleware automatically generates or uses an existing X-Request-ID
// header for request tracking and makes it available via c.Locals("request_id").
func LoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Generate request ID for tracking
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
			c.Set("X-Request-ID", requestID)
		}

		// Store request ID in context for handlers to access
		c.Locals("request_id", requestID)

		// Get logger from context with request ID
		log := utils.GetLoggerFromContext(c)

		// Log incoming request
		log.Info("Incoming API request",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
			zap.ByteString("body", c.Body()),
			zap.String("query", string(c.Request().URI().QueryString())),
		)

		// Process request
		err := c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log response
		statusCode := c.Response().StatusCode()
		logLevel := zap.InfoLevel
		if statusCode >= 500 {
			logLevel = zap.ErrorLevel
		} else if statusCode >= 400 {
			logLevel = zap.WarnLevel
		}

		log.Log(logLevel, "API request completed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
			zap.Int("response_size", len(c.Response().Body())),
		)

		return err
	}
}
