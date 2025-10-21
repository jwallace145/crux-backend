package docs

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/handlers"

	"github.com/jwallace145/crux-backend/internal/utils"
)

// ServeOpenAPISpec serves the OpenAPI specification file
func ServeOpenAPISpec(c *fiber.Ctx) error {
	apiName := "openapi_spec"
	log := utils.GetLoggerFromContext(c)

	log.Info("Serving OpenAPI specification",
		zap.String("api", apiName),
	)

	// Read the OpenAPI spec file
	specPath := "docs/openapi.yaml"
	specContent, err := os.ReadFile(specPath)
	if err != nil {
		log.Error("Failed to read OpenAPI spec file",
			zap.Error(err),
			zap.String("api", apiName),
			zap.String("path", specPath),
		)
		return handlers.InternalErrorResponse(c, apiName, "Failed to load API specification", nil)
	}

	// Set appropriate content type for YAML
	c.Set("Content-Type", "application/x-yaml")

	return c.Send(specContent)
}

// ServeAPIDocs serves the Scalar API documentation UI
func ServeAPIDocs(c *fiber.Ctx) error {
	apiName := "api_docs"
	log := utils.GetLoggerFromContext(c)

	log.Info("Serving API documentation UI",
		zap.String("api", apiName),
	)

	// HTML page with Scalar integration
	html := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Crux Backend API Documentation</title>
    <meta name="description" content="Interactive API documentation for the Crux Backend climbing activity tracker">
</head>
<body>
    <!-- Scalar API Reference -->
    <script
        id="api-reference"
        data-url="/docs/openapi.yaml"
        data-configuration='{"theme":"purple","layout":"modern","showSidebar":true}'
    ></script>

    <!-- Load Scalar from CDN -->
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
</body>
</html>`

	c.Set("Content-Type", "text/html; charset=utf-8")
	return c.SendString(html)
}
