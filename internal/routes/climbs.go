package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/jwallace145/crux-backend/internal/handlers/climbs"
)

func SetupClimbRoutes(app *fiber.App) {
	climbRoutes := app.Group("/climbs")
	climbRoutes.Get("/", climbs.GetClimbs)
	climbRoutes.Post("/", climbs.CreateClimb)
}
