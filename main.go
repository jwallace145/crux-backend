package main

import (
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
		return c.SendString("API is running!!!")
	})

	if err := app.Listen(":3000"); err != nil {
		utils.Logger.Fatal("Failed to start server", zap.Error(err))
	}
}
