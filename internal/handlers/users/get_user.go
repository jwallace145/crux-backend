package users

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// GetUser handles GET /users requests to retrieve the authenticated user
// Requires AuthMiddleware to be applied - reads user_id from context
func GetUser(c *fiber.Ctx) error {
	apiName := "get_user"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting get user process",
		zap.String("api", apiName),
	)

	// Get user ID from context (set by AuthMiddleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Error("User ID not found in context",
			zap.String("api", apiName),
		)
		return handlers.InternalErrorResponse(c, apiName, "Authentication context missing", nil)
	}

	log.Info("User ID retrieved from context",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
	)

	// Fetch user from db
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("User not found",
				zap.String("api", apiName),
				zap.Uint("user_id", userID),
			)
			return handlers.UnauthorizedResponse(c, apiName, "User not found")
		}
		log.Error("Database error while looking up user",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", userID),
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
