package email

import (
	"context"
	"fmt"
	"loficode/internal/config"
	"loficode/internal/templates/components"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailSender interface {
	SendVerificationEmail(to, token string, subscribe bool) error
}

type NoopEmailSender struct{}

func NewNoopEmailSender() EmailSender {
	return NoopEmailSender{}
}

func (s NoopEmailSender) SendVerificationEmail(to, token string, subscribe bool) error {
	log.Printf("Sending verification email to %s with token %s", to, token)
	log.Printf("Verify email link: http://localhost:8080/api/verify?token=%s&subscribe=%t", token, subscribe)
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

func (s AwsSesEmailSender) SendVerificationEmail(to, token string, subscribe bool) error {
	url := fmt.Sprintf("https://loficode.com/api/verify?token=%s&subscribe=%t", token, subscribe)

	var sb strings.Builder
	err := components.VerifyEmail(to, token).Render(context.Background(), &sb)
	if err != nil {
		log.Printf("Error rendering email template: %v", err)
		return err
	}
	s.client.SendEmail(context.Background(), &ses.SendEmailInput{
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
