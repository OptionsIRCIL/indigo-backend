package config

import (
	"fmt"
	"log"
	"net/http"

	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)
import c "myoptions.info/indigo/backend/internal/controller"

// Services Struct meant to pass all app dependencies from main.go,
type Services struct {
	Config *util.Config
	Ldap   *service.LdapConnection
	Jwt    *service.JwtTransformer
}

// muxWrapper ideally wraps around an [http.ServeMux] to abstract away some common middleware or routes
// such as logging, user authentication, or CORS headers.
type muxWrapper struct {
	mux    *http.ServeMux
	config *util.Config
	ldap   *service.LdapConnection
	jwt    *service.JwtTransformer
}

// HandleFunc adds multiple [http.HandleFunc] functions to the mux wrapped by a muxWrapper.
// It wraps each with common middleware and adds an OPTIONS method to ensure CORS preflight functionality.
func (m *muxWrapper) HandleFunc(route string, methods []string, handler http.HandlerFunc) {
	// Add OPTIONS method
	m.mux.HandleFunc(fmt.Sprintf("OPTIONS %s", route), c.ProvideOptions(methods))

	// Add other things I guess
	for _, method := range methods {
		m.mux.HandleFunc(fmt.Sprintf("%s %s", method, route), handler)
	}
}

// CreateRoutes takes in a [util.Config] and various service structs, and using them, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateRoutes(s Services) *muxWrapper {
	mux := muxWrapper{
		mux:    http.NewServeMux(),
		config: s.Config,
		ldap:   s.Ldap,
		jwt:    s.Jwt,
	}
	mux.HandleFunc("/api/v1/auth", []string{"POST"}, c.AuthEntry(s.Config, s.Ldap, s.Jwt))
	mux.HandleFunc("/api/v1", []string{"GET"}, middleware.RequireAuth(s.Jwt, s.Ldap, c.IndexHelloWorld))
	return &mux
}

func (m *muxWrapper) ListenAndServe(addr string) error {
	log.Printf("Starting HTTP server on %s using MuxWrapper...", addr)
	// This is the essential blocking call that starts the web server.
	return http.ListenAndServe(addr, m.mux)
}
