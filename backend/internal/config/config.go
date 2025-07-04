package config

import (
	"os"
	"strings"
)

type Config struct {
	DatabaseURL       string
	APIPort           string
	Environment       string
	CORSOrigins       []string
	PrometheusEnabled bool
	PrometheusPort    string
}

func Load() *Config {
	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://qafka_user:qafka_password@localhost:5432/qafka?sslmode=disable"),
		APIPort:           getEnv("API_PORT", "8080"),
		Environment:       getEnv("ENVIRONMENT", "development"),
		CORSOrigins:       strings.Split(getEnv("CORS_ORIGINS", "http://localhost:5173"), ","),
		PrometheusEnabled: getEnv("PROMETHEUS_ENABLED", "true") == "true",
		PrometheusPort:    getEnv("PROMETHEUS_PORT", "2112"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
