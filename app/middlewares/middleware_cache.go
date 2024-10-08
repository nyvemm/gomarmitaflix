package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cache"
)

func CacheSetup(app *fiber.App) {
	app.Use(cache.New(cache.Config{
		Next: func(c fiber.Ctx) bool {
			return c.Query("refresh") == "true"
		},
		Expiration: 10 * time.Minute,
	}))
}
