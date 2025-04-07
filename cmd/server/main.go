package main

import (
	"context"
	"loficode/internal/application"
	"loficode/internal/config"
	logger "loficode/internal/logger"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx)
	logger.Configure(cfg.LogLevel, logger.LogFormatConsole)

	r := chi.NewRouter()

	r.Use(logger.ZerologMiddleware())

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
	}))

	app := application.NewDevelopment(ctx, cfg)

	r.Route("/api/posts/{slug}/comments", func(r chi.Router) {
		r.Get("/", app.Comments)
		r.Post("/", app.CommentForm)
	})

	r.Get("/api/search-results", app.SearchResults)
	r.Post("/api/subscribe", app.Subscribe)
	r.Get("/api/verify", app.VerifyEmail)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
