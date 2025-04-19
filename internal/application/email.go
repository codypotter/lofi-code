package application

import (
	"loficode/internal/templates/components"
	"net/http"

	"github.com/rs/zerolog/log"
)

func (a application) Subscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	success, err := a.hCaptcha.VerifyHCaptcha(ip, r.FormValue("h-captcha-response"))
	if err != nil {
		log.Error().Err(err).Msg("Error verifying hCaptcha")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Error verifying hCaptcha",
		}).Render(r.Context(), w)
		return
	}
	if !success {
		log.Info().Msg("hCaptcha verification failed")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Captcha verification failed. Please try again.",
		}).Render(r.Context(), w)
		return
	}
	e := r.FormValue("email")
	verified, err := a.db.IsEmailVerified(ctx, e)
	if err != nil {
		log.Error().Err(err).Msg("Error checking email verification status")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Error checking email verification status",
		}).Render(r.Context(), w)
		return
	}
	if verified {
		log.Info().Msg("Email already verified")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Email already verified",
		}).Render(r.Context(), w)
		return
	}

	log := log.With().Str("handler", "Subscribe").Str("method", r.Method).Str("email", e).Logger()
	token, err := a.db.NewVerificationToken(ctx, e, 86400)
	if err != nil {
		log.Error().Err(err).Msg("Error creating verification token")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Error creating verification token",
		}).Render(r.Context(), w)
		return
	}
	log.Debug().Msg("Created verification token")

	err = a.emailSender.SendVerificationEmail(ctx, e, token, true)
	if err != nil {
		log.Error().Err(err).Msg("Error sending verification email")
		components.MailingListForm(components.MailingListConfig{
			ErrorMessage: "Error sending verification email",
		}).Render(r.Context(), w)
		return
	}
	log.Info().Msg("Sent verification email")
	components.Notification("is-link", "Verification email sent!").Render(r.Context(), w)
}

func (a application) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	success, err := a.hCaptcha.VerifyHCaptcha(ip, r.FormValue("h-captcha-response"))
	if err != nil {
		log.Error().Err(err).Msg("Error verifying hCaptcha")
		components.UnsubscribeForm(components.UnsubscribeFormConfig{
			ErrorMessage: "Error verifying hCaptcha",
		}).Render(r.Context(), w)
		return
	}
	if !success {
		log.Info().Msg("hCaptcha verification failed")
		components.UnsubscribeForm(components.UnsubscribeFormConfig{
			ErrorMessage: "Captcha verification failed. Please try again.",
		}).Render(r.Context(), w)
		return
	}
	e := r.FormValue("email")
	log := log.With().Str("handler", "Unsubscribe").Str("method", r.Method).Str("email", e).Logger()
	err = a.db.Unsubscribe(ctx, e)
	if err != nil {
		log.Error().Err(err).Msg("Error unsubscribing")
	} else {
		log.Debug().Msg("Unsubscribed successfully")
	}
	components.Notification("is-link", "If this email was subscribed, you've been unsubscribed successfully.").Render(r.Context(), w)
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
