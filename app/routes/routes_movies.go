package routes

import (
	"marmitaflix/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/movies", func(c *fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
	app.Get("/movies/:slug", func(c *fiber.Ctx) error {
		return controllers.GetMovie(c)
	})
}
