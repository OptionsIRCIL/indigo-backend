package db

import (
	"log"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/models" // import models
)

// RunMigrations performs GORM AutoMigrate for all defined models.
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// all models mapped to database
	err := db.AutoMigrate(
		&models.AddressPhone{},
		&models.Alias{},
		&models.CommunityServiceEvent{},
		&models.ConsumerService{},
		&models.DisabilityInfo{},
		&models.Goal{},
		&models.InformationAndReferral{},
		&models.Organization{},
		&models.Person{},
		&models.Place{},
		&models.RecordDef{},
		&models.ServicesOffered{},
	)

	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully for all models.")
	return nil
}
