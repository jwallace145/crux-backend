package db

import (
	"fmt"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	"github.com/jwallace145/crux-backend/internal/utils"

	"github.com/jwallace145/crux-backend/models"
)

var DB *gorm.DB

// GetDB returns the global PostgreSQL DB instance.
func GetDB() *gorm.DB {
	log := utils.Log

	if DB == nil {
		log.Warn("Database accessed before initialization!")
	}

	return DB
}

// ConnectDB establishes a connection to the PostgreSQL DB and performs schema migrations with GORM.
func ConnectDB() {
	log := utils.Log

	log.Info("Starting to establish PostgreSQL DB connection")

	// Load db configuration
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load DB configurations", zap.Error(err))
	}

	// Configure GORM
	gormConfig := &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	}

	// Connect to database
	log.Info("Opening DB connection")
	db, err := gorm.Open(postgres.Open(config.DSN()), gormConfig)
	if err != nil {
		log.Fatal("Failed to connect to DB", zap.Error(err))
	}
	log.Info("DB connection established")

	// Configure connection pool using config values
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get underlying SQL DB", zap.Error(err))
	}

	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)
	log.Info("DB Connection pool configured",
		zap.Int("maxIdle", config.MaxIdleConns),
		zap.Int("maxOpen", config.MaxOpenConns),
		zap.Duration("maxLifetime", config.ConnMaxLifetime),
	)

	// Verify connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("Database ping failed", zap.Error(err))
	}
	log.Info("Database ping successful")

	// Assign to global variable
	DB = db

	// Perform schema migrations
	log.Info("Starting schema migration")
	if err := migrateModels(DB); err != nil {
		log.Fatal("Schema migration failed", zap.Error(err))
	}

	log.Info("Database initialization complete")
}

// migrateModels performs automatic schema migration for all application models with GORM.
func migrateModels(db *gorm.DB) error {
	log := utils.Log

	modelsToMigrate := []interface{}{
		&models.User{},
		&models.Session{},
		&models.Crag{},
		&models.Wall{},
		&models.Route{},
		&models.Gym{},
		&models.Climb{},
		&models.TrainingSession{},
		&models.RopeClimb{},
		&models.Boulder{},
	}

	log.Info("Starting model migration", zap.Int("modelCount", len(modelsToMigrate)))

	for i, model := range modelsToMigrate {
		modelName := fmt.Sprintf("%T", model)
		log.Info("Migrating model",
			zap.Int("current", i+1),
			zap.Int("total", len(modelsToMigrate)),
			zap.String("model", modelName),
		)

		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", modelName, err)
		}
	}

	log.Info("All models migrated successfully", zap.Int("modelCount", len(modelsToMigrate)))
	return nil
}
