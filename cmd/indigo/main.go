//go:debug x509negativeserial=1
package main

import (
	"log"
	"net/http"

	c "myoptions.info/indigo/backend/internal/config"
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

	// Create routes
	mux := c.CreateMux(
		c.Services{
			Config: config,
			Ldap:   &l,
			Jwt:    &jwtTransformer,
		},
	)

	// Serve
	log.Printf("Serving on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
