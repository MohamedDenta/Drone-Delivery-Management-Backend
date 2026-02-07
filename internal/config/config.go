package config

import (
	"os"
)

type Config struct {
	DatabaseURL   string
	Port          string
	OTelCollector string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5433/drone_delivery?sslmode=disable"),
		Port:          getEnv("PORT", "8081"),
		OTelCollector: getEnv("OTEL_COLLECTOR_URL", "localhost:4317"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
