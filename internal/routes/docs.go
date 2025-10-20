package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/services/docs"
)

// SetupDocsRoutes configures all documentation-related routes
func SetupDocsRoutes(app *fiber.App) {
	// API documentation UI (Scalar)
	app.Get("/docs", docs.ServeAPIDocs)

	// OpenAPI specification file
	app.Get("/docs/openapi.yaml", docs.ServeOpenAPISpec)
}
