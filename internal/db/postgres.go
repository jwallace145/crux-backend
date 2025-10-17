package db

import (
	"fmt"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/jwallace145/crux-backend/internal/utils"

	"github.com/jwallace145/crux-backend/models"
)

// DB is the global database connection instance used throughout the application.
// It should be initialized once during application startup via ConnectDB().
var DB *gorm.DB

// ConnectDB establishes a connection to the PostgreSQL database and performs schema migrations.
// It loads database configuration from environment variables (with optional .env file support),
// creates a GORM database instance, and automatically migrates all required models.
//
// Environment Variables Required:
//   - DB_HOST: Database host address (e.g., "localhost")
//   - DB_USER: Database username
//   - DB_PASSWORD: Database password
//   - DB_NAME: Database name
//   - DB_PORT: Database port (typically "5432")
//
// The function will terminate the application if:
//   - Database connection cannot be established
//   - Schema migration fails
func ConnectDB() {
	utils.Logger.Info("Starting Database Connection")

	// Load environment variables from .env file if present
	if err := godotenv.Load(); err != nil {
		utils.Logger.Info("No .env file found, using system environment variables")
	} else {
		utils.Logger.Info("Successfully loaded .env file")
	}

	// Build PostgreSQL connection string
	dsn := buildDSN()
	utils.Logger.Info("Connecting to database",
		zap.String("host", os.Getenv("DB_HOST")),
		zap.String("port", os.Getenv("DB_PORT")),
		zap.String("database", os.Getenv("DB_NAME")),
	)

	// Configure GORM with custom logger for better visibility
	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	}

	// Attempt database connection
	utils.Logger.Info("Opening database connection")
	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		utils.Logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	utils.Logger.Info("Database connection established")

	// Configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		utils.Logger.Fatal("Failed to get underlying SQL DB", zap.Error(err))
	}

	// Set connection pool parameters for optimal performance
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	utils.Logger.Info("Connection pool configured",
		zap.Int("maxIdle", 10),
		zap.Int("maxOpen", 100),
		zap.Duration("maxLifetime", time.Hour),
	)

	// Verify connection is alive
	if err := sqlDB.Ping(); err != nil {
		utils.Logger.Fatal("Database ping failed", zap.Error(err))
	}
	utils.Logger.Info("Database ping successful")

	// Assign to global variable
	DB = db

	// Perform schema migrations
	utils.Logger.Info("Starting schema migration")
	if err := migrateModels(DB); err != nil {
		utils.Logger.Fatal("Schema migration failed", zap.Error(err))
	}

	utils.Logger.Info("Database initialization complete")
}

// buildDSN constructs the PostgreSQL Data Source Name from environment variables.
// It validates that all required environment variables are set and returns a
// formatted connection string.
//
// Returns:
//   - string: PostgreSQL DSN in the format required by the postgres driver
func buildDSN() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Validate required environment variables
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		utils.Logger.Warn("One or more database environment variables are not set",
			zap.Strings("required", []string{"DB_HOST", "DB_USER", "DB_PASSWORD", "DB_NAME", "DB_PORT"}),
		)
	}

	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
}

// migrateModels performs automatic schema migration for all application models.
// It creates tables, missing columns, indexes, and constraints based on model definitions.
// Existing data is preserved during migrations.
//
// Models migrated:
//   - User: User accounts and authentication
//   - Session: User session management
//   - Crag: Climbing area/location information
//   - Wall: Climbing wall within a crag
//   - Route: Individual climbing routes
//
// Parameters:
//   - db: GORM database instance
//
// Returns:
//   - error: nil on success, error describing the failure otherwise
func migrateModels(db *gorm.DB) error {
	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Session{},
		&models.Crag{},
		&models.Wall{},
		&models.Route{},
	}

	utils.Logger.Info("Starting model migration", zap.Int("modelCount", len(modelsToMigrate)))

	for i, model := range modelsToMigrate {
		modelName := fmt.Sprintf("%T", model)
		utils.Logger.Info("Migrating model",
			zap.Int("current", i+1),
			zap.Int("total", len(modelsToMigrate)),
			zap.String("model", modelName),
		)

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}
	}

	utils.Logger.Info("All models migrated successfully", zap.Int("modelCount", len(modelsToMigrate)))
	return nil
}

// GetDB returns the global database instance.
// This is useful for accessing the database connection in other packages
// without directly referencing the global variable.
//
// Returns:
//   - *gorm.DB: The global database connection, or nil if not initialized
func GetDB() *gorm.DB {
	if DB == nil {
		utils.Logger.Warn("Database accessed before initialization")
	}
	return DB
}
