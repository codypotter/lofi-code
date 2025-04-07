package router

import (
	"loficode/internal/application"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(app application.Application) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/posts/{slug}/comments", func(r chi.Router) {
		r.Get("/", app.Comments)
		r.Post("/", app.CommentForm)
	})

	r.Get("/api/search-results", app.SearchResults)
	r.Post("/api/subscribe", app.Subscribe)
	r.Get("/api/verify", app.VerifyEmail)

	return r
}
