package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/services"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// AuthMiddleware validates the access token from cookies and automatically
// refreshes it using the refresh token if expired.
//
// The middleware follows this flow:
// 1. Check for access_token in cookies
// 2. If valid, set user info in context and proceed
// 3. If expired but refresh_token is present:
//   - Validate refresh token
//   - Check session validity
//   - Generate new access token
//   - Set new access_token cookie
//   - Set user info in context and proceed
//
// 4. If no valid tokens, return 401 Unauthorized
//
// The middleware stores the following in context for use by handlers:
// - c.Locals("user_id") - The authenticated user's ID
// - c.Locals("username") - The authenticated user's username
// - c.Locals("email") - The authenticated user's email
// - c.Locals("session_id") - The session ID
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		log := utils.GetLoggerFromContext(c)

		log.Info("Authenticating request",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
		)

		// Get access token from cookie
		accessToken := c.Cookies("access_token", "")

		// Try to validate access token
		claims, err := services.ValidateAccessToken(accessToken)
		if err == nil {
			// Access token is valid
			log.Info("Access token is valid",
				zap.Uint("user_id", claims.UserID),
				zap.String("username", claims.Username),
			)

			// Store user info in context
			c.Locals("user_id", claims.UserID)
			c.Locals("username", claims.Username)
			c.Locals("email", claims.Email)
			c.Locals("session_id", claims.SessionID)

			return c.Next()
		}

		log.Info("Access token invalid or expired, attempting refresh",
			zap.Error(err),
		)

		// Access token is invalid or expired, try to refresh it
		refreshToken := c.Cookies("refresh_token", "")
		if refreshToken == "" {
			log.Warn("No refresh token found",
				zap.String("path", c.Path()),
			)
			return handlers.UnauthorizedResponse(c, "auth_middleware", "Not authenticated")
		}

		// Validate refresh token
		refreshClaims, err := services.ValidateRefreshToken(refreshToken)
		if err != nil {
			log.Warn("Invalid refresh token",
				zap.Error(err),
				zap.String("path", c.Path()),
			)
			return handlers.UnauthorizedResponse(c, "auth_middleware", "Invalid or expired refresh token")
		}

		log.Info("Refresh token validated successfully",
			zap.Uint("user_id", refreshClaims.UserID),
			zap.String("session_id", refreshClaims.SessionID),
		)

		// Verify session exists and is valid
		var session models.Session
		if err := db.DB.Where("session_id = ?", refreshClaims.SessionID).First(&session).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warn("Session not found",
					zap.String("session_id", refreshClaims.SessionID),
				)
				return handlers.UnauthorizedResponse(c, "auth_middleware", "Session not found")
			}
			log.Error("Database error while checking session",
				zap.Error(err),
				zap.String("session_id", refreshClaims.SessionID),
			)
			return handlers.InternalErrorResponse(c, "auth_middleware", "Failed to validate session", nil)
		}

		// Check if session is revoked
		if session.Revoked {
			log.Warn("Session is revoked",
				zap.String("session_id", refreshClaims.SessionID),
			)
			return handlers.UnauthorizedResponse(c, "auth_middleware", "Session has been revoked")
		}

		// Check if session has expired
		if time.Now().After(session.ExpiresAt) {
			log.Warn("Session has expired",
				zap.String("session_id", refreshClaims.SessionID),
				zap.Time("expires_at", session.ExpiresAt),
			)
			return handlers.UnauthorizedResponse(c, "auth_middleware", "Session has expired")
		}

		log.Info("Session is valid, generating new access token",
			zap.Uint("user_id", refreshClaims.UserID),
			zap.String("session_id", refreshClaims.SessionID),
		)

		// Generate new access token
		newAccessToken, err := services.GenerateAccessToken(
			refreshClaims.UserID,
			refreshClaims.Username,
			refreshClaims.Email,
			refreshClaims.SessionID,
		)
		if err != nil {
			log.Error("Failed to generate new access token",
				zap.Error(err),
				zap.Uint("user_id", refreshClaims.UserID),
			)
			return handlers.InternalErrorResponse(c, "auth_middleware", "Failed to generate access token", nil)
		}

		log.Info("New access token generated successfully",
			zap.Uint("user_id", refreshClaims.UserID),
		)

		// Set new access token cookie
		c.Cookie(&fiber.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			Path:     "/",
			MaxAge:   int(services.AccessTokenExpiry.Seconds()),
			HTTPOnly: true,
			Secure:   true,
			SameSite: "None",
		})

		log.Info("New access token cookie set successfully, token refreshed automatically",
			zap.Uint("user_id", refreshClaims.UserID),
			zap.String("path", c.Path()),
		)

		// Store user info in context
		c.Locals("user_id", refreshClaims.UserID)
		c.Locals("username", refreshClaims.Username)
		c.Locals("email", refreshClaims.Email)
		c.Locals("session_id", refreshClaims.SessionID)

		return c.Next()
	}
}
