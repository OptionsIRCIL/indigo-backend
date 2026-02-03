package config

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/model/entity"
	"myoptions.info/indigo/backend/model/schema"
)
import c "myoptions.info/indigo/backend/internal/controller"
import s "myoptions.info/indigo/backend/internal/service"

// registerRouterNode recursively registers a routerNode and its children to a mux.
// It also adds an OPTIONS method to support CORS preflight requests.
func (m *MuxWrapper) registerRouterNode(node RouterConfig, parentPath string) {
	// Concat the parent path with this node's path for fully-qualified path
	path := util.PathConcat(parentPath, node.Path)

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
	m.registerRouterNode(m.Routes, "")
}

// CreateMux takes in a [Services] struct, and using the contained utilities, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateMux(services Services) MuxWrapper {
	auth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireAuth(services.Jwt, services.Ldap, next)
	}

	mux := MuxWrapper{
		mux:      http.NewServeMux(),
		services: services,
		Routes: RouterConfig{
			Path: "/",
			Methods: []MethodConfig{
				{
					Method:  "GET",
					Summary: `Demo "Hello World" endpoint.`,
					Handler: auth(c.IndexHelloWorld),
				},
			},
			Children: []RouterConfig{
				{
					Path: "/auth",
					Methods: []MethodConfig{
						{
							Method:  "POST",
							Summary: "Authenticate against LDAP to request an authentication token.",
							InputDto: &DataTransferObject{
								Interface: schema.LoginCredentials{},
							},
							Responses: map[int]Response{
								204: {
									Description: "Successful Authentication",
								},
								422: {
									Description: "Invalid Credentials",
									Dto: &DataTransferObject{
										Interface: util.HttpError{},
									},
								},
							},
							Handler: c.AuthEntry(services.Config, services.Ldap, services.Jwt, services.Flags.AuthSameSite),
						},
						{
							Method:  "DELETE",
							Summary: "Delete current cookie",
							Handler: c.DeleteCookie(services.Config),
							Responses: map[int]Response{
								204: {
									Description: "yeag",
								},
							},
						},
					},
				},
				{
					Path: "/person",
					Methods: []MethodConfig{
						{
							Method:  "POST",
							Summary: "Create a new person",
							InputDto: &DataTransferObject{
								Interface: entity.Person{},
								Groups:    []string{"post"},
							},
							Handler: auth(c.PrimitivePost[entity.Person](
								services.Database,
								c.SerializationParameters{
									SerializationGroup:   []string{"post"},
									DeserializationGroup: []string{"get"},
								},
							)),
							Responses: map[int]Response{
								201: {
									Description: "Person successfully created",
									Dto: &DataTransferObject{
										Interface: entity.Person{},
										Groups:    []string{"get"},
									},
								},
								422: {
									Description: "Serialization error",
									Dto: &DataTransferObject{
										Interface: util.HttpError{},
									},
								},
							},
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

// ServeToSocket wraps http.Serve and creates a net.Listen listener under the unix network type.
// Also sets the ownership of the socket with os.Chown using uid:gid. If an ownership change is
// not desired, these can be set to -1 (flag default). Default permissions is 770.
func (m *MuxWrapper) ServeToSocket(socket string, uid int, gid int) error {
	listener, err := net.Listen("unix", socket)
	if err != nil {
		return err
	}

	// Set ownership, if applicable
	err = os.Chown(socket, uid, gid)
	if err != nil {
		return err
	}

	// Allow rwx for owner and group, nothing for everyone else
	err = os.Chmod(socket, 0770)
	if err != nil {
		return err
	}

	return http.Serve(listener, m.mux)
}

// Services is an aggregation of various tools that will be made available during assembly
// of a routerConfig.
type Services struct {
	Config   *util.Config
	Ldap     *s.LdapConnection
	Jwt      *s.JwtTransformer
	Flags    util.ServeRuntimeFlags
	Database *gorm.DB
}

// MuxWrapper ideally wraps around an [http.ServeMux] to abstract away some common middleware or Routes
// such as logging, user authentication, or CORS headers.
type MuxWrapper struct {
	mux      *http.ServeMux
	services Services
	Routes   RouterConfig
}

type DataTransferObject struct {
	Interface interface{}
	Groups    []string
}

type Response struct {
	Description string
	Dto         *DataTransferObject
}

// MethodConfig defines the behavior that a mux should follow for a Method invoked on a given route.
type MethodConfig struct {
	Method    string
	Summary   string
	InputDto  *DataTransferObject
	Responses map[int]Response
	Handler   http.HandlerFunc
}

// RouterConfig defines each route added to the application.
// TODO: Middleware?
type RouterConfig struct {
	Path     string         `json:"path"`               // The Path to assign methods to.
	Methods  []MethodConfig `json:"methods"`            // What to do for each available HTTP method.
	Children []RouterConfig `json:"children,omitempty"` // Each child will inherit the parent's Path.
}
