package config

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/caarlos0/env/v8"
)

type Config struct {
	BaseUrl        string     `json:"-" env:"BASE_URL"`
	LogLevel       string     `json:"-" env:"LOG_LEVEL" envDefault:"info"`
	Environment    string     `json:"-" env:"ENVIRONMENT" envDefault:"development"`
	AwsConfig      aws.Config `json:"-"`
	HCaptchaSecret string     `json:"-" env:"HCAPTCHA_SECRET"`
}

func (c *Config) String() string {
	return "Config{[REDACTED]}"
}

func New(ctx context.Context) *Config {
	var c Config
	err := env.Parse(&c)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	c.AwsConfig, err = config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	return &c
}
