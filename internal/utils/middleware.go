package utils

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

// RequestLogger is a middleware that logs all incoming requests and outgoing responses
// This provides a comprehensive audit trail for all API calls
func RequestLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Start timer
		start := time.Now()

		// Generate request ID for tracking
		requestID := c.Get("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
			c.Set("X-Request-ID", requestID)
		}

		// Log incoming request
		Logger.Info("Incoming API request",
			zap.String("request_id", requestID),
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

		Logger.Log(logLevel, "API request completed",
			zap.String("request_id", requestID),
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status_code", statusCode),
			zap.Duration("duration", duration),
			zap.Int("response_size", len(c.Response().Body())),
		)

		return err
	}
}

// generateRequestID generates a unique request ID for tracking
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// CORSMiddleware configures and returns CORS middleware for the API
// This configuration is permissive for development - should be locked down for production
func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		// Allow common development origins
		// NOTE: Cannot use "*" with AllowCredentials: true due to browser security restrictions
		// TODO(#1): Lock this down to specific production origins when deploying
		AllowOrigins: "http://localhost:3000,http://localhost:3001,http://localhost:3002,http://localhost:5173,http://localhost:8080,http://127.0.0.1:3000,http://127.0.0.1:5173,http://127.0.0.1:8080",

		// Allow credentials (cookies, authorization headers, etc.)
		// CRITICAL: This is required for JWT cookie-based authentication to work
		AllowCredentials: true,

		// Allow common HTTP methods including OPTIONS for preflight
		AllowMethods: "GET,POST,PUT,PATCH,DELETE,OPTIONS,HEAD",

		// Allow common headers including custom headers the client might send
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-Request-ID,X-Requested-With,Access-Control-Request-Method,Access-Control-Request-Headers",

		// Expose headers that client-side JavaScript can access
		ExposeHeaders: "X-Request-ID,Content-Length",

		// Cache preflight requests for 12 hours
		MaxAge: 43200,

		// CRITICAL: This ensures OPTIONS requests are always handled by CORS middleware
		// even if no route explicitly handles OPTIONS
		Next: func(c *fiber.Ctx) bool {
			// Never skip CORS for any request
			return false
		},
	})
}
