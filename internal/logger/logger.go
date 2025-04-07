package logger

import (
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logFormat string

const (
	LogFormatJson    logFormat = "json"
	LogFormatConsole logFormat = "console"
)

func Configure(levelStr string, format logFormat) {
	level, err := zerolog.ParseLevel(levelStr)
	if err != nil {
		log.Warn().Err(err).Msg("Invalid log level, defaulting to Info")
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
	if format == LogFormatConsole {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func ZerologMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			log.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Str("remote_addr", r.RemoteAddr).
				Msg("Incoming request")

			next.ServeHTTP(w, r)

			log.Info().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Dur("duration", time.Since(start)).
				Msg("Request completed")
		})
	}
}
