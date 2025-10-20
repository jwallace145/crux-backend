package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// Refresh handles POST /refresh requests to refresh access tokens
// Uses the refresh token from HTTP-only cookie to generate a new access token
func Refresh(c *fiber.Ctx) error {
	apiName := "refresh"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting token refresh process",
		zap.String("api", apiName),
	)

	// Get refresh token from cookie
	refreshToken := c.Cookies("refresh_token", "")
	if refreshToken == "" {
		log.Warn("No refresh token found in cookies",
			zap.String("api", apiName),
		)
		return utils.UnauthorizedResponse(c, apiName, "No refresh token provided")
	}

	log.Info("Refresh token found, validating",
		zap.String("api", apiName),
	)

	// Validate refresh token
	claims, err := utils.ValidateRefreshToken(refreshToken)
	if err != nil {
		log.Warn("Invalid refresh token",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return utils.UnauthorizedResponse(c, apiName, "Invalid or expired refresh token")
	}

	log.Info("Refresh token validated successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", claims.UserID),
		zap.String("session_id", claims.SessionID),
	)

	// Check if session is still valid and not revoked
	var session models.Session
	if err := db.DB.Where("session_id = ?", claims.SessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Session not found",
				zap.String("api", apiName),
				zap.String("session_id", claims.SessionID),
			)
			return utils.UnauthorizedResponse(c, apiName, "Session not found")
		}
		log.Error("Database error while checking session",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to validate session", nil)
	}

	// Check if session is revoked
	if session.Revoked {
		log.Warn("Session is revoked",
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
		)
		return utils.UnauthorizedResponse(c, apiName, "Session has been revoked")
	}

	// Check if session has expired
	if time.Now().After(session.ExpiresAt) {
		log.Warn("Session has expired",
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
			zap.Time("expires_at", session.ExpiresAt),
		)
		return utils.UnauthorizedResponse(c, apiName, "Session has expired")
	}

	log.Info("Session is valid and active",
		zap.String("api", apiName),
		zap.String("session_id", claims.SessionID),
	)

	// Get user information
	var user models.User
	if err := db.DB.First(&user, claims.UserID).Error; err != nil {
		log.Error("Failed to find user",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", claims.UserID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to find user", nil)
	}

	log.Info("User found, generating new access token",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	// Generate new access token
	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email, claims.SessionID)
	if err != nil {
		log.Error("Failed to generate new access token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", user.ID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to generate access token", nil)
	}

	log.Info("New access token generated successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
	)

	// Set new access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		Path:     "/",
		MaxAge:   int(utils.AccessTokenExpiry.Seconds()),
		HTTPOnly: true,
	})

	log.Info("New access token cookie set successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
	)

	// Prepare response
	response := &models.RefreshResponse{
		Message:   "Access token refreshed successfully",
		ExpiresAt: time.Now().Add(utils.AccessTokenExpiry).Format(time.RFC3339),
	}

	log.Info("Token refresh completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("session_id", claims.SessionID),
	)

	return utils.SuccessResponse(c, apiName, response, "Access token refreshed successfully")
}
