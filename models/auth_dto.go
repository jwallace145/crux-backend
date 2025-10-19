package models

// LoginRequest represents the request body for user login
// Supports login with either username or email
type LoginRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password" validate:"required"`
}

// LoginResponse represents the response for successful login
type LoginResponse struct {
	User      *UserResponse `json:"user"`
	SessionID string        `json:"session_id"`
	ExpiresAt string        `json:"expires_at"`
	Message   string        `json:"message"`
}

// RefreshRequest represents the request body for token refresh
// The refresh token will be sent via HTTP-only cookie
type RefreshRequest struct {
	// No body needed - token comes from cookie
}

// RefreshResponse represents the response for successful token refresh
type RefreshResponse struct {
	Message   string `json:"message"`
	ExpiresAt string `json:"expires_at"`
}

// LogoutResponse represents the response for successful logout
type LogoutResponse struct {
	Message string `json:"message"`
}
