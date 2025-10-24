//go:debug x509negativeserial=1
package main

import (
	c "backend/internal/config"
	s "backend/internal/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	log.Print("Hello, World!")

	// Load environment
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Populate LdapConnection
	l := s.LdapConnection{
		Base: os.Getenv("LDAP_BASE_DN"),
	}
	l.SetUrl(os.Getenv("LDAP_URL"))

	// Initialize connection
	err = l.Initialize(
		os.Getenv("LDAP_USERNAME"),
		os.Getenv("LDAP_PASSWORD"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Initialize JWT & secret
	jwtTransformer := s.JwtTransformer{}
	jwtInitErr := jwtTransformer.SetSecret([]byte(os.Getenv("APP_SECRET")))
	if jwtInitErr != nil {
		log.Fatal(jwtInitErr)
	}

	// Create routes
	mux := c.CreateRoutes(&l, &jwtTransformer)

	// Serve
	fmt.Printf("Serving on :8080\n")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
