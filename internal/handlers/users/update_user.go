package users

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	awsClient "github.com/jwallace145/crux-backend/internal/aws"
	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// UpdateUser handles PUT /users requests to update the authenticated user's information
// Supports updating Username, Email, FirstName, LastName, and profile picture
// Requires AuthMiddleware to be applied - reads user_id from context
// Accepts both JSON and multipart/form-data for profile picture uploads
func UpdateUser(c *fiber.Ctx) error {
	apiName := "update_user"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting update user process",
		zap.String("api", apiName),
	)

	// Get user ID from context (set by AuthMiddleware)
	userID, err := getUserIDFromContext(c, apiName)
	if err != nil {
		return err
	}

	log.Info("User ID retrieved from context",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
	)

	// Fetch user from db
	user, err := fetchUserByID(c, apiName, userID)
	if err != nil {
		return err
	}

	log.Info("User fetched successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	// Track what fields are being updated
	updates := make(map[string]interface{})

	// Check content type to determine how to parse the request
	contentType := c.Get("Content-Type")
	isMultipart := strings.Contains(contentType, "multipart/form-data")

	if isMultipart {
		if err := processMultipartRequest(c, apiName, user, userID, updates); err != nil {
			return err
		}
	} else {
		if err := processJSONRequest(c, apiName, user, userID, updates); err != nil {
			return err
		}
	}

	// Check if any fields were provided for update
	if len(updates) == 0 {
		log.Warn("No fields provided for update",
			zap.String("api", apiName),
		)
		return handlers.BadRequestResponse(c, apiName, "No fields provided for update", nil)
	}

	// Save updates to database
	if err := saveUserUpdates(c, apiName, user, updates); err != nil {
		return err
	}

	// Generate presigned URL if profile picture exists
	response := generateUserResponse(c, apiName, user)

	log.Info("Update user completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	return handlers.SuccessResponse(c, apiName, response, "User updated successfully")
}

// getUserIDFromContext retrieves the user ID from the request context
func getUserIDFromContext(c *fiber.Ctx, apiName string) (uint, error) {
	log := utils.GetLoggerFromContext(c)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		log.Error("User ID not found in context",
			zap.String("api", apiName),
		)
		return 0, handlers.InternalErrorResponse(c, apiName, "Authentication context missing", nil)
	}
	return userID, nil
}

// fetchUserByID fetches a user from the database by ID
func fetchUserByID(c *fiber.Ctx, apiName string, userID uint) (*models.User, error) {
	log := utils.GetLoggerFromContext(c)
	var user models.User
	if err := db.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("User not found",
				zap.String("api", apiName),
				zap.Uint("user_id", userID),
			)
			return nil, handlers.UnauthorizedResponse(c, apiName, "User not found")
		}
		log.Error("Database error while looking up user",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", userID),
		)
		return nil, handlers.InternalErrorResponse(c, apiName, "Failed to fetch user", nil)
	}
	return &user, nil
}

// processMultipartRequest handles multipart/form-data requests
func processMultipartRequest(c *fiber.Ctx, apiName string, user *models.User, userID uint, updates map[string]interface{}) error {
	log := utils.GetLoggerFromContext(c)
	log.Info("Processing multipart/form-data request",
		zap.String("api", apiName),
	)

	// Parse form fields
	if err := processUsernameField(c, apiName, c.FormValue("username"), userID, user, updates); err != nil {
		return err
	}

	if err := processEmailField(c, apiName, c.FormValue("email"), userID, user, updates); err != nil {
		return err
	}

	if err := processFirstNameField(c, apiName, c.FormValue("first_name"), user, updates); err != nil {
		return err
	}

	if err := processLastNameField(c, apiName, c.FormValue("last_name"), user, updates); err != nil {
		return err
	}

	// Handle profile picture upload
	return processProfilePictureUpload(c, apiName, user, userID, updates)
}

// processJSONRequest handles JSON requests
func processJSONRequest(c *fiber.Ctx, apiName string, user *models.User, userID uint, updates map[string]interface{}) error {
	log := utils.GetLoggerFromContext(c)

	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
	)

	// Validate and update fields
	if req.Username != nil {
		if err := processUsernameField(c, apiName, *req.Username, userID, user, updates); err != nil {
			return err
		}
	}

	if req.Email != nil {
		if err := processEmailField(c, apiName, *req.Email, userID, user, updates); err != nil {
			return err
		}
	}

	if req.FirstName != nil {
		if err := processFirstNameField(c, apiName, *req.FirstName, user, updates); err != nil {
			return err
		}
	}

	if req.LastName != nil {
		if err := processLastNameField(c, apiName, *req.LastName, user, updates); err != nil {
			return err
		}
	}

	return nil
}

// processUsernameField validates and processes username updates
func processUsernameField(c *fiber.Ctx, apiName, username string, userID uint, user *models.User, updates map[string]interface{}) error {
	if username == "" {
		return nil
	}
	if err := validateUsername(username); err != nil {
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}
	if err := checkUsernameAvailability(username, userID); err != nil {
		return handlers.BadRequestResponse(c, apiName, err.Error(), nil)
	}
	updates["username"] = username
	user.Username = username
	return nil
}

// processEmailField validates and processes email updates
func processEmailField(c *fiber.Ctx, apiName, email string, userID uint, user *models.User, updates map[string]interface{}) error {
	if email == "" {
		return nil
	}
	email = strings.ToLower(strings.TrimSpace(email))
	if err := validateEmail(email); err != nil {
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}
	if err := checkEmailAvailability(email, userID); err != nil {
		return handlers.BadRequestResponse(c, apiName, err.Error(), nil)
	}
	updates["email"] = email
	user.Email = email
	return nil
}

// processFirstNameField validates and processes first name updates
func processFirstNameField(c *fiber.Ctx, apiName, firstName string, user *models.User, updates map[string]interface{}) error {
	if firstName == "" {
		return nil
	}
	if len(firstName) > 100 {
		return handlers.ValidationErrorResponse(c, apiName, "First name must not exceed 100 characters", nil)
	}
	updates["first_name"] = firstName
	user.FirstName = firstName
	return nil
}

// processLastNameField validates and processes last name updates
func processLastNameField(c *fiber.Ctx, apiName, lastName string, user *models.User, updates map[string]interface{}) error {
	if lastName == "" {
		return nil
	}
	if len(lastName) > 100 {
		return handlers.ValidationErrorResponse(c, apiName, "Last name must not exceed 100 characters", nil)
	}
	updates["last_name"] = lastName
	user.LastName = lastName
	return nil
}

// processProfilePictureUpload handles profile picture file upload and S3 storage
func processProfilePictureUpload(c *fiber.Ctx, apiName string, user *models.User, userID uint, updates map[string]interface{}) error {
	log := utils.GetLoggerFromContext(c)

	fileHeader, err := c.FormFile("profile_picture")
	if err != nil {
		// No file uploaded, not an error
		return nil
	}

	log.Info("Profile picture file detected",
		zap.String("api", apiName),
		zap.String("filename", fileHeader.Filename),
		zap.Int64("size", fileHeader.Size),
	)

	// Validate file size (5MB max)
	if err := validateFileSize(fileHeader); err != nil {
		return handlers.BadRequestResponse(c, apiName, err.Error(), nil)
	}

	// Validate file type
	contentType := fileHeader.Header.Get("Content-Type")
	if !isValidImageType(contentType) {
		return handlers.BadRequestResponse(c, apiName, "Profile picture must be a JPEG, PNG, or GIF image", nil)
	}

	// Read file content
	fileBytes, err := readUploadedFile(fileHeader)
	if err != nil {
		log.Error("Failed to process uploaded file",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to process file", nil)
	}

	// Delete old profile picture if exists
	deleteOldProfilePicture(c, apiName, user.ProfilePictureURI)

	// Upload to S3
	s3URI, err := uploadProfilePictureToS3(c, apiName, userID, fileBytes, contentType)
	if err != nil {
		return err
	}

	updates["profile_picture_uri"] = s3URI
	user.ProfilePictureURI = s3URI

	log.Info("Profile picture uploaded successfully",
		zap.String("api", apiName),
		zap.String("s3_uri", s3URI),
	)

	return nil
}

// validateFileSize validates the uploaded file size
func validateFileSize(fileHeader *multipart.FileHeader) error {
	const maxFileSize = 5 * 1024 * 1024
	if fileHeader.Size > maxFileSize {
		return fiber.NewError(fiber.StatusBadRequest, "Profile picture must not exceed 5MB")
	}
	return nil
}

// readUploadedFile reads the content of an uploaded file
func readUploadedFile(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

// deleteOldProfilePicture deletes the old profile picture from S3 if it exists
func deleteOldProfilePicture(c *fiber.Ctx, apiName, profilePictureURI string) {
	log := utils.GetLoggerFromContext(c)

	if profilePictureURI == "" {
		return
	}

	oldKey := extractS3KeyFromURI(profilePictureURI)
	if oldKey == "" {
		return
	}

	log.Info("Deleting old profile picture",
		zap.String("api", apiName),
		zap.String("old_uri", profilePictureURI),
	)

	if err := awsClient.DeleteFile(c.Context(), "crux-project-dev", oldKey); err != nil {
		log.Warn("Failed to delete old profile picture",
			zap.Error(err),
			zap.String("api", apiName),
		)
	}
}

// uploadProfilePictureToS3 uploads a profile picture to S3 and returns the S3 URI
func uploadProfilePictureToS3(c *fiber.Ctx, apiName string, userID uint, fileBytes []byte, contentType string) (string, error) {
	log := utils.GetLoggerFromContext(c)

	pictureID := uuid.New().String()
	s3Key := fmt.Sprintf("users/%d/profile_picture/%s", userID, pictureID)

	log.Info("Uploading profile picture to S3",
		zap.String("api", apiName),
		zap.String("s3_key", s3Key),
	)

	s3URI, err := awsClient.UploadFile(c.Context(), "crux-project-dev", s3Key, bytes.NewReader(fileBytes), contentType)
	if err != nil {
		log.Error("Failed to upload profile picture to S3",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return "", handlers.InternalErrorResponse(c, apiName, "Failed to upload profile picture", nil)
	}

	return s3URI, nil
}

// saveUserUpdates saves the user updates to the database
func saveUserUpdates(c *fiber.Ctx, apiName string, user *models.User, updates map[string]interface{}) error {
	log := utils.GetLoggerFromContext(c)

	// Create a filtered version of updates for logging (exclude profile picture URI)
	loggableUpdates := filterUpdatesForLogging(updates)

	log.Info("Saving user updates to database",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.Any("updates", loggableUpdates),
	)

	if err := db.DB.Model(user).Updates(updates).Error; err != nil {
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
		zap.Any("updates", loggableUpdates),
	)

	return nil
}

// generateUserResponse generates the user response with presigned URL if applicable
func generateUserResponse(c *fiber.Ctx, apiName string, user *models.User) *models.UserResponse {
	log := utils.GetLoggerFromContext(c)

	if user.ProfilePictureURI == "" {
		return user.ToUserResponse()
	}

	s3Key := extractS3KeyFromURI(user.ProfilePictureURI)
	if s3Key == "" {
		return user.ToUserResponse()
	}

	presignedURL, err := awsClient.GeneratePresignedURL(c.Context(), "crux-project-dev", s3Key, 60)
	if err != nil {
		log.Error("Failed to generate presigned URL",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return user.ToUserResponse()
	}

	expiresAt := time.Now().Add(60 * time.Minute).Format("2006-01-02T15:04:05Z07:00")
	return user.ToUserResponseWithPresignedURL(presignedURL, expiresAt)
}

// validateUsername validates a username
func validateUsername(username string) error {
	if len(username) < 3 {
		return fiber.NewError(fiber.StatusBadRequest, "Username must be at least 3 characters")
	}
	if len(username) > 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Username must not exceed 50 characters")
	}
	return nil
}

// validateEmail validates an email address
func validateEmail(email string) error {
	if len(email) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Email must not exceed 100 characters")
	}
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}
	return nil
}

// checkUsernameAvailability checks if a username is available for the user
func checkUsernameAvailability(username string, currentUserID uint) error {
	var existingUser models.User
	err := db.DB.Where("username = ? AND id != ?", username, currentUserID).First(&existingUser).Error
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Username is already taken")
	}
	if err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check username availability")
	}
	return nil
}

// checkEmailAvailability checks if an email is available for the user
func checkEmailAvailability(email string, currentUserID uint) error {
	var existingUser models.User
	err := db.DB.Where("email = ? AND id != ?", email, currentUserID).First(&existingUser).Error
	if err == nil {
		return fiber.NewError(fiber.StatusBadRequest, "Email is already taken")
	}
	if err != gorm.ErrRecordNotFound {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check email availability")
	}
	return nil
}

// isValidImageType checks if the content type is a valid image type
func isValidImageType(contentType string) bool {
	validTypes := []string{"image/jpeg", "image/png", "image/gif", "image/jpg"}
	for _, validType := range validTypes {
		if contentType == validType {
			return true
		}
	}
	return false
}

// extractS3KeyFromURI extracts the S3 key from an S3 URI (s3://bucket/key)
func extractS3KeyFromURI(s3URI string) string {
	if !strings.HasPrefix(s3URI, "s3://") {
		return ""
	}
	// Remove s3:// prefix
	withoutPrefix := strings.TrimPrefix(s3URI, "s3://")
	// Split by first /
	parts := strings.SplitN(withoutPrefix, "/", 2)
	if len(parts) != 2 {
		return ""
	}
	return parts[1]
}

// filterUpdatesForLogging creates a copy of the updates map without profile picture URI
// to keep logs clean and avoid logging large S3 URIs
func filterUpdatesForLogging(updates map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})
	for key, value := range updates {
		if key == "profile_picture_uri" {
			filtered[key] = "[profile_picture_updated]"
		} else {
			filtered[key] = value
		}
	}
	return filtered
}
