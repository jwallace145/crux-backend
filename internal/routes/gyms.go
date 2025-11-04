package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/gyms"
)

func SetupGymRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	gymRoutes := app.Group("/gyms")

	// Protected routes (authentication required)
	gymRoutes.Get("/", authMiddleware, gyms.GetGyms)
	gymRoutes.Post("/", authMiddleware, gyms.CreateGym)
}
