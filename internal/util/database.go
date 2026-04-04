package util

import (
	"log"
	"net/url"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/model/entity"
)

// ConnectToDatabase connects to a MariaDB database and performs any needed migrations.
func ConnectToDatabase() *gorm.DB {
	// Parse charset and parseTime into URL params
	dsn, err := url.Parse(config.Config.Database.Dsn)
	if err != nil {
		log.Fatalf("Failed to read DSN: %s", err)
	}
	params := dsn.Query()

	if !params.Has("charset") {
		params.Add("charset", "utf8mb4")
	}

	if !params.Has("parseTime") {
		params.Add("parseTime", "True")
	}

	dsn.RawQuery = params.Encode()

	// Open connection
	database, err := gorm.Open(mysql.Open(dsn.String()), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
	})

	if err != nil {
		log.Fatalf("Could not connect to MariaDB database: %v", err)
	}
	log.Println("Successfully connected to database. Starting migration...")

	// all model mapped to database
	err = database.AutoMigrate(
		// Core entities
		&entity.Employee{},
		&entity.LocalUser{},
		&entity.Grant{},

		// People
		&entity.Person{},
		&entity.Alias{},
		&entity.Goal{},
		&entity.DisabilityInfo{},
		&entity.AddressPhone{},

		// Organizations
		&entity.Organization{},
		&entity.ServicesOffered{},
		&entity.ConsumerService{},

		// Information and Referral AKA Community Navigation AKA I&R AKA...
		&entity.InformationAndReferral{},
		&entity.InformationAndReferralEffort{},
		&entity.InformationAndReferralAttachment{},

		&entity.CommunityServiceEvent{},
	)

	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	log.Println("Database migration completed successfully for all model.")
	return database
}
