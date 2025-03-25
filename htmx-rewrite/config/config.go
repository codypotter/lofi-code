package config

import (
	"log"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	StorageUrl                string `json:"-" env:"STORAGE_URL"`
	BaseUrl                   string `json:"-" env:"BASE_URL"`
	LogLevel                  string `json:"-" env:"LOG_LEVEL"`
	FirebaseApiKey            string `json:"-" env:"FIREBASE_API_KEY"`
	FirebaseAuthDomain        string `json:"-" env:"FIREBASE_AUTH_DOMAIN"`
	FirebaseProjectId         string `json:"-" env:"FIREBASE_PROJECT_ID"`
	FirebaseStorageBucket     string `json:"-" env:"FIREBASE_STORAGE_BUCKET"`
	FirebaseMessagingSenderId string `json:"-" env:"FIREBASE_MESSAGING_SENDER_ID"`
	FirebaseAppId             string `json:"-" env:"FIREBASE_APP_ID"`
	FirebaseMeasurementId     string `json:"-" env:"FIREBASE_MEASUREMENT_ID"`
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
