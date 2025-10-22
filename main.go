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

	// Attach middleware for CORS, logging, and authentication
	app.Use(middleware.CORSMiddleware())
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.AuthMiddleware())

	// Connect to db and perform schema migrations
	db.ConnectDB()

	// Setup API routes
	routes.SetupHealthCheckRoute(app)
	routes.SetupAuthRoutes(app)
	routes.SetupUserRoutes(app)
	routes.SetupClimbRoutes(app)
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
