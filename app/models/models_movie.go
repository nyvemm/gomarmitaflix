package models

type ModelMovieLink struct {
	Label string `json:"label"`
	Link  string `json:"link"`
}

type ModelMovie struct {
	Title       string           `json:"title"`
	Slug        string           `json:"slug"`
	Image       string           `json:"image"`
	Description string           `json:"description"`
	Embed       string           `json:"embed"`
	Links       []ModelMovieLink `json:"links"`
}
