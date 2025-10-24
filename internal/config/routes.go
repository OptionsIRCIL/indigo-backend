package config

import (
	"backend/internal/service"
	"net/http"
)
import c "backend/internal/controller"
import m "backend/internal/middleware"

func CreateRoutes(ldap *service.LdapConnection, jwtTransformer *service.JwtTransformer) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/auth", c.AuthEntry(ldap, jwtTransformer))
	mux.HandleFunc("GET /", m.RequireAuth(jwtTransformer, c.IndexHelloWorld))
	return mux
}
