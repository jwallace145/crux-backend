package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: strings.Join([]string{
			// local development
			"http://localhost:3000",
			"http://localhost:3001",

			// dev api server for api docs
			"https://dev-api.cruxproject.io",

			// frontend
			"https://cruxproject.io",
			"https://www.cruxproject.io",
		}, ","),
		AllowMethods: strings.Join([]string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
			"HEAD",
		}, ","),
		AllowHeaders: strings.Join([]string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
			"X-Requested-With",
			"Access-Control-Request-Method",
			"Access-Control-Request-Headers",
		}, ","),
		AllowCredentials: true,
		ExposeHeaders: strings.Join([]string{
			"X-Request-ID",
			"Content-Length",
		}, ","),
		MaxAge: 43200, // cache preflight requests for 12 hours
	})
}
