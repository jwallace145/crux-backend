package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/climbs"
)

func SetupClimbRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	climbRoutes := app.Group("/climbs")

	// Protected routes (authentication required)
	climbRoutes.Get("/", authMiddleware, climbs.GetClimbs)
	climbRoutes.Post("/", authMiddleware, climbs.CreateClimb)
}
