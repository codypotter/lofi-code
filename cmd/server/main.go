package main

import (
	"context"
	"loficode/internal/application"
	"loficode/internal/assets"
	"loficode/internal/config"
	logger "loficode/internal/logger"
	"loficode/internal/router"
	"loficode/internal/templates/components/page"
	"net/http"

	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()
	cfg := config.New(ctx)
	logger.Configure(cfg.LogLevel, logger.LogFormatConsole)

	resolver, err := assets.NewResolver("public/assets/manifest.json", "/assets")
	if err != nil {
		log.Fatal().Err(err).Msg("Error loading asset manifest")
	}
	page.AssetURL = resolver.URL

	app := application.NewDevelopment(ctx, cfg)
	r := router.NewRouter(app)

	r.Get("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.StripPrefix("/", http.FileServer(http.Dir("public"))).ServeHTTP(w, r)
	}))

	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal().Err(err).Msg("Server failed to start")
	}
}
