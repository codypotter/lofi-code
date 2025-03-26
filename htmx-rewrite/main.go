package main

import (
	"loficode/templates/components"
	"loficode/templates/pages/home"
	"loficode/templates/pages/notfound"
	"loficode/templates/pages/post"
	"loficode/templates/pages/privacypolicy"
	"loficode/templates/pages/tos"
	"log"
	"net/http"
	"time"

	"github.com/a-h/templ"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets"))))
	http.Handle("/favicon.ico", http.FileServer(http.Dir("public/assets/images")))

	http.Handle("/admin/", http.StripPrefix("/admin/", http.FileServer(http.Dir("public/admin"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			notfound.NotFound().Render(r.Context(), w)
			return
		}
		home.Home().Render(r.Context(), w)
	})
	http.HandleFunc("/posts/{slug}", func(w http.ResponseWriter, r *http.Request) {
		post.Post(components.Post{
			Name:        "Hello, World!",
			Slug:        "hello-world",
			Description: "This is a description of the post.",
			Content:     "This is the content of the post.",
		}).Render(r.Context(), w)
	})
	http.Handle("/tos", templ.Handler(tos.TermsOfService()))
	http.Handle("/privacy-policy", templ.Handler(privacypolicy.PrivacyPolicy()))
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
				Slug:        "hello-world",
				Description: "This is a description of the post.",
				Content:     "This is the content of the post.",
			},
			{
				Name:        "Hello, World!",
				Slug:        "hello-world",
				Description: "This is a description of the post.",
				Content:     "This is the content of the post.",
			},
		}).Render(r.Context(), w)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))
}
