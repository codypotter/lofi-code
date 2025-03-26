package main

import (
	"loficode/templates/components"
	"log"
	"net/http"
	"time"
)

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	http.HandleFunc("/api/posts/{slug}/comments", func(w http.ResponseWriter, r *http.Request) {
		components.Comments([]components.Comment{
			{
				Text:      "This is a comment.",
				Timestamp: time.Now(),
				User:      "Alice",
			},
		}).Render(r.Context(), w)
	})
	http.HandleFunc("/api/posts/{slug}/related", func(w http.ResponseWriter, r *http.Request) {
		components.RelatedPosts([]components.Post{
			{
				Name:        "Hello, World!",
				Slug:        "hello-world",
				Description: "This is a description of the post.",
				Content:     "This is the content of the post.",
			},
			{
				Name:        "Hello, World!",
				Slug:        "hello-world-1",
				Description: "This is a description of the post.",
				Content:     "This is the content of the post.",
			},
			{
				Name:        "Hello, World!",
				Slug:        "hello-world-2",
				Description: "This is a description of the post.",
				Content:     "This is the content of the post.",
			},
		}).Render(r.Context(), w)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
