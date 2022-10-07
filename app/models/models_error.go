package models

type ModelError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
