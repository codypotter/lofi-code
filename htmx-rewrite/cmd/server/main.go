package main

import (
	"context"
	"loficode/application"
	"loficode/config"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
	}))

	ctx := context.Background()
	cfg := config.New(ctx)
	app := application.New(ctx, cfg)

	r.Route("/api/posts/{slug}/comments", func(r chi.Router) {
		r.Get("/", app.Comments)
		r.Post("/", app.CommentForm)
	})

	r.Get("/api/search-results", app.SearchResults)
	r.Get("/api/tags", app.Tags)
	r.Post("/api/subscribe", app.Subscribe)
	r.Get("/api/verify", app.VerifyEmail)

	log.Fatal(http.ListenAndServe(":8080", r))
}
