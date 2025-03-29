package application

import (
	"fmt"
	"loficode/model"
	"loficode/templates/components"
	"net/http"
	"net/url"
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
	components.PostPreviews([]model.Post{
		{
			Title:   "Hello, World!",
			Slug:    "hello-world",
			Summary: "This is a description of the post.",
		},
		{
			Title:   "Hello, World!",
			Slug:    "hello-world-1",
			Summary: "This is a description of the post.",
		},
		{
			Title:   "Hello, World!",
			Slug:    "hello-world-2",
			Summary: "This is a description of the post.",
		},
	}).Render(r.Context(), w)
}

func (a *application) SearchResults(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cursor := query.Get("cursor")
	currentURL := r.Header.Get("HX-Current-URL")
	if currentURL != "" {
		parsedURL, err := url.Parse(currentURL)
		if err == nil {
			query = parsedURL.Query()
		}
	}

	tags := query["tag"]
	fmt.Printf("SearchResults: %+v\n", tags)
	fmt.Printf("SearchResults: %+v\n", cursor)
	components.SearchResults([]model.Post{
		{
			Title:          "Hello, World!",
			Slug:           "being-right-is-overrated",
			Summary:        "This is a description of the post.",
			Tags:           []string{"foo", "bar"},
			Date:           time.Now(),
			HeaderImage:    "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Ftdtm4_clown%203x1.jpg?alt=media&token=c957fa6f-f715-4855-ae08-5c8fb0a564b4",
			OpenGraphImage: "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Fmkc6d_clown16x9.jpg?alt=media&token=e480a4c5-662d-41be-8d72-862bb1351e1f",
		},
		{
			Title:          "Hello, World!",
			Slug:           "being-right-is-overrated",
			Summary:        "This is a description of the post.",
			Tags:           []string{"foo", "bar"},
			Date:           time.Now(),
			HeaderImage:    "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Ftdtm4_clown%203x1.jpg?alt=media&token=c957fa6f-f715-4855-ae08-5c8fb0a564b4",
			OpenGraphImage: "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Fmkc6d_clown16x9.jpg?alt=media&token=e480a4c5-662d-41be-8d72-862bb1351e1f",
		},
		{
			Title:          "Hello, World!",
			Slug:           "being-right-is-overrated",
			Summary:        "This is a description of the post.",
			Tags:           []string{"foo", "bar"},
			Date:           time.Now(),
			HeaderImage:    "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Ftdtm4_clown%203x1.jpg?alt=media&token=c957fa6f-f715-4855-ae08-5c8fb0a564b4",
			OpenGraphImage: "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Fmkc6d_clown16x9.jpg?alt=media&token=e480a4c5-662d-41be-8d72-862bb1351e1f",
		},
		{
			Title:          "Hello, World!",
			Slug:           "being-right-is-overrated",
			Summary:        "This is a description of the post.",
			Tags:           []string{"foo", "bar"},
			Date:           time.Now(),
			HeaderImage:    "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Ftdtm4_clown%203x1.jpg?alt=media&token=c957fa6f-f715-4855-ae08-5c8fb0a564b4",
			OpenGraphImage: "https://firebasestorage.googleapis.com/v0/b/lofi-code.appspot.com/o/images%2Fmkc6d_clown16x9.jpg?alt=media&token=e480a4c5-662d-41be-8d72-862bb1351e1f",
		},
	}, "foo").Render(r.Context(), w)
}

func (a *application) Tags(w http.ResponseWriter, r *http.Request) {
	components.Tags(components.TagsConfig{
		Size:             "is-large",
		Tags:             []string{"tag1", "tag2", "tag3"},
		EnableNavigation: true,
	}).Render(r.Context(), w)
}
