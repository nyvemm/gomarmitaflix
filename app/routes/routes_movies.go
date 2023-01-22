package routes

import (
	"marmitaflix/app/controllers"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App) {
	app.Get("/movies/search/:search", func(c *fiber.Ctx) error {
		return controllers.SearchMovies(c)
	})
	app.Get("/movies/search/:search/:page", func(c *fiber.Ctx) error {
		return controllers.SearchMovies(c)
	})
	app.Get("/movies/:slug", func(c *fiber.Ctx) error {
		return controllers.GetMovie(c)
	})
	app.Get("/movies/all/:page", func(c *fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
	app.Get("/movies/all/:page", func(c *fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
	app.Get("/movies/categories/:category/:page", func(c *fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
}
