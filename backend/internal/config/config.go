package config

import (
	"errors"
	"fmt"
	"os"
)

// MinJWTSecretLength is the minimum accepted length (in bytes) of JWT_SECRET.
const MinJWTSecretLength = 32

// legacyDefaultSecret was the hard-coded fallback used by earlier versions.
// It is publicly known, so it is rejected explicitly.
const legacyDefaultSecret = "mcpark-secret-change-in-production"

type Config struct {
	Port               string
	ListenAddr         string
	DatabasePath       string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
	AdminEmail         string
}

// Load reads configuration from the process environment.
func Load() (*Config, error) {
	return LoadWithGetenv(os.Getenv)
}

// LoadWithGetenv reads configuration through the supplied lookup function so
// tests can inject values without touching the real environment.
func LoadWithGetenv(getenv func(string) string) (*Config, error) {
	get := func(key, fallback string) string {
		if v := getenv(key); v != "" {
			return v
		}
		return fallback
	}

	secret := getenv("JWT_SECRET")
	if err := validateJWTSecret(secret); err != nil {
		return nil, err
	}

	return &Config{
		Port:               get("PORT", "8086"),
		ListenAddr:         get("LISTEN_ADDR", "127.0.0.1"),
		DatabasePath:       get("DATABASE_PATH", "./data/mcpark.db"),
		JWTSecret:          secret,
		GoogleClientID:     get("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: get("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  get("GOOGLE_REDIRECT_URL", "http://localhost:8086/api/auth/google/callback"),
		FrontendURL:        get("FRONTEND_URL", "http://localhost:5173"),
		AdminEmail:         get("ADMIN_EMAIL", ""),
	}, nil
}

// validateJWTSecret never includes the secret value in the returned error.
func validateJWTSecret(secret string) error {
	switch {
	case secret == "":
		return errors.New("JWT_SECRET is required")
	case secret == legacyDefaultSecret:
		return errors.New("JWT_SECRET must not be the legacy default value")
	case len(secret) < MinJWTSecretLength:
		return fmt.Errorf("JWT_SECRET must be at least %d characters", MinJWTSecretLength)
	}
	return nil
}
