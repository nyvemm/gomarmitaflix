package controllers

import (
	"fmt"
	"marmitaflix/app/helpers"
	"marmitaflix/app/models"
	"strings"

	"github.com/anaskhan96/soup"
	"github.com/gofiber/fiber/v2"
)

func getSlugFromLink(link string) string {
	defaultUrl := helpers.GetEnv("DEFAULT_URL")
	slug := strings.Replace(link, defaultUrl, "", -1)
	return slug
}

func HandleNotFoundError(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Movie not found",
	})
}

func GetMovies(c *fiber.Ctx) error {
	defaultURL := helpers.GetEnv("DEFAULT_URL")

	resp, err := soup.Get(defaultURL)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)
	movies := doc.FindAll("div", "class", "capa_lista")

	var moviesList []models.ModelMovies

	for _, movie := range movies {
		movieTitle := movie.Find("a").Attrs()["title"]
		movieLink := movie.Find("a").Attrs()["href"]
		movieImage := movie.Find("img").Attrs()["src"]
		moviesList = append(moviesList, models.ModelMovies{
			Title: movieTitle,
			Image: movieImage,
			Slug:  getSlugFromLink(movieLink),
		})
	}

	return c.JSON(moviesList)
}

func GetMovie(c *fiber.Ctx) error {
	slug := c.Params("slug")
	defaultURL := helpers.GetEnv("DEFAULT_URL")
	url := fmt.Sprintf("%s%s/", defaultURL, slug)
	fmt.Println("URL: ", url)

	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)

	elementMovieImage := doc.Find("img", "class", "img-fluid")
	if elementMovieImage.Error != nil {
		fmt.Println("Error: ", elementMovieImage.Error)
		return HandleNotFoundError(c)
	}
	movieImage := elementMovieImage.Attrs()["src"]

	infoDiv := doc.Find("div", "id", "informacoes")
	if infoDiv.Error != nil {
		fmt.Println("Error: ", infoDiv.Error)
		return HandleNotFoundError(c)
	}
	infoDivFirstStrong := infoDiv.Find("strong")
	if infoDivFirstStrong.Error != nil {
		fmt.Println("Error: ", infoDivFirstStrong.Error)
		return HandleNotFoundError(c)
	}
	movieTitle := infoDivFirstStrong.Text()

	sinopseDiv := doc.Find("div", "id", "sinopse")
	if sinopseDiv.Error != nil {
		fmt.Println("Error: ", sinopseDiv.Error)
		return HandleNotFoundError(c)
	}

	sinopseDivFirstP := sinopseDiv.Find("p")
	if sinopseDivFirstP.Error != nil {
		fmt.Println("Error: ", sinopseDivFirstP.Error)
		return HandleNotFoundError(c)
	}
	movieSinopse := sinopseDivFirstP.Text()
	movieSinopse = movieSinopse[2:]
	movieSinopse = strings.ReplaceAll(movieSinopse, "\t", "")

	var movieEmbed string
	embedIframe := doc.Find("iframe", "class", "embed-responsive-item")
	if embedIframe.Error == nil {
		movieEmbed = embedIframe.Attrs()["src"]
	} else {
		movieEmbed = ""
	}

	var MovieLinks []models.ModelMovieLink

	downloadP := doc.Find("p", "id", "lista_download")
	if downloadP.Error != nil {
		fmt.Println("Error: ", downloadP.Error)
		return HandleNotFoundError(c)
	}

	downloadSpans := downloadP.FindAll("span")
	for _, downloadSpan := range downloadSpans {
		label := downloadSpan.Text()
		MovieLinks = append(MovieLinks, models.ModelMovieLink{
			Label: label,
			Link:  "",
		})
	}

	downloadLinks := downloadP.FindAll("a", "class", "btn")

	for index, downloadLink := range downloadLinks {
		link := downloadLink.Attrs()["href"]
		if index < len(MovieLinks) {
			MovieLinks[index].Link = link
		}
	}

	return c.JSON(models.ModelMovie{
		Title:       movieTitle,
		Image:       movieImage,
		Slug:        slug,
		Description: movieSinopse,
		Embed:       movieEmbed,
		Links:       MovieLinks,
	})
}

func SearchMovies(c *fiber.Ctx) error {
	search := c.Params("search")
	page := c.Params("page")

	if page == "" {
		page = "1"
	}

	defaultURL := helpers.GetEnv("DEFAULT_URL")
	url := fmt.Sprintf("%s%s/%s/", defaultURL, search, page)
	fmt.Println("URL: ", url)

	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)
	movies := doc.FindAll("div", "class", "capa_lista")

	var moviesList []models.ModelMovies

	for _, movie := range movies {
		movieTitle := movie.Find("a").Attrs()["title"]
		movieLink := movie.Find("a").Attrs()["href"]
		movieImage := movie.Find("img").Attrs()["src"]
		moviesList = append(moviesList, models.ModelMovies{
			Title: movieTitle,
			Image: movieImage,
			Slug:  getSlugFromLink(movieLink),
		})
	}

	return c.JSON(moviesList)
}
