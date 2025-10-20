package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/services/users"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(app *fiber.App) {
	// User routes group
	userRoutes := app.Group("/users")

	// GET /users - Get authenticated user
	userRoutes.Get("/", users.GetUser)

	// POST /users - Create a new user
	userRoutes.Post("/", users.CreateUser)
}
