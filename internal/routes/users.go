package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/users"
)

func SetupUserRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	userRoutes := app.Group("/users")

	// Public routes (no authentication required)
	userRoutes.Post("/", users.CreateUser)

	// Protected routes (authentication required)
	userRoutes.Get("/", authMiddleware, users.GetUser)
	userRoutes.Put("/", authMiddleware, users.UpdateUser)
}
