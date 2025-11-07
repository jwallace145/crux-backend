package users

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	awsClient "github.com/jwallace145/crux-backend/internal/aws"
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

	// Generate presigned URL if profile picture exists
	var response *models.UserResponse
	if user.ProfilePictureURI != "" {
		s3Key := extractS3KeyFromURI(user.ProfilePictureURI)
		if s3Key != "" {
			log.Info("Generating presigned URL for profile picture",
				zap.String("api", apiName),
				zap.String("s3_key", s3Key),
			)

			presignedURL, err := awsClient.GeneratePresignedURL(c.Context(), "crux-project-dev", s3Key, 60)
			if err != nil {
				log.Error("Failed to generate presigned URL",
					zap.Error(err),
					zap.String("api", apiName),
				)
				// Still return success but without presigned URL
				response = user.ToUserResponse()
			} else {
				expiresAt := time.Now().Add(60 * time.Minute).Format("2006-01-02T15:04:05Z07:00")
				response = user.ToUserResponseWithPresignedURL(presignedURL, expiresAt)

				log.Info("Presigned URL generated successfully",
					zap.String("api", apiName),
					zap.String("expires_at", expiresAt),
				)
			}
		} else {
			response = user.ToUserResponse()
		}
	} else {
		response = user.ToUserResponse()
	}

	log.Info("Get user completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return handlers.SuccessResponse(c, apiName, response, "User retrieved successfully")
}
