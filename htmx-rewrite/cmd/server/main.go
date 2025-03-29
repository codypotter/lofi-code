package main

import (
	"loficode/application"
	"log"
	"net/http"
)

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("public"))))

	app := application.New()
	http.HandleFunc("/api/posts/{slug}/comments", app.Comments)
	http.HandleFunc("/api/post-previews", app.PostPreviews)
	http.HandleFunc("/api/search-results", app.SearchResults)
	http.HandleFunc("/api/tags", app.Tags)
	http.HandleFunc("/api/subscribe", app.Subscribe)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
