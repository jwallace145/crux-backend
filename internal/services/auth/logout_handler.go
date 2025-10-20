package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// Logout handles POST /logout requests to log out users
// Invalidates the session and clears cookies
func Logout(c *fiber.Ctx) error {
	apiName := "logout"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting logout process",
		zap.String("api", apiName),
	)

	// Get access token from cookie
	accessToken := c.Cookies("access_token", "")

	if accessToken != "" {
		// Validate the access token to get session ID
		claims, err := utils.ValidateAccessToken(accessToken)
		if err == nil {
			log.Info("Access token validated, revoking session",
				zap.String("api", apiName),
				zap.String("session_id", claims.SessionID),
				zap.Uint("user_id", claims.UserID),
			)

			// Revoke the session in database
			if err := db.DB.Model(&models.Session{}).
				Where("session_id = ?", claims.SessionID).
				Update("revoked", true).Error; err != nil {
				log.Error("Failed to revoke session",
					zap.Error(err),
					zap.String("api", apiName),
					zap.String("session_id", claims.SessionID),
				)
				// Don't fail logout even if we can't revoke the session
			} else {
				log.Info("Session revoked successfully",
					zap.String("api", apiName),
					zap.String("session_id", claims.SessionID),
				)
			}
		} else {
			log.Warn("Invalid access token during logout",
				zap.Error(err),
				zap.String("api", apiName),
			)
			// Continue with logout even if token is invalid
		}
	} else {
		log.Info("No access token found during logout",
			zap.String("api", apiName),
		)
	}

	// Clear access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour), // Set to past time
	})

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Expires:  time.Now().Add(-time.Hour), // Set to past time
	})

	log.Info("Cookies cleared successfully",
		zap.String("api", apiName),
	)

	response := &models.LogoutResponse{
		Message: "Logout successful",
	}

	log.Info("Logout completed successfully",
		zap.String("api", apiName),
	)

	return utils.SuccessResponse(c, apiName, response, "Logout successful")
}
