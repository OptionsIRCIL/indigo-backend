//go:debug x509negativeserial=1
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	c "myoptions.info/indigo/backend/internal/config"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

type runtimeFlags struct {
	port       int
	socket     string
	dumpRoutes bool
}

var flags = runtimeFlags{}

func init() {
	flag.IntVar(&flags.port, "port", 8080, "specifies the port the http server runs on")
	flag.StringVar(&flags.socket, "socket", "", "specifies a socket to listen on, takes priority over -port")
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
	if flags.socket == "" {
		log.Printf("Serving on :%d\n", flags.port)
		log.Fatal(mux.ListenAndServe(fmt.Sprintf(":%d", flags.port)))
	} else {
		log.Printf("Serving on socket %s\n", flags.socket)
		log.Fatal(mux.ServeToSocket(flags.socket))
	}
}
