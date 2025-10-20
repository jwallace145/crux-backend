package users

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// CreateUser handles POST /users requests to create a new user
// It validates the request, checks for existing users with the same email/username,
// and persists the new user to the database
func CreateUser(c *fiber.Ctx) error {
	apiName := "create_user"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting user creation process",
		zap.String("api", apiName),
	)

	// Validate Content-Type header
	if err := utils.ValidateJSONContentType(c, apiName); err != nil {
		// Error response is already sent by ValidateJSONContentType
		return err
	}

	// Parse request body
	var req models.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
			zap.ByteString("raw_body", c.Body()),
		)
		return utils.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
		zap.String("username", req.Username),
		zap.String("email", req.Email),
		zap.String("first_name", req.FirstName),
		zap.String("last_name", req.LastName),
	)

	// Validate required fields
	log.Info("Validating request parameters",
		zap.String("api", apiName),
	)

	if err := validateCreateUserRequest(&req); err != nil {
		log.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("username", req.Username),
			zap.String("email", req.Email),
		)
		return utils.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	log.Info("Request validation passed",
		zap.String("api", apiName),
	)

	// Normalize email to lowercase for consistency
	originalEmail := req.Email
	originalUsername := req.Username
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Username = strings.TrimSpace(req.Username)

	log.Info("Normalized request parameters",
		zap.String("api", apiName),
		zap.String("original_email", originalEmail),
		zap.String("normalized_email", req.Email),
		zap.String("original_username", originalUsername),
		zap.String("normalized_username", req.Username),
	)

	// Check if user already exists by email
	log.Info("Checking for existing user by email",
		zap.String("api", apiName),
		zap.String("email", req.Email),
	)

	existingUser, err := getUserByEmail(req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Database error while checking for existing user by email",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("email", req.Email),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to check for existing user", nil)
	}

	if existingUser != nil {
		log.Warn("User creation rejected - email already exists",
			zap.String("api", apiName),
			zap.String("email", req.Email),
			zap.Uint("existing_user_id", existingUser.ID),
		)
		return utils.BadRequestResponse(c, apiName, "User with this email already exists", map[string]string{
			"field": "email",
			"value": req.Email,
		})
	}

	log.Info("Email uniqueness check passed",
		zap.String("api", apiName),
		zap.String("email", req.Email),
	)

	// Check if username is already taken
	log.Info("Checking for existing user by username",
		zap.String("api", apiName),
		zap.String("username", req.Username),
	)

	existingUsername, err := getUserByUsername(req.Username)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("Database error while checking for existing username",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("username", req.Username),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to check for existing username", nil)
	}

	if existingUsername != nil {
		log.Warn("User creation rejected - username already exists",
			zap.String("api", apiName),
			zap.String("username", req.Username),
			zap.Uint("existing_user_id", existingUsername.ID),
		)
		return utils.BadRequestResponse(c, apiName, "User with this username already exists", map[string]string{
			"field": "username",
			"value": req.Username,
		})
	}

	log.Info("Username uniqueness check passed",
		zap.String("api", apiName),
		zap.String("username", req.Username),
	)

	// Hash the password
	log.Info("Hashing user password",
		zap.String("api", apiName),
		zap.String("username", req.Username),
	)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("Failed to hash password",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("username", req.Username),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to process password", nil)
	}

	log.Info("Password hashed successfully",
		zap.String("api", apiName),
		zap.String("username", req.Username),
	)

	// Create new user
	log.Info("Creating new user in database",
		zap.String("api", apiName),
		zap.String("username", req.Username),
		zap.String("email", req.Email),
		zap.String("first_name", req.FirstName),
		zap.String("last_name", req.LastName),
	)

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := db.DB.Create(user).Error; err != nil {
		log.Error("Failed to create user in database",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("email", req.Email),
			zap.String("username", req.Username),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to create user", nil)
	}

	log.Info("User created successfully in database",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
		zap.String("username", user.Username),
		zap.Time("created_at", user.CreatedAt),
	)

	// Prepare response
	response := user.ToUserResponse()

	log.Info("User creation completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("email", user.Email),
		zap.Any("response", response),
	)

	// Return created user
	return utils.CreatedResponse(c, apiName, response, "User created successfully")
}

// validateCreateUserRequest validates the create user request
func validateCreateUserRequest(req *models.CreateUserRequest) error {
	if req.Username == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username is required")
	}
	if len(req.Username) < 3 {
		return fiber.NewError(fiber.StatusBadRequest, "Username must be at least 3 characters")
	}
	if len(req.Username) > 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Username must not exceed 50 characters")
	}
	if req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is required")
	}
	if len(req.Email) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Email must not exceed 100 characters")
	}
	// Basic email validation
	if !strings.Contains(req.Email, "@") || !strings.Contains(req.Email, ".") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}
	if req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Password is required")
	}
	if len(req.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters")
	}
	if len(req.Password) > 72 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must not exceed 72 characters")
	}
	if len(req.FirstName) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "First name must not exceed 100 characters")
	}
	if len(req.LastName) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Last name must not exceed 100 characters")
	}
	return nil
}

// getUserByEmail retrieves a user by email address
func getUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// getUserByUsername retrieves a user by username
func getUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := db.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
