package email

import (
	"context"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/templates/components"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/rs/zerolog/log"
)

type EmailSender interface {
	SendVerificationEmail(ctx context.Context, to, token string, subscribe bool) error
}

type NoopEmailSender struct{}

func NewNoopEmailSender() EmailSender {
	return NoopEmailSender{}
}

func (s NoopEmailSender) SendVerificationEmail(_ context.Context, to, token string, subscribe bool) error {
	log.Debug().Str("to", to).
		Str("token", token).
		Bool("subscribe", subscribe).
		Msg("NoopEmailSender: Sending verification email")
	log.Info().Msgf("Verify email link: http://localhost:8080/api/verify?token=%s&subscribe=%t", token, subscribe)
	return nil
}

type AwsSesEmailSender struct {
	client *ses.Client
}

func NewAwsSesEmailSender(cfg *config.Config) EmailSender {
	return AwsSesEmailSender{
		client: ses.NewFromConfig(cfg.AwsConfig),
	}
}

func (s AwsSesEmailSender) SendVerificationEmail(ctx context.Context, to, token string, subscribe bool) error {
	url := fmt.Sprintf("https://loficode.com/api/verify?token=%s&subscribe=%t", token, subscribe)

	var sb strings.Builder
	err := components.VerifyEmail(to, token).Render(ctx, &sb)
	if err != nil {
		log.Error().Err(err).Str("to", to).Msg("Error rendering email template")
		return err
	}
	s.client.SendEmail(ctx, &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data:    aws.String("Verify your email"),
				Charset: aws.String("UTF-8"),
			},
			Body: &types.Body{
				Html: &types.Content{
					Data:    aws.String(sb.String()),
					Charset: aws.String("UTF-8"),
				},
				Text: &types.Content{
					Data:    aws.String("Please verify your email by clicking the link below:\n" + url),
					Charset: aws.String("UTF-8"),
				},
			},
		},
		Source: aws.String("verify@loficode.com"),
	})
	return nil
}
