package routes

import (
	"marmitaflix/app/controllers"

	"github.com/gofiber/fiber/v3"
)

func Setup(app *fiber.App) {
	app.Get("/movies/search/:search", func(c fiber.Ctx) error {
		return controllers.SearchMovies(c)
	})
	app.Get("/lemon/search/:search", func(c fiber.Ctx) error {
		return controllers.SearchMoviesLemon(c)
	})
	app.Get("/movies/search/:search/:page", func(c fiber.Ctx) error {
		return controllers.SearchMovies(c)
	})
	app.Get("/movies/:slug", func(c fiber.Ctx) error {
		return controllers.GetMovie(c)
	})
	app.Get("/open/:slug", func(c fiber.Ctx) error {
		return controllers.GetMagnetMovies(c)
	})
	app.Get("/magnet/:magnet", func(c fiber.Ctx) error {
		return controllers.OpenMagnet(c)
	})

	app.Get("/movies/all/:page", func(c fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
	app.Get("/movies/all/:page", func(c fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
	app.Get("/movies/categories/:category/:page", func(c fiber.Ctx) error {
		return controllers.GetMovies(c)
	})
}
