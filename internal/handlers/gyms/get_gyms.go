package gyms

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// GetGyms handles GET /gyms requests to retrieve gyms
// Query parameters:
//   - id (optional): The ID of a specific gym to retrieve
//   - state (optional): Filter gyms by state
//   - city (optional): Filter gyms by city
//
// If id is not provided, returns all gyms (optionally filtered by state/city)
// If id is provided, returns the specific gym with that ID
func GetGyms(c *fiber.Ctx) error {
	apiName := "get_gyms"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting get gyms process",
		zap.String("api", apiName),
	)

	// Get optional query parameters
	idStr := c.Query("id")
	state := c.Query("state")
	city := c.Query("city")

	// If ID is provided, return specific gym
	if idStr != "" {
		return getGymByID(c, apiName, log, idStr)
	}

	// Otherwise, return all gyms (optionally filtered by state/city)
	return getAllGyms(c, apiName, log, state, city)
}

// getGymByID retrieves a specific gym by ID
func getGymByID(c *fiber.Ctx, apiName string, log *zap.Logger, idStr string) error {
	log.Info("Retrieving specific gym by ID",
		zap.String("api", apiName),
		zap.String("gym_id", idStr),
	)

	gymID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		log.Warn("Invalid gym ID format",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("gym_id", idStr),
		)
		return handlers.BadRequestResponse(c, apiName, "id must be a valid number", nil)
	}

	log.Info("Gym ID parsed from query parameter",
		zap.String("api", apiName),
		zap.Uint64("gym_id", gymID),
	)

	// Fetch gym from db
	var gym models.Gym
	if err := db.DB.Where("id = ?", uint(gymID)).First(&gym).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Warn("Gym not found",
				zap.String("api", apiName),
				zap.Uint64("gym_id", gymID),
			)
			return handlers.NotFoundResponse(c, apiName, "Gym not found")
		}
		log.Error("Database error while looking up gym",
			zap.Error(err),
			zap.String("api", apiName),
			zap.Uint64("gym_id", gymID),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to fetch gym", nil)
	}

	log.Info("Gym fetched successfully",
		zap.String("api", apiName),
		zap.Uint("gym_id", gym.ID),
		zap.String("name", gym.Name),
	)

	// Return gym data
	response := gym.ToFullGymResponse()

	log.Info("Get gym by ID completed successfully",
		zap.String("api", apiName),
		zap.Uint("gym_id", gym.ID),
		zap.String("name", gym.Name),
	)

	return handlers.SuccessResponse(c, apiName, response, "Gym retrieved successfully")
}

// getAllGyms retrieves all gyms from the database with optional filtering
func getAllGyms(c *fiber.Ctx, apiName string, log *zap.Logger, state, city string) error {
	logFields := []zap.Field{zap.String("api", apiName)}

	if state != "" {
		logFields = append(logFields, zap.String("state", state))
	}
	if city != "" {
		logFields = append(logFields, zap.String("city", city))
	}

	log.Info("Retrieving gyms", logFields...)

	// Build query with optional filters
	var gyms []models.Gym
	query := db.DB.Order("name ASC")

	// Apply state filter if provided
	if state != "" {
		query = query.Where("state = ?", state)
	}

	// Apply city filter if provided
	if city != "" {
		query = query.Where("city = ?", city)
	}

	// Execute query
	if err := query.Find(&gyms).Error; err != nil {
		log.Error("Database error while querying gyms",
			zap.Error(err),
			zap.String("api", apiName),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to retrieve gyms", nil)
	}

	log.Info("Gyms retrieved successfully",
		zap.String("api", apiName),
		zap.Int("count", len(gyms)),
	)

	// Convert gyms to response DTOs
	gymResponses := make([]*models.FullGymResponse, len(gyms))
	for i, gym := range gyms {
		gymResponses[i] = gym.ToFullGymResponse()
	}

	// Prepare response
	responseData := map[string]interface{}{
		"gyms":  gymResponses,
		"count": len(gymResponses),
	}

	log.Info("Get all gyms completed successfully",
		zap.String("api", apiName),
		zap.Int("count", len(gyms)),
	)

	return handlers.SuccessResponse(c, apiName, responseData, "Gyms retrieved successfully")
}
