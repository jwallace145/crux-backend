package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/users"
	"github.com/jwallace145/crux-backend/internal/middleware"
)

func SetupUserRoutes(app *fiber.App) {
	userRoutes := app.Group("/users")

	// Public routes (no authentication required)
	userRoutes.Post("/", users.CreateUser)

	// Protected routes (authentication required)
	userRoutes.Get("/", middleware.AuthMiddleware(), users.GetUser)
	userRoutes.Put("/", middleware.AuthMiddleware(), users.UpdateUser)
}
