package model

import "time"

type Post struct {
	Title          string    `json:"title"`
	Slug           string    `json:"slug"`
	Summary        string    `json:"summary"`
	Tags           []string  `json:"tags"`
	Date           time.Time `json:"date"`
	HeaderImage    string    `json:"headerImage"`
	OpenGraphImage string    `json:"openGraphImage"`
	Content        string    `json:"content"`
}
