package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/docs"
)

func SetupDocsRoutes(app *fiber.App) {
	// Public routes (no authentication required)
	app.Get("/docs", docs.ServeAPIDocs)
	app.Get("/docs/openapi.yaml", docs.ServeOpenAPISpec)
}
