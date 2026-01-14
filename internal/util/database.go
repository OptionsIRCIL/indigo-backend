package util

import (
	"log"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/model/entity"
)

// RunMigrations performs GORM AutoMigrate for all defined model.
func RunMigrations(db *gorm.DB) error {
	log.Println("Starting database migration...")

	// all model mapped to database
	err := db.AutoMigrate(
		// core/parents
		&entity.Organization{},
		&entity.Person{},
		&entity.Place{},
		&entity.RecordDef{},
		//reference above entities
		&entity.Alias{},
		&entity.Goal{},
		&entity.ConsumerService{},
		&entity.InformationAndReferral{},
		&entity.ServicesOffered{},
		&entity.DisabilityInfo{},
		//relies on others before migration possible
		&entity.AddressPhone{},
		&entity.CommunityServiceEvent{},
	)

	if err != nil {
		log.Printf("Database migration failed: %v", err)
		return err
	}

	log.Println("Database migration completed successfully for all model.")
	return nil
}
