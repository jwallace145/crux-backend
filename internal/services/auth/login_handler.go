package auth

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// Login handles POST /login requests to authenticate users
// Accepts either username+password or email+password
// Returns JWT tokens via secure HTTP-only cookies
func Login(c *fiber.Ctx) error {
	apiName := "login"
	requestID := c.Get("X-Request-ID", "unknown")

	utils.Logger.Info("Starting login process",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	// Validate Content-Type header
	if err := utils.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	// Parse request body
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		utils.Logger.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	utils.Logger.Info("Request body parsed successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	// Validate request
	if err := validateLoginRequest(&req); err != nil {
		utils.Logger.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	// Normalize inputs
	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}
	if req.Username != "" {
		req.Username = strings.TrimSpace(req.Username)
	}

	// Find user by email or username
	utils.Logger.Info("Looking up user",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.String("username", req.Username),
		zap.String("email", req.Email),
	)

	var user models.User
	var err error

	if req.Email != "" {
		err = db.DB.Where("email = ?", req.Email).First(&user).Error
	} else {
		err = db.DB.Where("username = ?", req.Username).First(&user).Error
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.Logger.Warn("User not found",
				zap.String("api", apiName),
				zap.String("request_id", requestID),
				zap.String("username", req.Username),
				zap.String("email", req.Email),
			)
			return utils.UnauthorizedResponse(c, apiName, "Invalid credentials")
		}
		utils.Logger.Error("Database error while looking up user",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.InternalErrorResponse(c, apiName, "Authentication failed", nil)
	}

	utils.Logger.Info("User found, verifying password",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
	)

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		utils.Logger.Warn("Invalid password",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("user_id", user.ID),
			zap.String("username", user.Username),
		)
		return utils.UnauthorizedResponse(c, apiName, "Invalid credentials")
	}

	utils.Logger.Info("Password verified successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", user.ID),
	)

	// Generate session ID
	sessionID := uuid.New().String()

	utils.Logger.Info("Generated session ID",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.String("session_id", sessionID),
		zap.Uint("user_id", user.ID),
	)

	// Create session in database
	expiresAt := time.Now().Add(7 * 24 * time.Hour) // 1 week
	session := &models.Session{
		UserID:    user.ID,
		SessionID: sessionID,
		IP:        c.IP(),
		UserAgent: c.Get("User-Agent", "unknown"),
		Device:    c.Get("User-Agent", "unknown"), // Could be enhanced with device detection
		ExpiresAt: expiresAt,
		Revoked:   false,
	}

	if err := db.DB.Create(session).Error; err != nil {
		utils.Logger.Error("Failed to create session",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("user_id", user.ID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to create session", nil)
	}

	utils.Logger.Info("Session created successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("session_id", session.ID),
		zap.String("session_uuid", sessionID),
		zap.Uint("user_id", user.ID),
	)

	// Generate JWT tokens
	accessToken, err := utils.GenerateAccessToken(user.ID, user.Username, user.Email, sessionID)
	if err != nil {
		utils.Logger.Error("Failed to generate access token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("user_id", user.ID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to generate access token", nil)
	}

	refreshToken, err := utils.GenerateRefreshToken(user.ID, user.Username, user.Email, sessionID)
	if err != nil {
		utils.Logger.Error("Failed to generate refresh token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("user_id", user.ID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to generate refresh token", nil)
	}

	utils.Logger.Info("JWT tokens generated successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", user.ID),
	)

	// Set secure HTTP-only cookies
	// Access token cookie (short-lived)
	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(utils.AccessTokenExpiry.Seconds()),
		HTTPOnly: true,
		Secure:   true, // HTTPS only in production
		SameSite: "Strict",
	})

	// Refresh token cookie (long-lived)
	c.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(utils.RefreshTokenExpiry.Seconds()),
		HTTPOnly: true,
		Secure:   true, // HTTPS only in production
		SameSite: "Strict",
	})

	utils.Logger.Info("Cookies set successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", user.ID),
	)

	// Prepare response
	response := &models.LoginResponse{
		User:      user.ToUserResponse(),
		SessionID: sessionID,
		ExpiresAt: expiresAt.Format(time.RFC3339),
		Message:   "Login successful",
	}

	utils.Logger.Info("Login completed successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", user.ID),
		zap.String("username", user.Username),
		zap.String("session_id", sessionID),
	)

	return utils.SuccessResponse(c, apiName, response, "Login successful")
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
