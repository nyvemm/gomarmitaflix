package controllers

import (
	"crypto/md5"
	"fmt"
	"marmitaflix/app/helpers"
	"marmitaflix/app/models"
	netURL "net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/anaskhan96/soup"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"

	"github.com/go-rod/rod"
	"github.com/go-rod/stealth"
)

func getSlugFromLink(link, defaultURL string) string {
	if defaultURL == "" {
		defaultURL = helpers.GetEnv("DEFAULT_URL")
	}
	slug := strings.Replace(link, defaultURL, "", -1)
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

	log.Info("URL: ", url)
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
			Slug:         getSlugFromLink(movieLink, ``),
			DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink, ``)),
		})
	}

	c.Set("Access-Control-Allow-Origin", "*")
	return c.JSON(moviesList)
}

func GetMagnetMovies(c fiber.Ctx) error {
	var htmlContent string
	var err error
	slug := c.Params("slug")
	defaultURL := helpers.GetEnv("DEFAULT_URL")
	url := fmt.Sprintf("%s%s/", defaultURL, slug)
	log.Info("URL: ", url)

	if !strings.Contains(url, "comando.la") {
		browser := rod.New().Timeout(time.Minute).MustConnect()
		defer browser.MustClose()

		log.Infof("js: %x\n\n", md5.Sum([]byte(stealth.JS)))

		mPage := stealth.MustPage(browser)

		mPage.MustNavigate(url)

		// TODO (il): improve how to get the element
		element := mPage.MustElement(`#main > div > div.post-block`)

		htmlContent, err = element.HTML()
		if err != nil {
			log.Error(err)
		}
	} else {
		htmlContent, err = soup.Get(url)
		if err != nil {
			panic(err)
		}
	}

	doc := soup.HTMLParse(htmlContent)
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
	log.Info("URL: ", url)

	resp, err := soup.Get(url)
	if err != nil {
		panic(err)
	}

	doc := soup.HTMLParse(resp)

	elementMovieImage := doc.Find("img", "class", "img-fluid")
	if elementMovieImage.Error != nil {
		log.Info("Error: ", elementMovieImage.Error)
		return HandleNotFoundError(c)
	}
	movieImage := elementMovieImage.Attrs()["src"]

	infoDiv := doc.Find("div", "id", "informacoes")
	if infoDiv.Error != nil {
		log.Info("Error: ", infoDiv.Error)
		return HandleNotFoundError(c)
	}
	infoDivFirstStrong := infoDiv.Find("strong")
	if infoDivFirstStrong.Error != nil {
		log.Info("Error: ", infoDivFirstStrong.Error)
		return HandleNotFoundError(c)
	}
	movieTitle := infoDivFirstStrong.Text()

	sinopseDiv := doc.Find("div", "id", "sinopse")
	if sinopseDiv.Error != nil {
		log.Info("Error: ", sinopseDiv.Error)
		return HandleNotFoundError(c)
	}

	sinopseDivFirstP := sinopseDiv.Find("p")
	if sinopseDivFirstP.Error != nil {
		log.Info("Error: ", sinopseDivFirstP.Error)
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
		log.Info("Error: ", downloadP.Error)
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

// SearchMoviesLemon tested in limontorrents
func SearchMoviesLemon(c fiber.Ctx) error {
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

	defaultURL := helpers.GetEnv("LEMON_URL")
	c.Set("Access-Control-Allow-Origin", "*")

	for {
		url = fmt.Sprintf("%spage/%d/?s=%s", defaultURL, pageInt, search)
		log.Info("URL: ", url)

		browser := rod.New().Timeout(time.Minute).MustConnect()

		log.Infof("js: %x\n\n", md5.Sum([]byte(stealth.JS)))

		mPage := stealth.MustPage(browser)

		mPage.MustNavigate(url)

		element := mPage.MustElement(`#main > div > div.movies-list`)

		htmlContent, err := element.HTML()
		if err != nil {
			log.Error(err)
		}

		doc := soup.HTMLParse(htmlContent)
		movies := doc.FindAll("div", "class", "item")

		if len(movies) == 0 {
			browser.MustClose()
			return c.JSON(moviesList)
		}

		for _, movie := range movies {
			movieTitle := movie.Find("div", "class", "title").Find("a").Text()
			movieLink := movie.Find("a").Attrs()["href"]
			movieImage := movie.Find("img").Attrs()["src"]
			moviesList = append(moviesList, models.ModelMovies{
				Title:        movieTitle,
				Image:        movieImage,
				Slug:         getSlugFromLink(movieLink, helpers.GetEnv("LEMON_URL")),
				DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink, helpers.GetEnv("LEMON_URL"))),
			})
		}
		if !loop {
			return c.JSON(moviesList)
		}
		browser.MustClose()

		if len(movies) < 21 {
			browser.MustClose()
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
			log.Info("URL: ", url)

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
					Slug:         getSlugFromLink(movieLink, ``),
					DownloadLink: fmt.Sprintf("%s/open/%s", c.BaseURL(), getSlugFromLink(movieLink, ``)),
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
