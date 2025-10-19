package main

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/routes"
	"github.com/jwallace145/crux-backend/internal/utils"
)

func main() {
	defer utils.Sync()

	app := fiber.New()

	// Setup CORS middleware FIRST to handle preflight requests
	// This must come before other middleware to ensure CORS headers are set properly
	app.Use(utils.CORSMiddleware())

	// Setup request logging middleware for audit trail
	app.Use(utils.RequestLogger())

	db.ConnectDB()

	// Global OPTIONS handler for CORS preflight requests
	// This ensures all OPTIONS requests receive proper CORS headers
	app.Options("/*", func(c *fiber.Ctx) error {
		// CORS middleware has already set the necessary headers
		// Just return 204 No Content for preflight requests
		return c.SendStatus(fiber.StatusNoContent)
	})

	// Health check endpoint
	app.Get("/health", func(c *fiber.Ctx) error {
		healthData := map[string]interface{}{
			"status": "healthy",
			"uptime": time.Now().UTC().String(),
		}
		return utils.SuccessResponse(c, "health_check", healthData, "CruxBackend API is running")
	})

	// Setup API routes
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupClimbRoutes(app)

	// Get port from environment variable or default to 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// CRITICAL: Listen on 0.0.0.0 instead of localhost for ECS compatibility
	address := "0.0.0.0:" + port
	utils.Logger.Info("Starting Crux API server", zap.String("address", address))

	if err := app.Listen(address); err != nil {
		utils.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
