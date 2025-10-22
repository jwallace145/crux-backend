package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/climbs"
	"github.com/jwallace145/crux-backend/internal/middleware"
)

func SetupClimbRoutes(app *fiber.App) {
	climbRoutes := app.Group("/climbs")

	// Public routes (no authentication required)
	climbRoutes.Get("/", climbs.GetClimbs)

	// Protected routes (authentication required)
	climbRoutes.Post("/", middleware.AuthMiddleware(), climbs.CreateClimb)
}
