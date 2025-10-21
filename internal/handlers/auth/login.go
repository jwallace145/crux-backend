package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/services"
	"github.com/jwallace145/crux-backend/internal/utils"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/models"
)

var (
	SessionExpiry = 7 * 24 * time.Hour
)

// Login handles POST /login requests to authenticate users
// Accepts either username+password or email+password
// Returns JWT tokens via secure HTTP-only cookies
func Login(c *fiber.Ctx) error {
	apiName := "login"
	log := utils.GetLoggerFromContext(c).With(zap.String("api", apiName))
	DB := db.GetDB()

	log.Info("Executing login API handler")

	// Validate Content-Type header
	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	// Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// Validate request
	if err := validateLoginRequest(&req); err != nil {
		log.Warn("Request validation failed",
			zap.Error(err),
		)
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	// Normalize inputs
	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}
	if req.Username != "" {
		req.Username = strings.TrimSpace(req.Username)
	}

	// Find user by email or username
	log.Info("Looking up user",
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	var user models.User
	var err error

	if req.Email != "" {
		err = DB.Where("email = ?", req.Email).First(&user).Error
	} else {
		err = DB.Where("username = ?", req.Username).First(&user).Error
	}

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Warn("User not found",
				zap.String("username", req.Username),
				zap.String("email", req.Email),
			)
			return handlers.UnauthorizedResponse(c, apiName, "Invalid credentials")
		}
		log.Error("Database error while looking up user",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.InternalErrorResponse(c, apiName, "Authentication failed", nil)
	}

	log.Info("User found, verifying password",
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn("Invalid password",
			zap.Uint("user_id", user.ID),
			zap.String("username", user.Username),
		)
		return handlers.UnauthorizedResponse(c, apiName, "Invalid credentials")
	}

	log.Info("Password verified successfully",
		zap.Uint("user_id", user.ID),
	)

	// Generate session ID
	sessionID := uuid.New().String()

	log.Info("Generated session ID",
		zap.String("session_id", sessionID),
		zap.Uint("user_id", user.ID),
	)

	// Create session in db
	expiresAt := time.Now().Add(SessionExpiry)
	session := &models.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent", "unknown"),
		Device:    c.Get("User-Agent", "unknown"), // Could be enhanced with device detection
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := DB.Create(session).Error; err != nil {
		log.Error("Failed to create session",
			zap.Error(err),
			zap.Uint("user_id", user.ID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to create session", nil)
	}

	log.Info("Session created successfully",
		zap.Uint("session_id", session.ID),
		zap.String("session_uuid", sessionID),
		zap.Uint("user_id", user.ID),
	)

	// Generate JWT tokens
	accessToken, err := services.GenerateAccessToken(user.ID, user.Username, user.Email, sessionID)
	if err != nil {
		log.Error("Failed to generate access token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", user.ID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to generate access token", nil)
	}

	refreshToken, err := services.GenerateRefreshToken(user.ID, user.Username, user.Email, sessionID)
	if err != nil {
		log.Error("Failed to generate refresh token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", user.ID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to generate refresh token", nil)
	}

	log.Info("JWT tokens generated successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
	)

	// Set secure HTTP-only cookies
	// Access token cookie (short-lived)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(services.AccessTokenExpiry.Seconds()),
		HTTPOnly: true,
	})

	// Refresh token cookie (long-lived)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(services.RefreshTokenExpiry.Seconds()),
		HTTPOnly: true,
	})

	log.Info("Cookies set successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
	)

	// Prepare response
	response := &models.LoginResponse{
		User:      user.ToUserResponse(),
		SessionID: sessionID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		Message:   "Login successful",
	}

	log.Info("Login completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("session_id", sessionID),
	)

	return handlers.SuccessResponse(c, apiName, response, "Login successful")
}

// validateLoginRequest validates the login request
func validateLoginRequest(req *models.LoginRequest) error {
	// Must provide either username or email
	if req.Username == "" && req.Email == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Username or email is required")
	}

	// Cannot provide both
	if req.Username != "" && req.Email != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Provide either username or email, not both")
	}

	// Password is required
	if req.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Password is required")
	}

	return nil
}
