package application

import (
	"loficode/internal/templates/components"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (a application) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	e := r.FormValue("email")
	log := log.With().Str("handler", "Subscribe").Str("method", r.Method).Str("email", e).Logger()
	token, err := a.db.NewVerificationToken(ctx, e, 86400)
	if err != nil {
		log.Error().Err(err).Msg("Error creating verification token")
		components.Notification("is-danger", "Error creating verification token").Render(r.Context(), w)
		return
	}
	log.Debug().Msg("Created verification token")

	err = a.emailSender.SendVerificationEmail(ctx, e, token, true)
	if err != nil {
		log.Error().Err(err).Msg("Error sending verification email")
		components.Notification("is-danger", "Error sending verification email").Render(r.Context(), w)
		return
	}
	log.Info().Msg("Sent verification email")
	components.Notification("is-link", "Verification email sent!").Render(r.Context(), w)
}

func (a application) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Redirect(w, r, "/error.html", http.StatusSeeOther)
		return
	}
	subscribe := r.URL.Query().Get("subscribe")
	log := log.With().Str("handler", "VerifyEmail").Str("method", r.Method).Str("subscribe", subscribe).Str("token", token).Logger()
	log.Debug().Msg("Verifying email")

	_, err := a.db.VerifyEmail(ctx, token, subscribe == "true")
	if err != nil {
		log.Error().Err(err).Msg("Error verifying email")
		http.Redirect(w, r, "/error.html", http.StatusSeeOther)
		return
	}

	log.Info().Msg("Email verified")
	http.Redirect(w, r, "/verified.html", http.StatusSeeOther)
}
