//go:debug x509negativeserial=1
package main

import (
	c "backend/internal/config"
	s "backend/internal/service"
	"backend/internal/util"
	"log"
	"net/http"
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

	// Initialize connection
	err := l.Initialize(
		config.LdapUsername,
		config.LdapPassword,
	)
	if err != nil {
		log.Fatal(err)
	}
	l.SetSecure(config.IndigoEnv == "prod")

	// Initialize JWT & secret
	jwtTransformer := s.JwtTransformer{}
	jwtInitErr := jwtTransformer.SetSecret([]byte(config.IndigoSecret))
	if jwtInitErr != nil {
		log.Fatal(jwtInitErr)
	}

	// Create routes
	mux := c.CreateRoutes(config, &l, &jwtTransformer)

	// Serve
	log.Printf("Serving on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
