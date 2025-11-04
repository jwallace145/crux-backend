package gyms

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/handlers"
	"github.com/jwallace145/crux-backend/internal/utils"
	"github.com/jwallace145/crux-backend/models"
)

// CreateGym handles POST /gyms requests to create a new climbing gym
// It validates the request, checks for required fields, and persists the new gym to the db
func CreateGym(c *fiber.Ctx) error {
	apiName := "create_gym"
	log := utils.GetLoggerFromContext(c)

	log.Info("Starting gym creation process",
		zap.String("api", apiName),
	)

	// Validate Content-Type header
	if err := handlers.ValidateJSONContentType(c, apiName); err != nil {
		// Error response is already sent by ValidateJSONContentType
		return err
	}

	// Parse request body
	var req models.CreateGymRequest
	if err := c.BodyParser(&req); err != nil {
		log.Error("Failed to parse request body",
			zap.Error(err),
			zap.String("api", apiName),
			zap.ByteString("raw_body", c.Body()),
		)
		return handlers.BadRequestResponse(c, apiName, "Invalid request body", err.Error())
	}

	log.Info("Request body parsed successfully",
		zap.String("api", apiName),
		zap.String("name", req.Name),
		zap.String("type", req.Type),
		zap.String("city", req.City),
		zap.String("country", req.Country),
	)

	// Validate required fields
	log.Info("Validating request parameters",
		zap.String("api", apiName),
	)

	if err := validateCreateGymRequest(&req); err != nil {
		log.Warn("Request validation failed",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("name", req.Name),
		)
		return handlers.ValidationErrorResponse(c, apiName, err.Error(), nil)
	}

	log.Info("Request validation passed",
		zap.String("api", apiName),
	)

	// Normalize fields
	originalName := req.Name
	req.Name = strings.TrimSpace(req.Name)
	req.City = strings.TrimSpace(req.City)
	req.Country = strings.TrimSpace(req.Country)
	if req.Email != "" {
		req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	}

	log.Info("Normalized request parameters",
		zap.String("api", apiName),
		zap.String("original_name", originalName),
		zap.String("normalized_name", req.Name),
	)

	// Create new gym
	log.Info("Creating new gym in db",
		zap.String("api", apiName),
		zap.String("name", req.Name),
		zap.String("type", req.Type),
		zap.String("city", req.City),
		zap.String("country", req.Country),
	)

	gym := &models.Gym{
		Name:            req.Name,
		Description:     req.Description,
		Type:            req.Type,
		Address:         req.Address,
		City:            req.City,
		State:           req.State,
		Province:        req.Province,
		Country:         req.Country,
		PostalCode:      req.PostalCode,
		Latitude:        req.Latitude,
		Longitude:       req.Longitude,
		Phone:           req.Phone,
		Email:           req.Email,
		Website:         req.Website,
		Hours:           req.Hours,
		HasBouldering:   req.HasBouldering,
		HasTopRope:      req.HasTopRope,
		HasLeadClimbing: req.HasLeadClimbing,
		HasAutoBelay:    req.HasAutoBelay,
		HasKidsArea:     req.HasKidsArea,
		HasTrainingArea: req.HasTrainingArea,
		HasYogaClasses:  req.HasYogaClasses,
		HasShower:       req.HasShower,
		HasParking:      req.HasParking,
		HasGearRental:   req.HasGearRental,
		HasProShop:      req.HasProShop,
		HasCafe:         req.HasCafe,
		WallHeight:      req.WallHeight,
		SquareFeet:      req.SquareFeet,
		DayPassPrice:    req.DayPassPrice,
		MonthlyPrice:    req.MonthlyPrice,
		YearlyPrice:     req.YearlyPrice,
		GearRentalPrice: req.GearRentalPrice,
		Notes:           req.Notes,
		Active:          req.Active,
	}

	if err := db.DB.Create(gym).Error; err != nil {
		log.Error("Failed to create gym in db",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("name", req.Name),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to create gym", nil)
	}

	log.Info("Gym created successfully in db",
		zap.String("api", apiName),
		zap.Uint("gym_id", gym.ID),
		zap.String("name", gym.Name),
		zap.String("city", gym.City),
		zap.String("country", gym.Country),
	)

	// Prepare response
	response := gym.ToFullGymResponse()

	log.Info("Gym creation completed successfully",
		zap.String("api", apiName),
		zap.Uint("gym_id", gym.ID),
		zap.String("name", gym.Name),
	)

	// Return created gym
	return handlers.CreatedResponse(c, apiName, response, "Gym created successfully")
}

// validateCreateGymRequest validates the create gym request
func validateCreateGymRequest(req *models.CreateGymRequest) error {
	if err := validateGymName(req.Name); err != nil {
		return err
	}
	if err := validateGymType(req.Type); err != nil {
		return err
	}
	if err := validateGymLocation(req.City, req.Country); err != nil {
		return err
	}
	if err := validateGymContactInfo(req.Email); err != nil {
		return err
	}
	if err := validateGymCoordinates(req.Latitude, req.Longitude); err != nil {
		return err
	}
	return nil
}

// validateGymName validates the gym name
func validateGymName(name string) error {
	if name == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Name is required")
	}
	if len(name) < 1 || len(name) > 200 {
		return fiber.NewError(fiber.StatusBadRequest, "Name must be between 1 and 200 characters")
	}
	return nil
}

// validateGymType validates the gym type
func validateGymType(gymType string) error {
	if gymType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Type is required")
	}
	if gymType != models.GymTypeBouldering && gymType != models.GymTypeRoped && gymType != models.GymTypeFull {
		return fiber.NewError(fiber.StatusBadRequest, "Type must be one of: bouldering, roped, full")
	}
	return nil
}

// validateGymLocation validates the city and country
func validateGymLocation(city, country string) error {
	if city == "" {
		return fiber.NewError(fiber.StatusBadRequest, "City is required")
	}
	if len(city) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "City must not exceed 100 characters")
	}
	if country == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Country is required")
	}
	if len(country) > 100 {
		return fiber.NewError(fiber.StatusBadRequest, "Country must not exceed 100 characters")
	}
	return nil
}

// validateGymContactInfo validates email format
func validateGymContactInfo(email string) error {
	if email != "" && !strings.Contains(email, "@") {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
	}
	return nil
}

// validateGymCoordinates validates latitude and longitude
func validateGymCoordinates(latitude, longitude *float64) error {
	if latitude != nil && (*latitude < -90 || *latitude > 90) {
		return fiber.NewError(fiber.StatusBadRequest, "Latitude must be between -90 and 90")
	}
	if longitude != nil && (*longitude < -180 || *longitude > 180) {
		return fiber.NewError(fiber.StatusBadRequest, "Longitude must be between -180 and 180")
	}
	return nil
}
