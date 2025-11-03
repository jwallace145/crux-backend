package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/health"
)

func SetupHealthCheckRoute(app *fiber.App) {
	app.Get("/health", health.CheckHealth)
}
