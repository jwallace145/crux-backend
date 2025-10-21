package db

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"

	"github.com/jwallace145/crux-backend/internal/utils"
)

// Config contains the connection parameters and configurations for the PostgreSQL database..
type Config struct {
	// Connection parameters
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string

	// Connection pool settings
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// LoadConfig loads db configurations from environment variables with sensible defaults.
// It supports loading from a .env file and applies intelligent defaults for SSL mode and
// connection pool settings.
func LoadConfig() (*Config, error) {
	// Load .env file if present (optional)
	if err := godotenv.Load(); err != nil {
		utils.Log.Info("No .env file found, using system environment variables")
	} else {
		utils.Log.Info("Successfully loaded .env file")
	}

	config := &Config{
		// Connection parameters
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Database: os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),

		// Connection pool settings with defaults
		MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
		MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
		ConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", time.Hour),
	}

	// Validate required fields
	if err := config.Validate(); err != nil {
		return nil, err
	}

	// Apply intelligent SSL defaults
	config.applySSLDefaults()

	utils.Log.Info("Database configuration loaded",
		zap.String("host", config.Host),
		zap.String("port", config.Port),
		zap.String("db", config.Database),
		zap.String("sslmode", config.SSLMode),
		zap.Int("maxIdleConns", config.MaxIdleConns),
		zap.Int("maxOpenConns", config.MaxOpenConns),
		zap.Duration("connMaxLifetime", config.ConnMaxLifetime),
	)

	return config, nil
}

// Validate checks that all required configuration fields are set.
func (c *Config) Validate() error {
	var missing []string

	if c.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if c.Port == "" {
		missing = append(missing, "DB_PORT")
	}
	if c.User == "" {
		missing = append(missing, "DB_USER")
	}
	if c.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if c.Database == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required db environment variables: %v", missing)
	}

	return nil
}

// applySSLDefaults applies intelligent SSL mode defaults based on the host.
// - If DB_SSLMODE is explicitly set, uses that value
// - If connecting to AWS RDS (*.rds.amazonaws.com), defaults to "require"
// - Otherwise, defaults to "disable" for local development
func (c *Config) applySSLDefaults() {
	if c.SSLMode != "" {
		utils.Log.Info("Using explicit SSL mode from environment",
			zap.String("sslmode", c.SSLMode),
		)
		return
	}

	if strings.Contains(c.Host, "rds.amazonaws.com") {
		c.SSLMode = "require"
		utils.Log.Info("AWS RDS detected, enabling SSL",
			zap.String("host", c.Host),
			zap.String("sslmode", c.SSLMode),
		)
	} else {
		c.SSLMode = "disable"
		utils.Log.Info("Local db detected, SSL disabled",
			zap.String("host", c.Host),
			zap.String("sslmode", c.SSLMode),
		)
	}
}

// DSN returns the PostgreSQL Data Source Name connection string.
func (c *Config) DSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Database, c.Port, c.SSLMode,
	)
}

// getEnvAsInt retrieves an environment variable as an integer.
// If the variable is not set or cannot be parsed, returns the default value.
func getEnvAsInt(key string, defaultVal int) int {
	if val := os.Getenv(key); val != "" {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
		utils.Log.Warn("Invalid integer value for environment variable, using default",
			zap.String("key", key),
			zap.String("value", val),
			zap.Int("default", defaultVal),
		)
	}
	return defaultVal
}

// getEnvAsDuration retrieves an environment variable as a time.Duration.
// Supports formats like "1h", "30m", "1h30m", etc.
// If the variable is not set or cannot be parsed, returns the default value.
func getEnvAsDuration(key string, defaultVal time.Duration) time.Duration {
	if val := os.Getenv(key); val != "" {
		if duration, err := time.ParseDuration(val); err == nil {
			return duration
		}
		utils.Log.Warn("Invalid duration value for environment variable, using default",
			zap.String("key", key),
			zap.String("value", val),
			zap.Duration("default", defaultVal),
		)
	}
	return defaultVal
}
