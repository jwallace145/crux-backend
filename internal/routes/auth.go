package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/services/auth"
)

// SetupAuthRoutes configures all authentication-related routes
func SetupAuthRoutes(app *fiber.App) {
	// Auth routes - no group prefix, these are top-level endpoints
	app.Post("/login", auth.Login)
	app.Post("/logout", auth.Logout)
	app.Post("/refresh", auth.Refresh)
}
