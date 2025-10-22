package climbs

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// CreateClimb handles POST /climbs requests to create a new climb log entry
// It validates the request, verifies optional route/gym references, and persists the climb
// Requires AuthMiddleware to be applied - reads user_id from context
func CreateClimb(c *fiber.Ctx) error {
	apiName := "create_climb"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting climb creation process",
		zap.String("api", apiName),
	)

	// Validate Content-Type header
	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	// Parse request body
	var req models.CreateClimbRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
		zap.String("climb_type", req.ClimbType),
		zap.String("grade", req.Grade),
	)

	// Validate request
	if err := validateCreateClimbRequest(&req); err != nil {
		log.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	log.Info("Request validation passed",
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

	// Verify route exists if RouteID is provided
	if req.RouteID != nil {
		log.Info("Verifying route exists",
			zap.String("api", apiName),
			zap.Uint("route_id", *req.RouteID),
		)

		var route models.Route
		if err := db.DB.First(&route, *req.RouteID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warn("Route not found",
					zap.String("api", apiName),
					zap.Uint("route_id", *req.RouteID),
				)
				return handlers.BadRequestResponse(c, apiName, "Route not found", map[string]interface{}{
					"route_id": *req.RouteID,
				})
			}
			log.Error("Database error while checking route",
				zap.Error(err),
				zap.String("api", apiName),
			)
			return handlers.InternalErrorResponse(c, apiName, "Failed to verify route", nil)
		}

		log.Info("Route verified",
			zap.String("api", apiName),
			zap.Uint("route_id", *req.RouteID),
		)
	}

	// Verify gym exists if GymID is provided
	if req.GymID != nil {
		log.Info("Verifying gym exists",
			zap.String("api", apiName),
			zap.Uint("gym_id", *req.GymID),
		)

		var gym models.Gym
		if err := db.DB.First(&gym, *req.GymID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				log.Warn("Gym not found",
					zap.String("api", apiName),
					zap.Uint("gym_id", *req.GymID),
				)
				return handlers.BadRequestResponse(c, apiName, "Gym not found", map[string]interface{}{
					"gym_id": *req.GymID,
				})
			}
			log.Error("Database error while checking gym",
				zap.Error(err),
				zap.String("api", apiName),
			)
			return handlers.InternalErrorResponse(c, apiName, "Failed to verify gym", nil)
		}

		log.Info("Gym verified",
			zap.String("api", apiName),
			zap.Uint("gym_id", *req.GymID),
		)
	}

	// Set defaults if not provided
	if req.Attempts == 0 {
		req.Attempts = 1
	}

	// Create climb
	log.Info("Creating climb in db",
		zap.String("api", apiName),
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
		log.Error("Failed to create climb in db",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", userID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to create climb", nil)
	}

	log.Info("Climb created successfully in db",
		zap.String("api", apiName),
		zap.Uint("climb_id", climb.ID),
		zap.Uint("user_id", userID),
		zap.Time("climb_date", climb.ClimbDate),
	)

	// Prepare response
	response := climb.ToClimbResponse()

	log.Info("Climb creation completed successfully",
		zap.String("api", apiName),
		zap.Uint("climb_id", climb.ID),
		zap.Uint("user_id", userID),
	)

	return handlers.CreatedResponse(c, apiName, response, "Climb created successfully")
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
