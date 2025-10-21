package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/users"
)

func SetupUserRoutes(app *fiber.App) {
	userRoutes := app.Group("/users")
	userRoutes.Get("/", users.GetUser)
	userRoutes.Post("/", users.CreateUser)
}
