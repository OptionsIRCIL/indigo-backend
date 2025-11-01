package config

import (
	"net/http"

	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)
import c "myoptions.info/indigo/backend/internal/controller"
import m "myoptions.info/indigo/backend/internal/middleware"

// CreateRoutes takes in a [util.Config] and various service structs, and using them, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateRoutes(config *util.Config, ldap *service.LdapConnection, jwtTransformer *service.JwtTransformer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth", c.AuthEntry(config, ldap, jwtTransformer))
	mux.HandleFunc("GET /", m.RequireAuth(jwtTransformer, c.IndexHelloWorld))
	return mux
}
