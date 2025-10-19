package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/services/climbs"
)

// SetupClimbRoutes configures all climb-related routes
func SetupClimbRoutes(app *fiber.App) {
	// Climb routes group
	climbRoutes := app.Group("/climbs")

	// GET /climbs - Retrieve user climbs with optional date filtering
	climbRoutes.Get("/", climbs.GetClimbs)

	// POST /climbs - Create a new climb log entry
	climbRoutes.Post("/", climbs.CreateClimb)
}
