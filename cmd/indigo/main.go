//go:debug x509negativeserial=1
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"

	c "myoptions.info/indigo/backend/internal/config"
	cntrl "myoptions.info/indigo/backend/internal/controller"
	"myoptions.info/indigo/backend/internal/db"
	"myoptions.info/indigo/backend/internal/repository"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

type runtimeFlags struct {
	port int
}

var flags = runtimeFlags{}

func init() {
	flag.IntVar(&flags.port, "port", 8080, "specifies the port the http server runs on")
	flag.Parse()
}

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

	// MARIADB DSN CONFIG
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		"127.0.0.1",
		"3306",
		"indigo_cil_dev",
	)

	// Configure GORM logger
	newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
	if config.IndigoEnv == "dev" {
		newLogger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	// Connect to MariaDB
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Fatalf("FATAL: Could not connect to MariaDB database: %v", err)
	}
	log.Printf("Successfully connected to MariaDB: %s", "indigo_cil_dev")

	// Run migrations
	if err := db.RunMigrations(gormDB); err != nil {
		log.Fatalf("FATAL: Database migration failed: %v", err)
	}

	// Initialize Repos, Services, Controllers
	repos := repository.NewRepositories(gormDB)
	services := s.NewServices(config.LdapUrl, config.IndigoEnv)
	// The blank identifier "_" is used to avoid the "declared and not used" error.
	_ = cntrl.NewControllers(repos, services)

	// Create routes using MuxWrapper
	mux := c.CreateRoutes(
		c.Services{
			Config: config,
			Ldap:   &l,
			Jwt:    &jwtTransformer,
		},
	)

	// Serve
	log.Printf("Serving on :%d\n", flags.port)
	log.Fatal(mux.ListenAndServe(fmt.Sprintf(":%d", flags.port)))
}
