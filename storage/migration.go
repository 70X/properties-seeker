package storage

import (
	"log"

	"github.com/70X/properties-seeker/storage/src/models"
)

var MigrationModels = []interface{}{
	&models.Property{},
	&models.Execution{},
	&models.ScraperFilter{},
	&models.User{},
}

func RunMigrations(db *DB) error {
	log.Println("Starting database migrations...")

	for _, model := range MigrationModels {
		if err := db.AutoMigrate(model); err != nil {
			return err
		}
	}

	log.Println("Database migrations completed successfully")
	return nil
}
