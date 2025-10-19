package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

var (
	// JWT secrets loaded from environment - REQUIRED
	AccessTokenSecret  = getRequiredEnv("ACCESS_TOKEN_SECRET_KEY")
	RefreshTokenSecret = getRequiredEnv("REFRESH_TOKEN_SECRET_KEY")

	// Token expiration times
	AccessTokenExpiry  = 15 * time.Minute   // Short-lived access token
	RefreshTokenExpiry = 7 * 24 * time.Hour // 7 days for refresh token
)

// getRequiredEnv retrieves a required environment variable.
// If the variable is not set or empty, the application will terminate with a fatal error.
// This prevents the application from running with insecure default values.
func getRequiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		Logger.Fatal("Required environment variable not set",
			zap.String("variable", key),
			zap.String("error", "environment variable is required for security"),
		)
	}
	return value
}

// TokenClaims represents the JWT claims structure
type TokenClaims struct {
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	SessionID string `json:"session_id"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a new access token for the user
func GenerateAccessToken(userID uint, username, email, sessionID string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		SessionID: sessionID,
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "crux-backend",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AccessTokenSecret))
}

// GenerateRefreshToken creates a new refresh token for the user
func GenerateRefreshToken(userID uint, username, email, sessionID string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserID:    userID,
		Username:  username,
		Email:     email,
		SessionID: sessionID,
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(RefreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "crux-backend",
			Subject:   username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(RefreshTokenSecret))
}

// ValidateAccessToken validates and parses an access token
func ValidateAccessToken(tokenString string) (*TokenClaims, error) {
	return validateToken(tokenString, AccessTokenSecret, "access")
}

// ValidateRefreshToken validates and parses a refresh token
func ValidateRefreshToken(tokenString string) (*TokenClaims, error) {
	return validateToken(tokenString, RefreshTokenSecret, "refresh")
}

// validateToken is a helper function to validate tokens
func validateToken(tokenString, secret, expectedType string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Verify token type
	if claims.TokenType != expectedType {
		return nil, errors.New("invalid token type")
	}

	return claims, nil
}
