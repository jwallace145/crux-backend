package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/auth"
)

func SetupAuthRoutes(app *fiber.App) {
	app.Post("/login", auth.Login)
	app.Post("/logout", auth.Logout)
	app.Post("/refresh", auth.Refresh)
}
