package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/training_sessions"
)

func SetupTrainingSessionRoutes(app *fiber.App, authMiddleware fiber.Handler) {
	trainingSessionRoutes := app.Group("/training-sessions")

	// Protected routes (authentication required)
	trainingSessionRoutes.Get("/", authMiddleware, training_sessions.GetTrainingSessions)
	trainingSessionRoutes.Post("/", authMiddleware, training_sessions.CreateTrainingSession)
}
