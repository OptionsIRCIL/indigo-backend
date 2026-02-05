package routine

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	c "myoptions.info/indigo/backend/internal/config"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

// RunServe serves the application proper. Initializes all services and listens on either a port or a socket.
func RunServe(flags util.ServeRuntimeFlags) int {
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
	l.SetSecure((config.IndigoEnv == "prod") && !flags.AllowInsecureLdap)
	defer l.Connection.Close()

	// Initialize JWT & secret
	jwtTransformer := s.JwtTransformer{}
	jwtInitErr := jwtTransformer.SetSecret([]byte(config.IndigoSecret))
	if jwtInitErr != nil {
		log.Fatal(jwtInitErr)
	}

	// Configure GORM logger
	newLogger := gormlogger.Default.LogMode(gormlogger.Silent)
	if config.IndigoEnv == "dev" {
		newLogger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	// Connect to MariaDB
	database, err := gorm.Open(mysql.Open(config.DatabaseDsn), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
			NoLowerCase:   false,
		},
	})
	if err != nil {
		log.Fatalf("FATAL: Could not connect to MariaDB database: %v", err)
	}
	log.Printf("Successfully connected to MariaDB: %s", "indigo_cil_dev")

	// Run migrations
	if err := util.RunMigrations(database); err != nil {
		log.Fatalf("FATAL: Database migration failed: %v", err)
	}

	// Create routes using MuxWrapper
	mux := c.CreateMux(
		c.Services{
			Config:   config,
			Ldap:     &l,
			Jwt:      &jwtTransformer,
			Flags:    flags,
			Database: database,
		},
	)

	// Serve
	if flags.Socket == "" {
		log.Printf("Serving on :%d\n", flags.Port)
		log.Fatal(mux.ListenAndServe(fmt.Sprintf(":%d", flags.Port)))
	} else {
		log.Printf("Serving on socket %s\n", flags.Socket)
		log.Fatal(mux.ServeToSocket(flags.Socket, flags.SocketUid, flags.SocketGid))
	}

	return 0
}
