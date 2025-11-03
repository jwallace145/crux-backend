package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/auth"
)

func SetupAuthRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	// Public routes (no authentication required)
	app.Post("/login", auth.Login)
	app.Post("/refresh", auth.Refresh)

	// Protected routes (authentication required)
	app.Post("/logout", authMiddleware, auth.Logout)
}
