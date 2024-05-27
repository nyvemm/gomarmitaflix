package main

import (
	"marmitaflix/app/helpers"
	"marmitaflix/app/middlewares"
	"marmitaflix/app/routes"
	"os"

	"github.com/gofiber/fiber/v3"

	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	helpers.LoadEnv()
	app := fiber.New()

	middlewares.CorsSetup(app)
	app.Use(recover.New())
	app.Use(logger.New())
	middlewares.CacheSetup(app)

	routes.Setup(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	app.Listen(":" + port)
}
