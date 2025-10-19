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
	requestID := c.Get("X-Request-ID", "unknown")

	utils.Logger.Info("Starting logout process",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	// Get access token from cookie
	accessToken := c.Cookies("access_token", "")

	if accessToken != "" {
		// Validate the access token to get session ID
		claims, err := utils.ValidateAccessToken(accessToken)
		if err == nil {
			utils.Logger.Info("Access token validated, revoking session",
				zap.String("api", apiName),
				zap.String("request_id", requestID),
				zap.String("session_id", claims.SessionID),
				zap.Uint("user_id", claims.UserID),
			)

			// Revoke the session in database
			if err := db.DB.Model(&models.Session{}).
				Where("session_id = ?", claims.SessionID).
				Update("revoked", true).Error; err != nil {
				utils.Logger.Error("Failed to revoke session",
					zap.Error(err),
					zap.String("api", apiName),
					zap.String("request_id", requestID),
					zap.String("session_id", claims.SessionID),
				)
				// Don't fail logout even if we can't revoke the session
			} else {
				utils.Logger.Info("Session revoked successfully",
					zap.String("api", apiName),
					zap.String("request_id", requestID),
					zap.String("session_id", claims.SessionID),
				)
			}
		} else {
			utils.Logger.Warn("Invalid access token during logout",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("request_id", requestID),
			)
			// Continue with logout even if token is invalid
		}
	} else {
		utils.Logger.Info("No access token found during logout",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
	}

	// Clear access token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour), // Set to past time
	})

	// Clear refresh token cookie
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		HTTPOnly: true,
		Secure:   true,
		SameSite: "Strict",
		Expires:  time.Now().Add(-time.Hour), // Set to past time
	})

	utils.Logger.Info("Cookies cleared successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	response := &models.LogoutResponse{
		Message: "Logout successful",
	}

	utils.Logger.Info("Logout completed successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	return utils.SuccessResponse(c, apiName, response, "Logout successful")
}
