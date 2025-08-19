package config

import (
	"os"
)

type Config struct {
	HttpPort string
}

func MustLoad() Config {
	cfg := Config{
		HttpPort: getEnv("HTTP_PORT", "8080"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	if defaultValue != "" {
		return defaultValue
	}

	return ""
}
