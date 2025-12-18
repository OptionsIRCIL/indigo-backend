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
		// core/parents
		&models.Organization{},
		&models.Person{},
		&models.Place{},
		&models.RecordDef{},
		//reference above entities
		&models.Alias{},
		&models.Goal{},
		&models.ConsumerService{},
		&models.InformationAndReferral{},
		&models.ServicesOffered{},
		&models.DisabilityInfo{},
		//relies on others before migration possible
		&models.AddressPhone{},
		&models.CommunityServiceEvent{},
	)

	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully for all models.")
	return nil
}
