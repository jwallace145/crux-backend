package climbs

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// GetClimbs handles GET /climbs requests to retrieve user climbs
// Query parameters:
//   - user_id (required): The ID of the user whose climbs to retrieve
//   - start_date (optional): Start date in RFC3339 format (e.g., "2024-01-01T00:00:00Z")
//   - end_date (optional): End date in RFC3339 format (e.g., "2024-12-31T23:59:59Z")
//
// If start_date is not provided, returns climbs from the beginning of time
// If end_date is not provided, returns climbs up to now
func GetClimbs(c *fiber.Ctx) error {
	apiName := "get_climbs"
	requestID := c.Get("X-Request-ID", "unknown")

	utils.Logger.Info("Starting get climbs process",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
	)

	// Get user_id query parameter (required)
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		utils.Logger.Warn("Missing user_id query parameter",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
		)
		return utils.BadRequestResponse(c, apiName, "user_id query parameter is required", nil)
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		utils.Logger.Warn("Invalid user_id format",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.String("user_id", userIDStr),
		)
		return utils.BadRequestResponse(c, apiName, "user_id must be a valid number", nil)
	}

	utils.Logger.Info("User ID parsed from query parameter",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint64("user_id", userID),
	)

	// Get optional start_date and end_date query parameters
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	// Set default dates
	// Use Unix epoch (January 1, 1970) as minimum date
	startDate := time.Unix(0, 0).UTC()
	// Use current time as maximum date
	endDate := time.Now().UTC()

	// Parse start_date if provided
	if startDateStr != "" {
		parsedStartDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			utils.Logger.Warn("Invalid start_date format",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("request_id", requestID),
				zap.String("start_date", startDateStr),
			)
			return utils.BadRequestResponse(c, apiName, "start_date must be in RFC3339 format (e.g., 2024-01-01T00:00:00Z)", nil)
		}
		startDate = parsedStartDate.UTC()
		utils.Logger.Info("Start date parsed",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Time("start_date", startDate),
		)
	} else {
		utils.Logger.Info("No start_date provided, using minimum date",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Time("start_date", startDate),
		)
	}

	// Parse end_date if provided
	if endDateStr != "" {
		parsedEndDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			utils.Logger.Warn("Invalid end_date format",
				zap.Error(err),
				zap.String("api", apiName),
				zap.String("request_id", requestID),
				zap.String("end_date", endDateStr),
			)
			return utils.BadRequestResponse(c, apiName, "end_date must be in RFC3339 format (e.g., 2024-12-31T23:59:59Z)", nil)
		}
		endDate = parsedEndDate.UTC()
		utils.Logger.Info("End date parsed",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Time("end_date", endDate),
		)
	} else {
		utils.Logger.Info("No end_date provided, using current time",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Time("end_date", endDate),
		)
	}

	// Validate that start_date is before end_date
	if startDate.After(endDate) {
		utils.Logger.Warn("Start date is after end date",
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Time("start_date", startDate),
			zap.Time("end_date", endDate),
		)
		return utils.BadRequestResponse(c, apiName, "start_date must be before or equal to end_date", nil)
	}

	// Query climbs from database
	utils.Logger.Info("Querying climbs from database",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint64("user_id", userID),
		zap.Time("start_date", startDate),
		zap.Time("end_date", endDate),
	)

	var climbs []models.Climb
	query := db.DB.Where("user_id = ? AND climb_date >= ? AND climb_date <= ?", uint(userID), startDate, endDate).
		Order("climb_date DESC"). // Most recent first
		Find(&climbs)

	if query.Error != nil {
		utils.Logger.Error("Database error while querying climbs",
			zap.Error(query.Error),
			zap.String("api", apiName),
			zap.String("request_id", requestID),
			zap.Uint64("user_id", userID),
		)
		return utils.InternalErrorResponse(c, apiName, "Failed to retrieve climbs", nil)
	}

	utils.Logger.Info("Climbs retrieved successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint64("user_id", userID),
		zap.Int("count", len(climbs)),
	)

	// Convert climbs to response DTOs
	climbResponses := make([]*models.ClimbResponse, len(climbs))
	for i, climb := range climbs {
		climbResponses[i] = climb.ToClimbResponse()
	}

	// Prepare response
	responseData := map[string]interface{}{
		"climbs":     climbResponses,
		"count":      len(climbResponses),
		"user_id":    uint(userID),
		"start_date": startDate.Format(time.RFC3339),
		"end_date":   endDate.Format(time.RFC3339),
	}

	utils.Logger.Info("Get climbs completed successfully",
		zap.String("api", apiName),
		zap.String("request_id", requestID),
		zap.Uint64("user_id", userID),
		zap.Int("count", len(climbs)),
	)

	return utils.SuccessResponse(c, apiName, responseData, "Climbs retrieved successfully")
}
