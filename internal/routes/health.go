package routes

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers"
)

func SetupHealthCheckRoute(app *fiber.App) {
	app.Get("/health", func(c *fiber.Ctx) error {
		healthData := map[string]interface{}{
			"status": "healthy",
			"uptime": time.Now().UTC().String(),
		}
		return handlers.SuccessResponse(c, "health_check", healthData, "CruxProject API is running!")
	})
}
