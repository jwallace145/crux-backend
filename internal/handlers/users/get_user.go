package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/services"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// GetUser handles GET /users requests to retrieve the authenticated user
// Validates the access_token from cookies and returns the user data
// Returns 401 if token is invalid, expired, or session is revoked
func GetUser(c *fiber.Ctx) error {
	apiName := "get_user"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting get user process",
		zap.String("api", apiName),
	)

	// Get access token from cookie
	accessToken := c.Cookies("access_token")
	if accessToken == "" {
		log.Warn("Missing access token cookie",
			zap.String("api", apiName),
		)
		return handlers.UnauthorizedResponse(c, apiName, "Not authenticated")
	}

	log.Info("Access token found in cookie",
		zap.String("api", apiName),
	)

	// Validate access token
	claims, err := services.ValidateAccessToken(accessToken)
	if err != nil {
		log.Warn("Invalid access token",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.UnauthorizedResponse(c, apiName, "Invalid or expired token")
	}

	log.Info("Access token validated successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", claims.UserID),
		zap.String("session_id", claims.SessionID),
	)

	// Verify session exists and is valid
	var session models.Session
	if err := db.DB.Where("session_id = ?", claims.SessionID).First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Session not found",
				zap.String("api", apiName),
				zap.String("session_id", claims.SessionID),
			)
			return handlers.UnauthorizedResponse(c, apiName, "Session not found")
		}
		log.Error("Database error while looking up session",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to validate session", nil)
	}

	log.Info("Session found",
		zap.String("api", apiName),
		zap.Uint("session_id", session.ID),
		zap.String("session_uuid", claims.SessionID),
	)

	// Check if session is revoked
	if session.Revoked {
		log.Warn("Session is revoked",
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
		)
		return handlers.UnauthorizedResponse(c, apiName, "Session has been revoked")
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		log.Warn("Session is expired",
			zap.String("api", apiName),
			zap.String("session_id", claims.SessionID),
			zap.Time("expires_at", session.ExpiresAt),
		)
		return handlers.UnauthorizedResponse(c, apiName, "Session has expired")
	}

	log.Info("Session is valid",
		zap.String("api", apiName),
		zap.String("session_id", claims.SessionID),
	)

	// Fetch user from db
	var user models.User
	if err := db.DB.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("User not found",
				zap.String("api", apiName),
				zap.Uint("user_id", claims.UserID),
			)
			return handlers.UnauthorizedResponse(c, apiName, "User not found")
		}
		log.Error("Database error while looking up user",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", claims.UserID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to fetch user", nil)
	}

	log.Info("User fetched successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	// Return user data
	response := user.ToUserResponse()

	log.Info("Get user completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return handlers.SuccessResponse(c, apiName, response, "User retrieved successfully")
}
