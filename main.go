package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/db"
	"github.com/jwallace145/crux-backend/internal/utils"
)

func main() {
	defer utils.Sync()

	app := fiber.New()

	db.ConnectDB()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("CruxBackend API is running!")
	})

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
