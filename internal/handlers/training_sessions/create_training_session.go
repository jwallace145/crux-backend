package training_sessions

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

// CreateTrainingSession handles POST /training-sessions requests to create a new training session
// It validates the request, verifies gym and partner references, and persists the session with climbs
// Requires AuthMiddleware to be applied - reads user_id from context
func CreateTrainingSession(c *fiber.Ctx) error {
	apiName := "create_training_session"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting training session creation process",
		zap.String("api", apiName),
	)

	// Validate Content-Type header
	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		return err
	}

	// Parse request body
	var req models.CreateTrainingSessionRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
		zap.Uint("gym_id", req.GymID),
		zap.Time("session_date", req.SessionDate),
	)

	// Validate request
	if err := validateCreateTrainingSessionRequest(&req); err != nil {
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

	// Verify gym exists
	if err := verifyGymExists(c, apiName, req.GymID); err != nil {
		return err
	}

	// Verify all partners exist if PartnerIDs are provided
	partners, err := verifyPartnersExist(c, apiName, req.PartnerIDs)
	if err != nil {
		return err
	}

	// Create training session
	log.Info("Creating training session in db",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
		zap.Uint("gym_id", req.GymID),
	)

	trainingSession := &models.TrainingSession{
		UserID:      userID,
		GymID:       req.GymID,
		SessionDate: req.SessionDate,
		Description: req.Description,
	}

	// Create training session with all related entities in a transaction
	err = createTrainingSessionWithRelations(trainingSession, partners, &req)
	if err != nil {
		log.Error("Failed to create training session in db",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint("user_id", userID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to create training session", nil)
	}

	log.Info("Training session created successfully in db",
		zap.String("api", apiName),
		zap.Uint("training_session_id", trainingSession.ID),
		zap.Uint("user_id", userID),
		zap.Time("session_date", trainingSession.SessionDate),
	)

	// Load the complete training session with all relationships for the response
	if err := db.DB.
		Preload("Gym").
		Preload("Partners").
		Preload("Boulders").
		Preload("RopeClimbs").
		First(trainingSession, trainingSession.ID).Error; err != nil {
		log.Error("Failed to load training session relationships",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to load training session", nil)
	}

	// Prepare response
	response := trainingSession.ToTrainingSessionResponse()

	log.Info("Training session creation completed successfully",
		zap.String("api", apiName),
		zap.Uint("training_session_id", trainingSession.ID),
		zap.Uint("user_id", userID),
		zap.Int("total_climbs", response.TotalClimbs),
		zap.Int("total_sends", response.TotalSends),
	)

	return handlers.CreatedResponse(c, apiName, response, "Training session created successfully")
}

// validateCreateTrainingSessionRequest validates the create training session request
func validateCreateTrainingSessionRequest(req *models.CreateTrainingSessionRequest) error {
	if err := validateSessionBasicInfo(req); err != nil {
		return err
	}
	if err := validateBoulders(req.Boulders); err != nil {
		return err
	}
	if err := validateRopeClimbs(req.RopeClimbs); err != nil {
		return err
	}
	return nil
}

// validateSessionBasicInfo validates the basic session information
func validateSessionBasicInfo(req *models.CreateTrainingSessionRequest) error {
	if req.GymID == 0 {
		return fiber.NewError(fiber.StatusBadRequest, "Gym ID is required")
	}
	if req.SessionDate.IsZero() {
		return fiber.NewError(fiber.StatusBadRequest, "Session date is required")
	}
	if req.SessionDate.After(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "Session date cannot be in the future")
	}
	if len(req.Description) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "Description must not exceed 1000 characters")
	}
	return nil
}

// validateBoulders validates all boulders in the request
func validateBoulders(boulders []models.BoulderRequest) error {
	for _, boulder := range boulders {
		if err := validateBoulder(boulder); err != nil {
			return err
		}
	}
	return nil
}

// validateBoulder validates a single boulder
func validateBoulder(boulder models.BoulderRequest) error {
	if boulder.Grade == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Boulder grade is required")
	}
	if len(boulder.Grade) > 20 {
		return fiber.NewError(fiber.StatusBadRequest, "Boulder grade must not exceed 20 characters")
	}
	if boulder.Outcome != models.BoulderOutcomeFell &&
		boulder.Outcome != models.BoulderOutcomeFlash &&
		boulder.Outcome != models.BoulderOutcomeOnSight &&
		boulder.Outcome != models.BoulderOutcomeRedpoint {
		return fiber.NewError(fiber.StatusBadRequest, "Boulder outcome must be 'Fell', 'Flash', 'Onsite', or 'Redpoint'")
	}
	if boulder.ColorTag != nil && len(*boulder.ColorTag) > 50 {
		return fiber.NewError(fiber.StatusBadRequest, "Boulder color tag must not exceed 50 characters")
	}
	if len(boulder.Notes) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "Boulder notes must not exceed 1000 characters")
	}
	return nil
}

// validateRopeClimbs validates all rope climbs in the request
func validateRopeClimbs(ropeClimbs []models.RopeClimbRequest) error {
	for _, ropeClimb := range ropeClimbs {
		if err := validateRopeClimb(ropeClimb); err != nil {
			return err
		}
	}
	return nil
}

// validateRopeClimb validates a single rope climb
func validateRopeClimb(ropeClimb models.RopeClimbRequest) error {
	if ropeClimb.ClimbType != models.RopeClimbTypeTopRope && ropeClimb.ClimbType != models.RopeClimbTypeLead {
		return fiber.NewError(fiber.StatusBadRequest, "Rope climb type must be 'TR' or 'Lead'")
	}
	if ropeClimb.Grade == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Rope climb grade is required")
	}
	if len(ropeClimb.Grade) > 20 {
		return fiber.NewError(fiber.StatusBadRequest, "Rope climb grade must not exceed 20 characters")
	}
	if ropeClimb.Outcome != models.RopeClimbOutcomeFell &&
		ropeClimb.Outcome != models.RopeClimbOutcomeHung &&
		ropeClimb.Outcome != models.RopeClimbOutcomeFlash &&
		ropeClimb.Outcome != models.RopeClimbOutcomeOnSight &&
		ropeClimb.Outcome != models.RopeClimbOutcomeRedpoint {
		return fiber.NewError(fiber.StatusBadRequest, "Rope climb outcome must be 'Fell', 'Hung', 'Flash', 'Onsite', or 'Redpoint'")
	}
	if len(ropeClimb.Notes) > 1000 {
		return fiber.NewError(fiber.StatusBadRequest, "Rope climb notes must not exceed 1000 characters")
	}
	return nil
}

// verifyGymExists verifies that the gym with the given ID exists
func verifyGymExists(c *fiber.Ctx, apiName string, gymID uint) error {
	log := utils.GetLoggerFromContext(c)

	log.Info("Verifying gym exists",
		zap.String("api", apiName),
		zap.Uint("gym_id", gymID),
	)

	var gym models.Gym
	if err := db.DB.First(&gym, gymID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Gym not found",
				zap.String("api", apiName),
				zap.Uint("gym_id", gymID),
			)
			return handlers.BadRequestResponse(c, apiName, "Gym not found", map[string]interface{}{
				"gym_id": gymID,
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
		zap.Uint("gym_id", gymID),
	)

	return nil
}

// verifyPartnersExist verifies that all partners with the given IDs exist
func verifyPartnersExist(c *fiber.Ctx, apiName string, partnerIDs []uint) ([]models.User, error) {
	log := utils.GetLoggerFromContext(c)

	var partners []models.User
	if len(partnerIDs) == 0 {
		return partners, nil
	}

	log.Info("Verifying training partners exist",
		zap.String("api", apiName),
		zap.Int("partner_count", len(partnerIDs)),
	)

	if err := db.DB.Where("id IN ?", partnerIDs).Find(&partners).Error; err != nil {
		log.Error("Database error while checking partners",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return nil, handlers.InternalErrorResponse(c, apiName, "Failed to verify partners", nil)
	}

	if len(partners) != len(partnerIDs) {
		log.Warn("One or more partners not found",
			zap.String("api", apiName),
			zap.Int("requested", len(partnerIDs)),
			zap.Int("found", len(partners)),
		)
		return nil, handlers.BadRequestResponse(c, apiName, "One or more partners not found", map[string]interface{}{
			"partner_ids": partnerIDs,
		})
	}

	log.Info("All partners verified",
		zap.String("api", apiName),
		zap.Int("partner_count", len(partners)),
	)

	return partners, nil
}

// createTrainingSessionWithRelations creates a training session with all related entities in a transaction
func createTrainingSessionWithRelations(trainingSession *models.TrainingSession, partners []models.User, req *models.CreateTrainingSessionRequest) error {
	return db.DB.Transaction(func(tx *gorm.DB) error {
		// Create the training session
		if err := tx.Create(trainingSession).Error; err != nil {
			return err
		}

		// Associate partners if any
		if len(partners) > 0 {
			if err := tx.Model(trainingSession).Association("Partners").Append(partners); err != nil {
				return err
			}
		}

		// Create boulders if any
		if len(req.Boulders) > 0 {
			boulders := make([]models.Boulder, len(req.Boulders))
			for i, boulderReq := range req.Boulders {
				boulder := boulderReq.ToBoulder()
				boulder.TrainingSessionID = trainingSession.ID
				boulders[i] = *boulder
			}
			if err := tx.Create(&boulders).Error; err != nil {
				return err
			}
			trainingSession.Boulders = boulders
		}

		// Create rope climbs if any
		if len(req.RopeClimbs) > 0 {
			ropeClimbs := make([]models.RopeClimb, len(req.RopeClimbs))
			for i, ropeClimbReq := range req.RopeClimbs {
				ropeClimb := ropeClimbReq.ToRopeClimb()
				ropeClimb.TrainingSessionID = trainingSession.ID
				ropeClimbs[i] = *ropeClimb
			}
			if err := tx.Create(&ropeClimbs).Error; err != nil {
				return err
			}
			trainingSession.RopeClimbs = ropeClimbs
		}

		return nil
	})
}
