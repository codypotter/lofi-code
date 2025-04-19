package application

import (
	"loficode/internal/model"
	"loficode/internal/templates/components"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func (a application) Comments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slug := chi.URLParam(r, "slug")
	log := log.With().Str("handler", "Comments").Str("method", r.Method).Str("slug", slug).Logger()

	comments, err := a.db.GetCommentsBySlug(ctx, slug)
	if err != nil {
		log.Error().Err(err).Msg("Error getting comments")
		components.Notification("is-warning", "Error getting comments").Render(r.Context(), w)
		return
	}

	log.Debug().Int("comments", len(comments)).Msg("Got comments")
	components.Comments(comments).Render(r.Context(), w)
}

func (a application) CommentForm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	fc := components.CommentFormConfig{
		Slug:    chi.URLParam(r, "slug"),
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Comment: r.FormValue("comment"),
	}
	log := log.With().Str("handler", "CommentForm").Str("method", r.Method).Any("formConfig", fc).Logger()

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		ip = r.RemoteAddr
	}
	success, err := a.hCaptcha.VerifyHCaptcha(ip, r.FormValue("h-captcha-response"))
	if err != nil {
		log.Error().Err(err).Msg("Error verifying hCaptcha")
		components.Notification("is-danger", "Error verifying hCaptcha").Render(ctx, w)
		components.CommentForm(fc).Render(ctx, w)
		return
	}
	if !success {
		log.Info().Msg("hCaptcha verification failed")
		components.Notification("is-danger", "Captcha verification failed. Please try again.").Render(ctx, w)
		components.CommentForm(fc).Render(ctx, w)
		return
	}

	mailingList := r.FormValue("mailingList") == "on"

	if fc.Name == "" || fc.Email == "" || fc.Comment == "" {
		log.Info().Msg("Required fields are empty")
		fc.Notification = components.Notification("is-danger", "All fields are required")
		components.CommentForm(fc).Render(ctx, w)
		return
	}

	verified, err := a.db.IsEmailVerified(ctx, fc.Email)
	if err != nil {
		log.Error().Err(err).Msg("Error checking email verification")
		fc.Notification = components.Notification("is-danger", "Error checking email verification")
		components.CommentForm(fc).Render(ctx, w)
		return
	}
	if !verified {
		log.Debug().Msg("Email is not verified, making a new verification token")
		token, err := a.db.NewVerificationToken(ctx, fc.Email, 86400)
		if err != nil {
			log.Printf("Error creating verification token: %v\n", err)
			fc.Notification = components.Notification("is-danger", "Error creating verification token")
			components.CommentForm(fc).Render(ctx, w)
			return
		}
		log.Debug().Msg("Sending verification email")
		err = a.emailSender.SendVerificationEmail(ctx, fc.Email, token, mailingList)
		if err != nil {
			log.Error().Err(err).Msg("Error sending verification email")
			fc.Notification = components.Notification("is-danger", "Error sending verification email")
			components.CommentForm(fc).Render(ctx, w)
			return
		}
		log.Debug().Msg("Verification email sent")
		fc.Notification = components.Notification("is-link", "Email not verified. Please check your inbox for a verification email. Once verified, try again.")
		log.Warn().Any("fc", fc).Msg("Comment form config")
		components.CommentForm(fc).Render(ctx, w)
		return
	}

	log.Debug().Msg("Email is verified, adding comment")
	if err := a.db.AddComment(ctx, fc.Slug, model.Comment{
		Name:  fc.Name,
		Email: fc.Email,
		Text:  fc.Comment,
		Date:  time.Now(),
	}); err != nil {
		log.Error().Err(err).Msg("Error adding comment")
		fc.Notification = components.Notification("is-danger", "Error adding comment")
		components.CommentForm(fc).Render(ctx, w)
		return
	}
	fc.Notification = components.Notification("is-success", "Comment added!")
	fc.Success = true
	fc.Name = ""
	fc.Email = ""
	fc.Comment = ""
	log.Info().Msg("Comment added successfully")
	components.CommentForm(fc).Render(ctx, w)
}
