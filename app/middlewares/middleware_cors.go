package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CorsSetup(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
}
