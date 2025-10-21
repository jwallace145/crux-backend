package config

import "os"

type AppConfig struct {
	Port    string
	Address string
}

func Load() *AppConfig {
	port := getEnvOrDefault("PORT", "3000")

	return &AppConfig{
		Port:    port,
		Address: "0.0.0.0:" + port,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
