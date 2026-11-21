package main

import (
	"fmt"
	"go-auth-service/handlers"
	"go-auth-service/utils"
	"os"
	"time"
)

const (
	defaultPort     = "8081"
	defaultIssuer   = "expense-wise-auth"
	defaultAudience = "expense-wise-api"
	tokenTTL        = 24 * time.Hour
)

// Config is loaded and validated once at startup, so a missing variable
// stops the service immediately instead of failing on the first login.
type Config struct {
	Port string
	Auth handlers.AuthConfig
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// LoadConfig reads the environment and returns an error describing
// every missing or invalid setting.
func LoadConfig() (Config, error) {
	var problems []string

	username := os.Getenv("ADMIN_USERNAME")
	password := os.Getenv("ADMIN_PASSWORD")
	secret := os.Getenv("JWT_SECRET")

	if username == "" {
		problems = append(problems, "ADMIN_USERNAME is not set")
	}
	if password == "" {
		problems = append(problems, "ADMIN_PASSWORD is not set")
	}
	if secret == "" {
		problems = append(problems, "JWT_SECRET is not set")
	} else if len(secret) < utils.MinSecretBytes {
		problems = append(problems, fmt.Sprintf("JWT_SECRET must be at least %d bytes (got %d)", utils.MinSecretBytes, len(secret)))
	}

	if len(problems) > 0 {
		return Config{}, fmt.Errorf("invalid configuration: %v", problems)
	}

	return Config{
		Port: getEnv("PORT", defaultPort),
		Auth: handlers.AuthConfig{
			AdminUsername: username,
			AdminPassword: password,
			Token: utils.TokenConfig{
				Secret:   []byte(secret),
				Issuer:   getEnv("JWT_ISSUER", defaultIssuer),
				Audience: getEnv("JWT_AUDIENCE", defaultAudience),
				TTL:      tokenTTL,
			},
		},
	}, nil
}
