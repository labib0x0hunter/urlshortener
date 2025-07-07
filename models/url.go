package models

import "time"

type reuqest struct {
	URL string `json:"url"`
}

type response struct {
	URL       string    `json:"url"`
	ShortURL  string    `json:"short_url"`
	CreatedAt time.Time `json:"created_at"`
	Expire    time.Time `json:"expire"`
}
