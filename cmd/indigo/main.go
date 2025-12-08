//go:debug x509negativeserial=1
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

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

type runtimeFlags struct {
	port       int
	dumpRoutes bool
}

var flags = runtimeFlags{}

func init() {
	flag.IntVar(&flags.port, "port", 8080, "specifies the port the http server runs on")
	flag.BoolVar(&flags.dumpRoutes, "dump_routes", false, "dump the configured routes to stdout and exit")
	flag.Parse()
}

func main() {
	log.Print("Indigo CIL, v0.0.0")

	// Some flags may result in diagnostic data being dumped rather than
	// the server fully starting up. These flags should not require a connection to the
	// database or to LDAP.
	delayConnection := flags.dumpRoutes

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

	if !delayConnection {
		// Initialize connection
		err := l.Initialize()
		if err != nil {
			log.Fatal(err)
		}
		l.SetSecure(config.IndigoEnv == "prod")
		defer l.Connection.Close()
	}

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
	mux := c.CreateMux(
		c.Services{
			Config: config,
			Ldap:   &l,
			Jwt:    &jwtTransformer,
		},
	)

	if flags.dumpRoutes {
		fmt.Println(mux.DumpRoutes())
		os.Exit(0)
	}

	// Serve
	log.Printf("Serving on :%d\n", flags.port)
	log.Fatal(mux.ListenAndServe(fmt.Sprintf(":%d", flags.port)))
}
