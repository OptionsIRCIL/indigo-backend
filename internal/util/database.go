package util

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/model/entity"
)

// ConnectToDatabase connects to a MariaDB database and performs any needed migrations.
func ConnectToDatabase() *gorm.DB {
	database, err := gorm.Open(mysql.Open(config.Config.Database.Dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
	})

	if err != nil {
		log.Fatalf("Could not connect to MariaDB database: %v", err)
	}
	log.Printf("Successfully connected to database")

	log.Println("Starting database migration...")

	// all model mapped to database
	err = database.AutoMigrate(
		&entity.Employee{},
		&entity.LocalUser{},
		&entity.Organization{},
		&entity.Person{},
		&entity.Place{},
		&entity.RecordDef{},

		&entity.Alias{},
		&entity.Goal{},
		&entity.ConsumerService{},
		&entity.InformationAndReferral{},
		&entity.ServicesOffered{},
		&entity.DisabilityInfo{},

		&entity.AddressPhone{},
		&entity.CommunityServiceEvent{},
	)

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migration completed successfully for all model.")
	return database
}
