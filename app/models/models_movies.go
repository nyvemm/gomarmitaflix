package models

type ModelMovies struct {
	Title        string `json:"title"`
	Slug         string `json:"slug"`
	Image        string `json:"image"`
	DownloadLink string `json:"download_link"`
}
