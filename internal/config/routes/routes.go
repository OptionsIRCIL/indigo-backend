package routes

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"

	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/middleware"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/model/entity"
	"myoptions.info/indigo/backend/model/schema"
)
import c "myoptions.info/indigo/backend/internal/controller"
import s "myoptions.info/indigo/backend/internal/service"

// generates a RouterConfig subtree for any entity provided as [T]
func generateCrudRoutes[T interface{}](database *gorm.DB, path string) RouterConfig {
	return RouterConfig{
		Path: path,
		Methods: []MethodConfig{
			{Method: "GET", Handler: c.PrimitiveGetCollection[T](database, nil)},
			{Method: "POST", Handler: c.PrimitivePost[T](database, nil)},
		},
		Children: []RouterConfig{
			{
				Path: "/{id}",
				Methods: []MethodConfig{
					{Method: "GET", Handler: c.PrimitiveGetOne[T](database, nil)},
					{Method: "PUT", Handler: c.PrimitivePut[T](database, nil)},
					{Method: "DELETE", Handler: c.PrimitiveDelete[T](database)},
				},
			},
		},
	}
}

// registerRouterNode recursively registers a routerNode and its children to a mux.
// It also adds an OPTIONS method to support CORS preflight requests.
func (m *MuxWrapper) registerRouterNode(node RouterConfig, parentPath string) {
	// Concat the parent path with this node's path for fully-qualified path
	path := util.PathConcat(parentPath, node.Path)

	// Perform any needed substitutions
	for _, sub := range node.PathValueSubstitutions {
		path = strings.Replace(path, "{"+sub.Original+"}", "{"+sub.New+"}", 1)
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
	m.registerRouterNode(m.Routes, "")
}

// CreateMux takes in a [Services] struct, and using the contained utilities, constructs an [http.ServeMux]
// that aggregates the various handler functions used in the application.
func CreateMux(services Services) MuxWrapper {
	auth := func(next http.HandlerFunc) http.HandlerFunc {
		return middleware.RequireAuth(services.Ldap, next)
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
								200: {
									Description: "Successful Authentication",
									Dto: &DataTransferObject{
										Interface: entity.Employee{},
										Groups:    []string{"get"},
									},
								},
								422: {
									Description: "Invalid Credentials",
									Dto: &DataTransferObject{
										Interface: util.HttpError{},
									},
								},
							},
							Handler: c.AuthEntry(services.Ldap, services.Database, services.Flags.AuthSameSite),
						},
						{
							Method:  "DELETE",
							Summary: "Delete current cookie",
							Handler: c.DeleteCookie(),
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
							Method:  "GET",
							Summary: "Get all persons",
							Handler: auth(c.PrimitiveGetCollection[entity.Person](services.Database, []string{"get"})),
						},
						{
							Method:  "POST",
							Summary: "Create a new person",
							InputDto: &DataTransferObject{
								Interface: entity.Person{},
								Groups:    []string{"post"},
							},
							Handler: auth(c.PrimitivePost[entity.Person](
								services.Database,
								&c.SerializationParameters{
									SerializationGroup:   []string{"get"},
									DeserializationGroup: []string{"post"},
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
					Children: []RouterConfig{
						{
							Path: "/{id}",
							Methods: []MethodConfig{
								{
									Method:  "GET",
									Summary: "Get a single person",
									Handler: auth(c.PrimitiveGetOne[entity.Person](
										services.Database,
										[]string{"get"},
									)),
									Responses: map[int]Response{
										200: {
											Description: "Person found",
											Dto: &DataTransferObject{
												Interface: entity.Person{},
												Groups:    []string{"get"},
											},
										},
										404: {
											Description: "Not found",
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
				{
					Path: "/information-and-referral",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all Information and Referral entries",
							Handler: auth(c.PrimitiveGetCollection[entity.InformationAndReferral](services.Database, []string{"get"})),
						},
						{
							Method:  "POST",
							Summary: "Create a new Information and Referral entry",
							InputDto: &DataTransferObject{
								Interface: entity.InformationAndReferral{},
								Groups:    []string{"post"},
							},
							Handler: auth(c.PrimitivePost[entity.InformationAndReferral](
								services.Database,
								&c.SerializationParameters{
									SerializationGroup:   []string{"get"},
									DeserializationGroup: []string{"post"},
								},
							)),
							Responses: map[int]Response{
								201: {
									Description: "I&R successfully created",
									Dto: &DataTransferObject{
										Interface: entity.InformationAndReferral{},
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
					Children: []RouterConfig{
						{
							Path: "/{id}",
							Methods: []MethodConfig{
								{
									Method:  "GET",
									Summary: "Get an information and referral record",
									Handler: auth(c.PrimitiveGetOne[entity.InformationAndReferral](
										services.Database,
										[]string{"get"},
									)),
									Responses: map[int]Response{
										200: {
											Description: "I&R found",
											Dto: &DataTransferObject{
												Interface: entity.InformationAndReferral{},
												Groups:    []string{"get"},
											},
										},
										404: {
											Description: "Not found",
											Dto: &DataTransferObject{
												Interface: util.HttpError{},
											},
										},
									},
								},
							},
							Children: []RouterConfig{
								{
									Path: "/effort",
									PathValueSubstitutions: []PathValueSubstitution{
										{
											Original: "id",
											New:      "informationAndReferralId",
										},
									},
									Methods: []MethodConfig{
										{
											Method:  "POST",
											Summary: "Log effort to an Information and Referral record",
											InputDto: &DataTransferObject{
												Interface: entity.InformationAndReferralEffort{},
												Groups:    []string{"post"},
											},
											Handler: auth(c.InformationAndReferralEffortPost(services.Database)),
											Responses: map[int]Response{
												201: {
													Description: "Effort log successfully created",
													Dto: &DataTransferObject{
														Interface: entity.InformationAndReferralEffort{},
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
								{
									Path: "/attachment",
									PathValueSubstitutions: []PathValueSubstitution{
										{
											Original: "id",
											New:      "informationAndReferralId",
										},
									},
									Methods: []MethodConfig{
										{
											Method:       "POST",
											Summary:      "Create a new Information and Referral attachment",
											IsAttachment: true,
											Handler:      auth(c.InformationAndReferralAttachmentPost(services.Database)),
											Responses: map[int]Response{
												200: {
													Description: "Attachment successfully created",
													Dto: &DataTransferObject{
														Interface: entity.InformationAndReferralAttachment{},
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
					},
				},
				{
					Path: "/information-and-referral-attachment/{id}",
					Methods: []MethodConfig{
						{
							Method:  "HEAD",
							Summary: "Get attachment details",
							Handler: auth(c.InformationAndReferralAttachmentGet(services.Database)),
							Responses: map[int]Response{
								200: {
									Description: "File details",
								},
								404: {
									Description: "Not found",
									Dto: &DataTransferObject{
										Interface: util.HttpError{},
									},
								},
							},
						},
						{
							Method:  "GET",
							Summary: "Get attachment details",
							Handler: auth(c.InformationAndReferralAttachmentGet(services.Database)),
							Responses: map[int]Response{
								200: {
									Description: "File contents",
								},
								404: {
									Description: "Not found",
									Dto: &DataTransferObject{
										Interface: util.HttpError{},
									},
								},
							},
						},
					},
				},
				{
					Path: "/employee",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all employees",
							Handler: auth(c.PrimitiveGetCollection[entity.Employee](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/community-service-event",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all community service events",
							Handler: auth(c.PrimitiveGetCollection[entity.CommunityServiceEvent](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/consumer-service",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all consumer services",
							Handler: auth(c.PrimitiveGetCollection[entity.ConsumerService](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/services-offered",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all services offered",
							Handler: auth(c.PrimitiveGetCollection[entity.ServicesOffered](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/organization",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all organizations",
							Handler: auth(c.PrimitiveGetCollection[entity.Organization](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/place",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all places",
							Handler: auth(c.PrimitiveGetCollection[entity.Place](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/record-def",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all record definitions",
							Handler: auth(c.PrimitiveGetCollection[entity.RecordDef](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/address-phone",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all address and phone records",
							Handler: auth(c.PrimitiveGetCollection[entity.AddressPhone](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/alias",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all aliases",
							Handler: auth(c.PrimitiveGetCollection[entity.Alias](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/disability-information",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all disability info records",
							Handler: auth(c.PrimitiveGetCollection[entity.DisabilityInfo](services.Database, []string{"get"})),
						},
					},
				},
				{
					Path: "/goals",
					Methods: []MethodConfig{
						{
							Method:  "GET",
							Summary: "Get all goals",
							Handler: auth(c.PrimitiveGetCollection[entity.Goal](services.Database, []string{"get"})),
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
	Ldap     *s.LdapConnection
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
	Description  string
	IsAttachment bool
	Dto          *DataTransferObject
}

// MethodConfig defines the behavior that a mux should follow for a Method invoked on a given route.
type MethodConfig struct {
	Method       string
	Summary      string
	IsAttachment bool // If set, the content type is set to multipart/form-data with a single key "attachment"
	InputDto     *DataTransferObject
	Responses    map[int]Response
	Handler      http.HandlerFunc
}

// A PathValueSubstitution may be used to rewrite PathValue keys in parent routes. Particularly
// useful if using PathValue keys as foreign keys.
type PathValueSubstitution struct {
	Original string
	New      string
}

// RouterConfig defines each route added to the application.
// TODO: Middleware?
type RouterConfig struct {
	Path                   string                  // The Path to assign methods to.
	PathValueSubstitutions []PathValueSubstitution // Any substitutions to apply to previously specified PathValues.
	Methods                []MethodConfig          // What to do for each available HTTP method.
	Children               []RouterConfig          // Each child will inherit the parent's Path.
}
