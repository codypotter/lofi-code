package application

import (
	"context"
	"loficode/internal/config"
	"loficode/internal/db"
	"loficode/internal/email"
	"net/http"

	"github.com/rs/zerolog/log"
)

type Application interface {
	Comments(w http.ResponseWriter, r *http.Request)
	CommentForm(w http.ResponseWriter, r *http.Request)
	SearchResults(w http.ResponseWriter, r *http.Request)
	Subscribe(w http.ResponseWriter, r *http.Request)
	VerifyEmail(w http.ResponseWriter, r *http.Request)
}

type application struct {
	db          *db.Db
	cfg         *config.Config
	emailSender email.EmailSender
}

func New(ctx context.Context, cfg *config.Config) Application {
	log.Debug().Msg("Bootstrapping application")
	app := application{
		db:          db.New(ctx, cfg),
		cfg:         cfg,
		emailSender: email.NewAwsSesEmailSender(cfg),
	}
	return app
}

func NewDevelopment(ctx context.Context, cfg *config.Config) Application {
	log.Debug().Msg("Bootstrapping application in development mode")
	app := application{
		db:          db.NewDevelopment(ctx, cfg),
		cfg:         cfg,
		emailSender: email.NewNoopEmailSender(),
	}
	return app
}
