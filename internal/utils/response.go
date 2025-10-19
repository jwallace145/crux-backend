package utils

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/models"
)

var (
	// ServiceName is the name of the API service
	ServiceName = "crux-backend"

	// APIVersion is the current API version
	APIVersion = getEnvOrDefault("API_VERSION", "1.0.0")

	// Environment is the deployment environment (development, staging, production)
	Environment = getEnvOrDefault("APP_ENV", "development")
)

// ResponseBuilder provides a fluent interface for building API responses
type ResponseBuilder struct {
	response *models.APIResponse
}

// NewResponse creates a new ResponseBuilder with default values
func NewResponse(apiName string) *ResponseBuilder {
	return &ResponseBuilder{
		response: &models.APIResponse{
			ServiceName: ServiceName,
			Version:     APIVersion,
			Environment: Environment,
			APIName:     apiName,
			Timestamp:   time.Now().UTC(),
			Status:      models.StatusSuccess,
		},
	}
}

// WithData sets the data field of the response
func (rb *ResponseBuilder) WithData(data interface{}) *ResponseBuilder {
	rb.response.Data = data
	return rb
}

// WithMessage sets the message field of the response
func (rb *ResponseBuilder) WithMessage(message string) *ResponseBuilder {
	rb.response.Message = message
	return rb
}

// WithStatus sets the status field of the response
func (rb *ResponseBuilder) WithStatus(status string) *ResponseBuilder {
	rb.response.Status = status
	return rb
}

// WithError sets the error field and status to error
func (rb *ResponseBuilder) WithError(code, message string, details interface{}) *ResponseBuilder {
	rb.response.Status = models.StatusError
	rb.response.Error = &models.APIError{
		Code:    code,
		Message: message,
		Details: details,
	}
	return rb
}

// Build returns the constructed APIResponse
func (rb *ResponseBuilder) Build() *models.APIResponse {
	return rb.response
}

// SendJSON sends the response as JSON with the given status code
// It also logs the response for audit trail purposes
func (rb *ResponseBuilder) SendJSON(c *fiber.Ctx, statusCode int) error {
	requestID := c.Get("X-Request-ID", "unknown")

	// Log the response before sending
	if rb.response.Status == models.StatusError {
		if statusCode >= 500 {
			Logger.Error("API response error",
				zap.String("api", rb.response.APIName),
				zap.String("request_id", requestID),
				zap.Int("status_code", statusCode),
				zap.String("error_code", rb.response.Error.Code),
				zap.String("error_message", rb.response.Error.Message),
			)
		} else {
			Logger.Warn("API response warning",
				zap.String("api", rb.response.APIName),
				zap.String("request_id", requestID),
				zap.Int("status_code", statusCode),
				zap.String("error_code", rb.response.Error.Code),
				zap.String("error_message", rb.response.Error.Message),
			)
		}
	} else {
		Logger.Info("API response success",
			zap.String("api", rb.response.APIName),
			zap.String("request_id", requestID),
			zap.Int("status_code", statusCode),
		)
	}

	return c.Status(statusCode).JSON(rb.response)
}

// Helper functions for common response patterns

// SuccessResponse creates a successful response with data
func SuccessResponse(c *fiber.Ctx, apiName string, data interface{}, message string) error {
	return NewResponse(apiName).
		WithData(data).
		WithMessage(message).
		SendJSON(c, fiber.StatusOK)
}

// CreatedResponse creates a successful creation response
func CreatedResponse(c *fiber.Ctx, apiName string, data interface{}, message string) error {
	return NewResponse(apiName).
		WithData(data).
		WithMessage(message).
		SendJSON(c, fiber.StatusCreated)
}

// ErrorResponse creates an error response
func ErrorResponse(c *fiber.Ctx, apiName string, statusCode int, errorCode, message string, details interface{}) error {
	return NewResponse(apiName).
		WithError(errorCode, message, details).
		SendJSON(c, statusCode)
}

// BadRequestResponse creates a 400 Bad Request response
func BadRequestResponse(c *fiber.Ctx, apiName string, message string, details interface{}) error {
	return ErrorResponse(c, apiName, fiber.StatusBadRequest, models.ErrorCodeInvalidInput, message, details)
}

// NotFoundResponse creates a 404 Not Found response
func NotFoundResponse(c *fiber.Ctx, apiName string, message string) error {
	return ErrorResponse(c, apiName, fiber.StatusNotFound, models.ErrorCodeNotFound, message, nil)
}

// UnauthorizedResponse creates a 401 Unauthorized response
func UnauthorizedResponse(c *fiber.Ctx, apiName string, message string) error {
	return ErrorResponse(c, apiName, fiber.StatusUnauthorized, models.ErrorCodeUnauthorized, message, nil)
}

// ForbiddenResponse creates a 403 Forbidden response
func ForbiddenResponse(c *fiber.Ctx, apiName string, message string) error {
	return ErrorResponse(c, apiName, fiber.StatusForbidden, models.ErrorCodeForbidden, message, nil)
}

// InternalErrorResponse creates a 500 Internal Server Error response
func InternalErrorResponse(c *fiber.Ctx, apiName string, message string, details interface{}) error {
	return ErrorResponse(c, apiName, fiber.StatusInternalServerError, models.ErrorCodeInternalError, message, details)
}

// ValidationErrorResponse creates a 422 Unprocessable Entity response for validation errors
func ValidationErrorResponse(c *fiber.Ctx, apiName string, message string, details interface{}) error {
	return ErrorResponse(c, apiName, fiber.StatusUnprocessableEntity, models.ErrorCodeValidationFail, message, details)
}

// getEnvOrDefault gets an environment variable or returns a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ValidateJSONContentType validates that the request has a JSON content-type header
// It performs case-insensitive matching on "application/json"
func ValidateJSONContentType(c *fiber.Ctx, apiName string) error {
	contentType := c.Get("Content-Type")

	// Log the content-type header for audit purposes
	Logger.Info("Validating content-type header",
		zap.String("api", apiName),
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("content_type", contentType),
	)

	// Check if content-type is empty
	if contentType == "" {
		Logger.Warn("Missing content-type header",
			zap.String("api", apiName),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		)
		return BadRequestResponse(c, apiName, "Content-Type header is required", map[string]string{
			"required": "application/json",
			"received": "none",
		})
	}

	// Case-insensitive check for application/json
	// Also handle charset parameter (e.g., "application/json; charset=utf-8")
	contentTypeLower := ""
	for i := 0; i < len(contentType); i++ {
		if contentType[i] == ';' {
			break
		}
		if contentType[i] >= 'A' && contentType[i] <= 'Z' {
			contentTypeLower += string(contentType[i] + 32)
		} else {
			contentTypeLower += string(contentType[i])
		}
	}

	// Trim whitespace
	contentTypeLower = trimWhitespace(contentTypeLower)

	if contentTypeLower != "application/json" {
		Logger.Warn("Invalid content-type header",
			zap.String("api", apiName),
			zap.String("request_id", c.Get("X-Request-ID", "unknown")),
			zap.String("content_type", contentType),
			zap.String("expected", "application/json"),
		)
		return BadRequestResponse(c, apiName, "Invalid Content-Type header", map[string]string{
			"required": "application/json",
			"received": contentType,
		})
	}

	Logger.Info("Content-type validation passed",
		zap.String("api", apiName),
		zap.String("request_id", c.Get("X-Request-ID", "unknown")),
		zap.String("content_type", contentType),
	)

	return nil
}

// trimWhitespace removes leading and trailing whitespace from a string
func trimWhitespace(s string) string {
	start := 0
	end := len(s)

	// Trim leading whitespace
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	// Trim trailing whitespace
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
