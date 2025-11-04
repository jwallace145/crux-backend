package main

import (
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/utils"

	"github.com/jwallace145/crux-backend/internal/config"
	"github.com/jwallace145/crux-backend/internal/middleware"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/routes"
)

func main() {
	// Initialize app, config, and logger
	app := fiber.New()
	cfg := config.Load()
	log := utils.Log

	// Initialize middleware
	corsMiddelware := middleware.CORSMiddleware()
	loggerMiddleware := middleware.LoggerMiddleware()
	authMiddleware := middleware.AuthMiddleware()

	// Attach global middleware for CORS and logging
	app.Use(corsMiddelware)
	app.Use(loggerMiddleware)

	// Connect to db and perform schema migrations
	db.ConnectDB()

	// Setup routes
	routes.SetupHealthCheckRoute(app)
	routes.SetupAuthRoutes(app, authMiddleware)
	routes.SetupUserRoutes(app, authMiddleware)
	routes.SetupClimbRoutes(app, authMiddleware)
	routes.SetupGymRoutes(app, authMiddleware)
	routes.SetupTrainingSessionRoutes(app, authMiddleware)
	routes.SetupDocsRoutes(app)

	log.Info("Starting CruxProject API server",
		zap.String("address", cfg.Address),
		zap.String("docs", "http://localhost:"+cfg.Port+"/docs"),
	)

	// Start server and listen for requests on the specified address
	if err := app.Listen(cfg.Address); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}
}
