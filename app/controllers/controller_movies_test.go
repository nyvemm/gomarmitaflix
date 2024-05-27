package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

func TestGetMovies(t *testing.T) {
	t.Setenv("DEFAULT_URL", "https://ondebaixo.com/")

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{{
		description:  "Should return 200",
		route:        "/movies",
		expectedCode: 200,
	},
	}

	app := fiber.New()

	app.Get("/movies", GetMovies)

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)

		res, _ := app.Test(req)

		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}

func TestGetMovie(t *testing.T) {
	t.Setenv("DEFAULT_URL", "https://ondebaixo.com/")

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Should return 200",
			route:        "/movies/morte-morte-morte-legendado-torrent-baixar-download/",
			expectedCode: 200,
		},
	}

	app := fiber.New()

	app.Get("/movies/:slug", GetMovie)

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)

		res, _ := app.Test(req)

		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}

func TestSearchMovies(t *testing.T) {
	t.Setenv("DEFAULT_URL", "https://ondebaixo.com/")

	tests := []struct {
		description  string
		route        string
		expectedCode int
	}{
		{
			description:  "Should return 200",
			route:        "/movies/search/morte",
			expectedCode: 200,
		},
	}

	app := fiber.New()

	app.Get("/movies/search/:search", SearchMovies)

	for _, test := range tests {
		req := httptest.NewRequest("GET", test.route, nil)

		res, _ := app.Test(req)

		assert.Equal(t, test.expectedCode, res.StatusCode, test.description)
	}
}
