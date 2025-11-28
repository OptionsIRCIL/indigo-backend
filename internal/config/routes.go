package config

import (
	"encoding/json"
	"fmt"
	"net/http"

	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/util"
)
import c "myoptions.info/indigo/backend/internal/controller"
import s "myoptions.info/indigo/backend/internal/service"

// registerRouterNode recursively registers a routerNode and its children to a mux.
// It also adds an OPTIONS method to support CORS preflight requests.
func (m *MuxWrapper) registerRouterNode(node RouterConfig, parentPath string) {
	// Concat the parent path with this node's path for fully-qualified path
	path := parentPath
	if parentPath != "" {
		if path[len(path)-1] == '/' && node.Path[0] == '/' {
			path = path + node.Path[1:]
		} else if path[len(path)-1] != '/' && node.Path[0] != '/' {
			path = path + "/" + node.Path
		} else {
			path = path + node.Path
		}
	} else {
		path = node.Path
	}

	// Register path to declared operations in router and collect list of methods
	methods := make([]string, len(node.Methods))
	for i, operation := range node.Methods {
		methods[i] = operation.Method
		m.mux.HandleFunc(
			fmt.Sprintf("%s %s", operation.Method, path),
			operation.Handler,
		)
	}

	// Add OPTIONS method for CORS preflight
	m.mux.HandleFunc(fmt.Sprintf("OPTIONS %s", path), c.ProvideOptions(methods))

	// Register children
	for _, child := range node.Children {
		m.registerRouterNode(child, path)
	}
}

// Initialize begins adding the root route and its children to the mux.
func (m *MuxWrapper) Initialize() {
	m.registerRouterNode(m.routes, "")
}

// CreateMux takes in a [Services] struct, and using the contained utilities, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateMux(services Services) MuxWrapper {
	mux := MuxWrapper{
		mux:      http.NewServeMux(),
		services: services,
		routes: RouterConfig{
			Path: "/",
			Methods: []methodConfig{
				{
					Method:  "GET",
					Handler: middleware.RequireAuth(services.Jwt, services.Ldap, c.IndexHelloWorld),
				},
			},
			Children: []RouterConfig{
				{
					Path: "/auth",
					Methods: []methodConfig{
						{
							Method:  "POST",
							Handler: c.AuthEntry(services.Config, services.Ldap, services.Jwt),
						},
						{
							Method:  "DELETE",
							Handler: c.DeleteCookie(services.Config),
						},
					},
				},
			},
		},
	}

	mux.Initialize()

	return mux
}

// ListenAndServe wraps http.ListenAndServe.
func (m *MuxWrapper) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, m.mux)
}

// DumpRoutes dumps the routing config to JSON. Useful for verifying documentation
// in CI jobs.
func (m *MuxWrapper) DumpRoutes() string {
	encoded, _ := json.MarshalIndent(m.routes, "", "  ")
	return string(encoded)
}

// Services is an aggregation of various tools that will be made available during assembly
// of a routerConfig.
type Services struct {
	Config *util.Config
	Ldap   *s.LdapConnection
	Jwt    *s.JwtTransformer
}

// MuxWrapper ideally wraps around an [http.ServeMux] to abstract away some common middleware or routes
// such as logging, user authentication, or CORS headers.
type MuxWrapper struct {
	mux      *http.ServeMux
	services Services
	routes   RouterConfig
}

// methodConfig defines the behavior that a mux should follow for a Method invoked on a given route.
type methodConfig struct {
	Method  string
	Handler http.HandlerFunc
}

// MarshalJSON abstracts away a methodConfig to only return the method string
// when being printed by the JSON marshaller.
func (m *methodConfig) MarshalJSON() ([]byte, error) {
	return []byte("\"" + m.Method + "\""), nil
}

// RouterConfig defines each route added to the application.
// TODO: Middleware?
type RouterConfig struct {
	Path     string         `json:"path"`               // The Path to assign methods to.
	Methods  []methodConfig `json:"methods"`            // What to do for each available HTTP method.
	Children []RouterConfig `json:"children,omitempty"` // Each child will inherit the parent's Path.
}
