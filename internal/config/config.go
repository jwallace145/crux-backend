package config

import "os"

type AppConfig struct {
	ServiceName string
	Version     string
	Environment string
	Port        string
	Address     string
}

func Load() *AppConfig {
	port := getEnvOrDefault("PORT", "3000")

	return &AppConfig{
		ServiceName: "CruxProject API",
		Version:     "1.0.0",
		Environment: "development",
		Port:        port,
		Address:     "0.0.0.0:" + port,
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
