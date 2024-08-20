package controllers

import (
	"fmt"
	"marmitaflix/app/helpers"
	"marmitaflix/app/models"
	netURL "net/url"
	"strconv"
	"strings"
	"sync"

	"github.com/anaskhan96/soup"
	"github.com/gofiber/fiber/v3"
)

func getSlugFromLink(link string) string {
	defaultUrl := helpers.GetEnv("DEFAULT_URL")
	slug := strings.Replace(link, defaultUrl, "", -1)
	return slug
}

func HandleNotFoundError(c fiber.Ctx) error {
	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(fiber.Map{
		"status":  404,
		"message": "Movie not found",
	})
}

func GetMovies(c fiber.Ctx) error {
	defaultURL := helpers.GetEnv("DEFAULT_URL")
	page := c.Params("page")
	category := c.Params("category")

	if page == "" {
		page = "1"
	}

	url := fmt.Sprintf("%spage/%s", defaultURL, page)
	if category != "" {
		url = fmt.Sprintf("%s/category/%s/page/%s/", defaultURL, category, page)
	}

	fmt.Println("URL: ", url)
	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)
	movies := doc.FindAll("article", "class", "post")

	var moviesList []models.ModelMovies

	for _, movie := range movies {
		movieTitle := movie.Find("header").Find("h2").Find("a").Text()
		movieLink := movie.Find("a").Attrs()["href"]
		movieImage := movie.Find("img").Attrs()["src"]
		moviesList = append(moviesList, models.ModelMovies{
			Title:        movieTitle,
			Image:        movieImage,
			Slug:         getSlugFromLink(movieLink),
			DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink)),
		})
	}

	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(moviesList)
}

func GetMagnetMovies(c fiber.Ctx) error {
	slug := c.Params("slug")
	defaultURL := helpers.GetEnv("DEFAULT_URL")
	url := fmt.Sprintf("%s%s/", defaultURL, slug)
	fmt.Println("URL: ", url)

	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)
	magnetLinks := searchMagnetLinks(doc)

	var movies []models.ModelMovieLink

	for _, magnetLink := range magnetLinks {
		movies = append(movies, models.ModelMovieLink{
			Label: getTitleFromMagnetLink(magnetLink),
			Link:  fmt.Sprintf("%s/magnet/%s", c.BaseURL(), magnetLink),
		})
	}

	c.Set("Access-Control-Allow-Origin", "*")
	return c.Status(200).JSON(movies)

}

func OpenMagnet(c fiber.Ctx) error {
	magnet := strings.Replace(c.OriginalURL(), "/magnet/", "", -1)
	c.Set("Access-Control-Allow-Origin", "*")
	return c.Redirect().Status(302).To(magnet)
}

func GetMovie(c fiber.Ctx) error {
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

	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(models.ModelMovie{
		Title:       movieTitle,
		Image:       movieImage,
		Slug:        slug,
		Description: movieSinopse,
		Embed:       movieEmbed,
		Links:       MovieLinks,
	})
}

func SearchMoviesSync(c fiber.Ctx) error {
	var moviesList []models.ModelMovies
	var url string
	var loop bool

	search := c.Params("search")
	page := c.Params("page")

	pageInt, atoiErr := strconv.Atoi(page)
	if atoiErr != nil {
		loop = true
		pageInt = 1
	}

	defaultURL := helpers.GetEnv("DEFAULT_URL")
	c.Set("Access-Control-Allow-Origin", "*")

	for {
		url = fmt.Sprintf("%spage/%d/?s=%s", defaultURL, pageInt, search)
		fmt.Println("URL: ", url)

		resp, err := soup.Get(url)
		if err != nil {
			panic(err)
		}

		doc := soup.HTMLParse(resp)
		movies := doc.FindAll("article", "class", "post")

		if len(movies) == 0 {
			return c.JSON(moviesList)
		}

		for _, movie := range movies {
			movieTitle := movie.Find("header").Find("h2").Find("a").Text()
			movieLink := movie.Find("a").Attrs()["href"]
			movieImage := movie.Find("img").Attrs()["src"]
			moviesList = append(moviesList, models.ModelMovies{
				Title:        movieTitle,
				Image:        movieImage,
				Slug:         getSlugFromLink(movieLink),
				DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink)),
			})
		}
		if !loop {
			return c.JSON(moviesList)
		}
		pageInt++
	}
}

// TODO (il): implement async paging scrapper
func SearchMovies(c fiber.Ctx) error {
	var moviesList []models.ModelMovies

	search := c.Params("search")
	page := c.Params("page")

	pageInt, atoiErr := strconv.Atoi(page)
	if atoiErr != nil {
		pageInt = 1
	}

	defaultURL := helpers.GetEnv("DEFAULT_URL")
	c.Set("Access-Control-Allow-Origin", "*")

	var wg sync.WaitGroup
	var mu sync.Mutex

	errorChan := make(chan error, 1)
	stopChan := make(chan struct{})

	var once sync.Once

	pagesToFetch := 10
	for i := pageInt; i < pageInt+pagesToFetch; i++ {
		wg.Add(1)
		go func(p int) {
			defer wg.Done()

			select {
			case <-stopChan:
				return
			default:
			}

			url := fmt.Sprintf("%spage/%d/?s=%s", defaultURL, p, search)
			fmt.Println("URL: ", url)

			resp, err := soup.Get(url)
			if err != nil {
				errorChan <- err
				return
			}

			doc := soup.HTMLParse(resp)
			movies := doc.FindAll("article", "class", "post")

			if len(movies) == 0 {
				once.Do(func() {
					close(stopChan)
				})
				return
			}

			mu.Lock()
			defer mu.Unlock()

			for _, movie := range movies {
				movieTitle := movie.Find("header").Find("h2").Find("a").Text()
				movieLink := movie.Find("a").Attrs()["href"]
				movieImage := movie.Find("img").Attrs()["src"]
				moviesList = append(moviesList, models.ModelMovies{
					Title:        movieTitle,
					Image:        movieImage,
					Slug:         getSlugFromLink(movieLink),
					DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink)),
				})
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(errorChan)
	}()

	for err := range errorChan {
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}
	}

	return c.JSON(moviesList)
}

func sanitizeURLText(encodedText string) (string, error) {
	decodedText, err := netURL.QueryUnescape(encodedText)
	if err != nil {
		return "", err
	}
	return strings.Replace(decodedText, ".", " ", -1), nil
}

func searchMagnetLinks(doc soup.Root) []string {
	var magnets []string
	finder := doc.FindAll("a")
	for _, link := range finder {
		if strings.Contains(link.Attrs()["href"], "magnet") {
			magnets = append(magnets, link.Attrs()["href"])
		}
	}
	return magnets
}

func getTitleFromMagnetLink(link string) string {
	split := strings.Split(link, "&")
	for _, s := range split {
		if strings.Contains(s, "dn=") {
			title, err := sanitizeURLText(strings.Replace(s, "dn=", "", -1))
			if err != nil {
				return ""
			}
			return title
		}
	}
	return ""
}
