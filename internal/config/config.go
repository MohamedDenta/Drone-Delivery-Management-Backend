package config

import (
	"os"
)

type Config struct {
	DatabaseURL   string
	RedisURL      string
	RabbitMQURL   string
	Port          string
	OTelCollector string
}

func Load() *Config {
	return &Config{
		DatabaseURL:   getEnv("DATABASE_URL", "postgres://user:password@localhost:5433/drone_delivery?sslmode=disable"),
		RedisURL:      getEnv("REDIS_URL", "localhost:6379"),
		RabbitMQURL:   getEnv("RABBITMQ_URL", "amqp://user:password@localhost:5672/"),
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
