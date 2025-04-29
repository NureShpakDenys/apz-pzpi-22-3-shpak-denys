// Package repository implements the repository interfaces
package repository // import "wayra/internal/adapter/repository"

import (
	"wayra/internal/adapter/config"
	"wayra/internal/core/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewGORMDB creates a new GORM database connection
// connectionString: connection string to the database
// returns: *gorm.DB, error
func NewGORMDB(cfg *config.Config) (*gorm.DB, error) {
	return gorm.Open(postgres.Open(cfg.StoragePath), &gorm.Config{})
}

// AutoMigrate runs the auto migration for the models
// db: database connection
// returns: error
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Company{},
		&models.Delivery{},
		&models.Role{},
		&models.Route{},
		&models.SensorData{},
		&models.User{},
		&models.Waypoint{},
		&models.UserCompany{},
		&models.Product{},
		&models.ProductCategory{},
		&models.Log{},
	)
}
