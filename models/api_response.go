package models

import (
	"time"
)

// APIResponse is the standardized response structure for all API endpoints.
// It provides consistent metadata about the service and request along with
// endpoint-specific data in the Data field.
type APIResponse struct {
	// ServiceName identifies the API service (e.g., "crux-backend")
	ServiceName string `json:"service_name"`

	// Version is the API version (e.g., "1.0.0", "v1")
	Version string `json:"version"`

	// Environment indicates the deployment environment (e.g., "development", "production")
	Environment string `json:"environment"`

	// APIName identifies the specific endpoint (e.g., "get_user", "create_route")
	APIName string `json:"api_name"`

	// Timestamp is when the response was generated (RFC3339 format)
	Timestamp time.Time `json:"timestamp"`

	// Status indicates the response status (e.g., "success", "error")
	Status string `json:"status"`

	// Message provides a human-readable description of the response
	Message string `json:"message,omitempty"`

	// Data contains endpoint-specific response data (optional)
	Data interface{} `json:"data,omitempty"`

	// Error contains error details if Status is "error" (optional)
	Error *APIError `json:"error,omitempty"`
}

// APIError provides structured error information
type APIError struct {
	// Code is the error code (e.g., "INVALID_INPUT", "NOT_FOUND")
	Code string `json:"code"`

	// Message is a human-readable error message
	Message string `json:"message"`

	// Details provides additional error context (optional)
	Details interface{} `json:"details,omitempty"`
}

// ResponseStatus constants for common response statuses
const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// Common error codes
const (
	ErrorCodeInvalidInput   = "INVALID_INPUT"
	ErrorCodeNotFound       = "NOT_FOUND"
	ErrorCodeUnauthorized   = "UNAUTHORIZED"
	ErrorCodeForbidden      = "FORBIDDEN"
	ErrorCodeInternalError  = "INTERNAL_ERROR"
	ErrorCodeDatabaseError  = "DATABASE_ERROR"
	ErrorCodeValidationFail = "VALIDATION_FAILED"
)
