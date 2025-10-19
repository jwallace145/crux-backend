package climbs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// CreateClimb handles POST /climbs requests to create a new climb log entry
// It validates the request, verifies optional route/gym references, and persists the climb
func CreateClimb(c *fiber.Ctx) error {
	apiName := "create_climb"
	requestID := c.Get("X-Request-ID", "unknown")

	utils.Logger.Info("Starting climb creation process",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	// Validate Content-Type header
	if err := utils.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	// Parse request body
	var req models.CreateClimbRequest
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
		zap.String("climb_type", req.ClimbType),
		zap.String("grade", req.Grade),
	)

	// Validate request
	if err := validateCreateClimbRequest(&req); err != nil {
		utils.Logger.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	utils.Logger.Info("Request validation passed",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	// Get authenticated user ID from token
	// For now, we'll extract from access token cookie
	accessToken := c.Cookies("access_token", "")
	if accessToken == "" {
		utils.Logger.Warn("No access token found",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.UnauthorizedResponse(c, apiName, "Authentication required")
	}

	claims, err := utils.ValidateAccessToken(accessToken)
	if err != nil {
		utils.Logger.Warn("Invalid access token",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.UnauthorizedResponse(c, apiName, "Invalid or expired token")
	}

	userID := claims.UserID

	utils.Logger.Info("User authenticated",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", userID),
	)

	// Verify route exists if RouteID is provided
	if req.RouteID != nil {
		utils.Logger.Info("Verifying route exists",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("route_id", *req.RouteID),
		)

		var route models.Route
		if err := db.DB.First(&route, *req.RouteID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.Logger.Warn("Route not found",
					zap.String("api", apiName),
					zap.String("request_id", requestID),
					zap.Uint("route_id", *req.RouteID),
				)
				return utils.BadRequestResponse(c, apiName, "Route not found", map[string]interface{}{
					"route_id": *req.RouteID,
				})
			}
			utils.Logger.Error("Database error while checking route",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("request_id", requestID),
			)
			return utils.InternalErrorResponse(c, apiName, "Failed to verify route", nil)
		}

		utils.Logger.Info("Route verified",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("route_id", *req.RouteID),
		)
	}

	// Verify gym exists if GymID is provided
	if req.GymID != nil {
		utils.Logger.Info("Verifying gym exists",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("gym_id", *req.GymID),
		)

		var gym models.Gym
		if err := db.DB.First(&gym, *req.GymID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.Logger.Warn("Gym not found",
					zap.String("api", apiName),
					zap.String("request_id", requestID),
					zap.Uint("gym_id", *req.GymID),
				)
				return utils.BadRequestResponse(c, apiName, "Gym not found", map[string]interface{}{
					"gym_id": *req.GymID,
				})
			}
			utils.Logger.Error("Database error while checking gym",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("request_id", requestID),
			)
			return utils.InternalErrorResponse(c, apiName, "Failed to verify gym", nil)
		}

		utils.Logger.Info("Gym verified",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("gym_id", *req.GymID),
		)
	}

	// Set defaults if not provided
	if req.Attempts == 0 {
		req.Attempts = 1
	}

	// Create climb
	utils.Logger.Info("Creating climb in database",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("user_id", userID),
		zap.String("climb_type", req.ClimbType),
		zap.String("grade", req.Grade),
	)

	climb := &models.Climb{
		UserID:    userID,
		RouteID:   req.RouteID,
		GymID:     req.GymID,
		ClimbType: req.ClimbType,
		ClimbDate: req.ClimbDate,
		Grade:     req.Grade,
		Style:     req.Style,
		Completed: req.Completed,
		Attempts:  req.Attempts,
		Falls:     req.Falls,
		Rating:    req.Rating,
		Notes:     req.Notes,
	}

	if err := db.DB.Create(climb).Error; err != nil {
		utils.Logger.Error("Failed to create climb in database",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint("user_id", userID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to create climb", nil)
	}

	utils.Logger.Info("Climb created successfully in database",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("climb_id", climb.ID),
		zap.Uint("user_id", userID),
		zap.Time("climb_date", climb.ClimbDate),
	)

	// Prepare response
	response := climb.ToClimbResponse()

	utils.Logger.Info("Climb creation completed successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint("climb_id", climb.ID),
		zap.Uint("user_id", userID),
	)

	return utils.CreatedResponse(c, apiName, response, "Climb created successfully")
}

// validateCreateClimbRequest validates the create climb request
func validateCreateClimbRequest(req *models.CreateClimbRequest) error {
	// Validate climb type
	if req.ClimbType != models.ClimbTypeIndoor && req.ClimbType != models.ClimbTypeOutdoor {
		return fiber.NewError(fiber.StatusBadRequest, "Climb type must be 'indoor' or 'outdoor'")
	}

	// Validate grade is provided
	if req.Grade == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Grade is required")
	}

	if len(req.Grade) > 20 {
		return fiber.NewError(fiber.StatusBadRequest, "Grade must not exceed 20 characters")
	}

	// Validate climb date
	if req.ClimbDate.IsZero() {
		return fiber.NewError(fiber.StatusBadRequest, "Climb date is required")
	}

	// Don't allow future dates
	if req.ClimbDate.After(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "Climb date cannot be in the future")
	}

	// Validate style if provided
	if req.Style != "" && len(req.Style) > 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Style must not exceed 50 characters")
	}

	// Validate rating if provided
	if req.Rating < 0 || req.Rating > 5 {
		return fiber.NewError(fiber.StatusBadRequest, "Rating must be between 0 and 5")
	}

	// Validate notes length if provided
	if len(req.Notes) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "Notes must not exceed 1000 characters")
	}

	// Validate attempts
	if req.Attempts < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Attempts must be a positive number")
	}

	// Validate falls
	if req.Falls < 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Falls must be a positive number")
	}

	return nil
}
