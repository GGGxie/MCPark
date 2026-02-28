package config

import "os"

type Config struct {
	Port               string
	DatabasePath       string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
	FrontendURL        string
	AdminEmail         string
}

func Load() *Config {
	return &Config{
		Port:               getEnv("PORT", "8086"),
		DatabasePath:       getEnv("DATABASE_PATH", "./data/mcpark.db"),
		JWTSecret:          getEnv("JWT_SECRET", "mcpark-secret-change-in-production"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8086/api/auth/google/callback"),
		FrontendURL:        getEnv("FRONTEND_URL", "http://localhost:5173"),
		AdminEmail:         getEnv("ADMIN_EMAIL", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
