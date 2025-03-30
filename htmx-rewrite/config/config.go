package config

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	StorageUrl       string `json:"-" env:"STORAGE_URL"`
	BaseUrl          string `json:"-" env:"BASE_URL"`
	LogLevel         string `json:"-" env:"LOG_LEVEL"`
	DynamoDBRegion   string `json:"-" env:"DYNAMODB_REGION"`
	DynamoDBEndpoint string `json:"-" env:"DYNAMODB_ENDPOINT"`
}

func (c *Config) String() string {
	return "Config{[REDACTED]}"
}

func New() *Config {
	c := &Config{}
	err := env.Parse(c)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return c
}
