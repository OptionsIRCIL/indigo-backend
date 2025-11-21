//go:debug x509negativeserial=1
package main

import (
	"log"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	c "myoptions.info/indigo/backend/internal/config"
	cntrl "myoptions.info/indigo/backend/internal/controller"
	"myoptions.info/indigo/backend/internal/db"
	"myoptions.info/indigo/backend/internal/repository"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

func main() {
	log.Print("Indigo CIL, v0.0.0")

	// Load environment
	config := util.LoadConfig()
	log.Println("Initialized in environment", config.IndigoEnv)
	if config.IndigoEnv == "dev" {
		log.Println(
			"WARNING! Running in development mode removes various safeguards and encryption features.",
			"If you intend to deploy this software in a production environment, please use INDIGO_ENV=prod.",
		)
	}

	// Populate LdapConnection
	l := s.LdapConnection{
		Base:   config.LdapSearchBase,
		Domain: config.LdapDomain,
	}
	l.SetUrl(config.LdapUrl)
	l.SetCredentials(
		config.LdapUsername,
		config.LdapPassword,
	)

	// Initialize connection
	err := l.Initialize()
	if err != nil {
		log.Fatal(err)
	}
	l.SetSecure(config.IndigoEnv == "prod")
	defer l.Connection.Close()

	// Initialize JWT & secret
	jwtTransformer := s.JwtTransformer{}
	jwtInitErr := jwtTransformer.SetSecret([]byte(config.IndigoSecret))
	if jwtInitErr != nil {
		log.Fatal(jwtInitErr)
	}

	//
	dbFile := config.LdapUrl
	if dbFile == "" {
		dbFile = "indigo.db"
		log.Printf("INFO: config.LdapUrl not set, defaulting to SQLite file: %s", dbFile)
	}

	// Configure GORM logger to be quiet for production, more detailed for development
	newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
	if config.IndigoEnv == "dev" {
		newLogger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	gormDB, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("FATAL: Could not connect to database (%s): %v", dbFile, err)
	}
	log.Printf("Successfully connected to database: %s", dbFile)

	// run actual migrations
	if err := db.RunMigrations(gormDB); err != nil {
		log.Fatalf("FATAL: Database migration failed: %v", err)
	}

	// Initialize Repos
	repos := repository.NewRepositories(gormDB)

	// Initialize Services
	services := s.NewServices(config.LdapUrl, config.IndigoEnv)

	// Initialize Controllers
	// blank "_" for now until I spend some time on where to initialize correctly.
	_ = cntrl.NewControllers(repos, services)

	// Create routes
	mux := c.CreateRoutes(config, &l, &jwtTransformer)

	// Serve
	log.Printf("Serving on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
