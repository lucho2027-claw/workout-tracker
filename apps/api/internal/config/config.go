package config

import "os"

type Config struct {
	Port          string
	DatabaseURL   string
	JWTSecret     string
	AllowedOrigin string
	AppEnv        string
}

func Load() Config {
	return Config{
		Port:          getEnv("PORT", "8080"),
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://workout:workout@localhost:5432/workout_tracker?sslmode=disable"),
		JWTSecret:     getEnv("JWT_SECRET", "change-me"),
		AllowedOrigin: getEnv("ALLOWED_ORIGIN", "http://localhost:3000"),
		AppEnv:        getEnv("APP_ENV", "development"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
