package main

import (
	"context"
	"loficode/application"
	"loficode/config"
	"loficode/db"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx)
	db.New(ctx, cfg)
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
	}))

	app := application.New()

	r.Route("/api/posts/{slug}/comments", func(r chi.Router) {
		r.Get("/", app.Comments)
		r.Post("/", app.CommentForm)
	})

	r.Get("/api/post-previews", app.PostPreviews)
	r.Get("/api/search-results", app.SearchResults)
	r.Get("/api/tags", app.Tags)
	r.Post("/api/subscribe", app.Subscribe)

	log.Fatal(http.ListenAndServe(":8080", r))
}
