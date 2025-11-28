package config

import (
	"fmt"
	"net/http"

	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/util"
)
import c "myoptions.info/indigo/backend/internal/controller"
import s "myoptions.info/indigo/backend/internal/service"

// registerRouterNode recursively registers a routerNode and its children to a mux.
// It also adds an OPTIONS method to support CORS preflight requests.
func (m *muxWrapper) registerRouterNode(node routerConfig, parentPath string) {
	// Concat the parent path with this node's path for fully-qualified path
	path := parentPath + node.path

	// Register path to declared operations in router and collect list of methods
	methods := make([]string, len(node.methods))
	for i, operation := range node.methods {
		methods[i] = operation.method
		m.mux.HandleFunc(
			fmt.Sprintf("%s %s", operation.method, path),
			operation.handler,
		)
	}

	// Add OPTIONS method for CORS preflight
	m.mux.HandleFunc(fmt.Sprintf("OPTIONS %s", path), c.ProvideOptions(methods))

	// Register children
	for _, child := range node.children {
		m.registerRouterNode(child, path)
	}
}

// Initialize begins adding the root route and its children to the mux.
func (m *muxWrapper) Initialize() {
	m.registerRouterNode(m.routes, "")
}

// CreateMux takes in a [Services] struct, and using the contained utilities, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateMux(services Services) *http.ServeMux {
	mux := muxWrapper{
		mux:      http.NewServeMux(),
		services: services,
		routes: routerConfig{
			path: "/api/v1",
			methods: []methodConfig{
				{
					method:  "GET",
					handler: middleware.RequireAuth(services.Jwt, services.Ldap, c.IndexHelloWorld),
				},
			},
			children: []routerConfig{
				{
					path: "/auth",
					methods: []methodConfig{
						{
							method:  "POST",
							handler: c.AuthEntry(services.Config, services.Ldap, services.Jwt),
						},
						{
							method:  "DELETE",
							handler: c.DeleteCookie(services.Config),
						},
					},
				},
			},
		},
	}

	mux.Initialize()

	return mux.mux
}

// Services is an aggregation of various tools that will be made available during assembly
// of a routerConfig.
type Services struct {
	Config *util.Config
	Ldap   *s.LdapConnection
	Jwt    *s.JwtTransformer
}

// muxWrapper ideally wraps around an [http.ServeMux] to abstract away some common middleware or routes
// such as logging, user authentication, or CORS headers.
// TODO: fully wrap mux?
type muxWrapper struct {
	mux      *http.ServeMux
	services Services
	routes   routerConfig
}

// methodConfig defines the behavior that a mux should follow for a method invoked on a given route.
type methodConfig struct {
	method  string
	handler http.HandlerFunc
}

// routerConfig defines each route added to the application.
// TODO: Middleware?
type routerConfig struct {
	path     string         // The path to assign methods to.
	methods  []methodConfig // What to do for each available HTTP method.
	children []routerConfig // Each child will inherit the parent's path.
}
