package application

import (
	"loficode/templates/components"
	"net/http"
	"time"
)

type application struct{}

func New() application {
	return application{}
}

func (a *application) Comments(w http.ResponseWriter, r *http.Request) {
	components.Comments([]components.Comment{
		{
			Text:      "This is a comment.",
			Timestamp: time.Now(),
			User:      "Alice",
		},
	}).Render(r.Context(), w)
}

func (a *application) PostPreviews(w http.ResponseWriter, r *http.Request) {
	components.PostPreviews([]components.Post{
		{
			Name:        "Hello, World!",
			Slug:        "hello-world",
			Description: "This is a description of the post.",
		},
		{
			Name:        "Hello, World!",
			Slug:        "hello-world-1",
			Description: "This is a description of the post.",
		},
		{
			Name:        "Hello, World!",
			Slug:        "hello-world-2",
			Description: "This is a description of the post.",
		},
	}).Render(r.Context(), w)
}

func (a *application) SearchResults(w http.ResponseWriter, r *http.Request) {
	components.SearchResults([]components.Post{
		{
			Name:        "Hello, World!",
			Slug:        "hello-world",
			Description: "This is a description of the post.",
		},
	}).Render(r.Context(), w)
}

func (a *application) Tags(w http.ResponseWriter, r *http.Request) {
	components.Tags(components.TagsConfig{
		Size:             "is-large",
		Tags:             []string{"tag1", "tag2", "tag3"},
		EnableNavigation: true,
	}).Render(r.Context(), w)
}
