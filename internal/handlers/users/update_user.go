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

// UpdateUser handles PUT /users requests to update the authenticated user's information
// Currently supports updating FirstName and LastName fields
// Requires AuthMiddleware to be applied - reads user_id from context
func UpdateUser(c *fiber.Ctx) error {
	apiName := "update_user"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting update user process",
		zap.String("api", apiName),
	)

	// Validate Content-Type header
	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		// Error response is already sent by ValidateJSONContentType
		return err
	}

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

	// Parse request body
	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
			zap.ByteString("raw_body", c.Body()),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
	)

	// Validate request fields
	if err := validateUpdateUserRequest(&req); err != nil {
		log.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	log.Info("Request validation passed",
		zap.String("api", apiName),
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

	// Track what fields are being updated
	updates := make(map[string]interface{})

	// Update fields if provided
	if req.FirstName != nil {
		log.Info("Updating first name",
			zap.String("api", apiName),
			zap.String("old_value", user.FirstName),
			zap.String("new_value", *req.FirstName),
		)
		updates["first_name"] = *req.FirstName
		user.FirstName = *req.FirstName
	}

	if req.LastName != nil {
		log.Info("Updating last name",
			zap.String("api", apiName),
			zap.String("old_value", user.LastName),
			zap.String("new_value", *req.LastName),
		)
		updates["last_name"] = *req.LastName
		user.LastName = *req.LastName
	}

	// Check if any fields were provided for update
	if len(updates) == 0 {
		log.Warn("No fields provided for update",
			zap.String("api", apiName),
		)
		return handlers.BadRequestResponse(c, apiName, "No fields provided for update", nil)
	}

	// Save updates to database
	log.Info("Saving user updates to database",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.Any("updates", updates),
	)

	if err := db.DB.Model(&user).Updates(updates).Error; err != nil {
		log.Error("Failed to update user in database",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", user.ID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to update user", nil)
	}

	log.Info("User updated successfully in database",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.Any("updates", updates),
	)

	// Return updated user data
	response := user.ToUserResponse()

	log.Info("Update user completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.Any("response", response),
	)

	return handlers.SuccessResponse(c, apiName, response, "User updated successfully")
}

// validateUpdateUserRequest validates the update user request
func validateUpdateUserRequest(req *models.UpdateUserRequest) error {
	if req.FirstName != nil && len(*req.FirstName) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "First name must not exceed 100 characters")
	}
	if req.LastName != nil && len(*req.LastName) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Last name must not exceed 100 characters")
	}
	return nil
}
