package training_sessions

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// GetTrainingSessions handles GET /training-sessions requests to retrieve training sessions for a user
// Query parameters:
//   - start_date (optional): The start date for filtering training sessions (RFC3339 format, e.g., "2024-01-01T00:00:00Z")
//   - end_date (optional): The end date for filtering training sessions (RFC3339 format, e.g., "2024-12-31T23:59:59Z")
//
// If start_date is not provided, queries from the beginning of time
// If end_date is not provided, queries until now
// Requires AuthMiddleware to be applied - reads user_id from context
func GetTrainingSessions(c *fiber.Ctx) error {
	apiName := "get_training_sessions"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting get training sessions process",
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

	// Parse optional query parameters for date range
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Default to min time if start_date not provided
	var startDate time.Time
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			log.Warn("Invalid start_date format",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("start_date", startDateStr),
			)
			return handlers.BadRequestResponse(c, apiName, "start_date must be in RFC3339 format (e.g., 2024-01-01T00:00:00Z)", nil)
		}
		startDate = parsedStartDate
		log.Info("Start date parsed from query parameter",
			zap.String("api", apiName),
			zap.Time("start_date", startDate),
		)
	} else {
		// Use minimum time value (beginning of time)
		startDate = time.Time{}
		log.Info("No start_date provided, using minimum time",
			zap.String("api", apiName),
		)
	}

	// Default to current time if end_date not provided
	var endDate time.Time
	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			log.Warn("Invalid end_date format",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("end_date", endDateStr),
			)
			return handlers.BadRequestResponse(c, apiName, "end_date must be in RFC3339 format (e.g., 2024-12-31T23:59:59Z)", nil)
		}
		endDate = parsedEndDate
		log.Info("End date parsed from query parameter",
			zap.String("api", apiName),
			zap.Time("end_date", endDate),
		)
	} else {
		// Use current time as default
		endDate = time.Now()
		log.Info("No end_date provided, using current time",
			zap.String("api", apiName),
			zap.Time("end_date", endDate),
		)
	}

	// Validate that start_date is before end_date
	if !startDate.IsZero() && startDate.After(endDate) {
		log.Warn("start_date is after end_date",
			zap.String("api", apiName),
			zap.Time("start_date", startDate),
			zap.Time("end_date", endDate),
		)
		return handlers.BadRequestResponse(c, apiName, "start_date must be before or equal to end_date", nil)
	}

	log.Info("Querying training sessions for user within date range",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	// Query training sessions for the user within the date range
	var trainingSessions []models.TrainingSession
	query := db.DB.
		Where("user_id = ?", userID).
		Where("session_date >= ?", startDate).
		Where("session_date <= ?", endDate).
		Order("session_date DESC").
		Preload("Gym").
		Preload("Partners").
		Preload("IndoorBoulders").
		Preload("RopeClimbs").
		Find(&trainingSessions)

	if query.Error != nil {
		log.Error("Database error while querying training sessions",
			zap.Error(query.Error),
			zap.String("api", apiName),
			zap.Uint("user_id", userID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to retrieve training sessions", nil)
	}

	log.Info("Training sessions retrieved successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
		zap.Int("count", len(trainingSessions)),
	)

	// Convert training sessions to response DTOs
	sessionResponses := make([]*models.TrainingSessionResponse, len(trainingSessions))
	for i, session := range trainingSessions {
		sessionResponses[i] = session.ToTrainingSessionResponse()
	}

	// Prepare response
	responseData := map[string]interface{}{
		"training_sessions": sessionResponses,
		"count":             len(sessionResponses),
		"start_date":        startDate.Format(time.RFC3339),
		"end_date":          endDate.Format(time.RFC3339),
	}

	log.Info("Get training sessions completed successfully",
		zap.String("api", apiName),
		zap.Uint("user_id", userID),
		zap.Int("count", len(trainingSessions)),
	)

	return handlers.SuccessResponse(c, apiName, responseData, "Training sessions retrieved successfully")
}
