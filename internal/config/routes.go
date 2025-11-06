package config

import (
	"fmt"
	"net/http"

	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)
import c "myoptions.info/indigo/backend/internal/controller"

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
func CreateRoutes(config *util.Config, ldap *service.LdapConnection, jwtTransformer *service.JwtTransformer) *http.ServeMux {
	mux := muxWrapper{
		mux:    http.NewServeMux(),
		config: config,
		ldap:   ldap,
		jwt:    jwtTransformer,
	}
	mux.HandleFunc("/api/v1/auth", []string{"POST"}, c.AuthEntry(config, ldap, jwtTransformer))
	mux.HandleFunc("/api/v1", []string{"GET"}, middleware.RequireAuth(jwtTransformer, ldap, c.IndexHelloWorld))
	return mux.mux
}
