package routine

import (
	"encoding/json"
	"fmt"
	"maps"
	"reflect"
	"slices"
	"strconv"
	"strings"

	c "myoptions.info/indigo/backend/internal/config/routes"
	"myoptions.info/indigo/backend/internal/schema/openApi"
	s "myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
)

func getDtoName(dto *c.DataTransferObject) string {
	if dto == nil {
		return "nil"
	}

	reflection := reflect.TypeOf(dto.Interface)
	if reflection == nil {
		return "unknown"
	}

	name := reflection.Name()
	if reflection.Kind() == reflect.Slice {
		name = "array_" + reflection.Elem().Name()
	}

	if dto.Groups == nil || len(dto.Groups) == 0 {
		return name
	}

	sorted := slices.Clone(dto.Groups)
	slices.Sort(sorted)
	suffix := strings.Join(sorted, ".")

	return name + "." + suffix
}

func routerConfigToMethodsElement(config *c.RouterConfig, path string, schemata map[string]openApi.SchemaType) map[string]openApi.Method {
	methods := make(map[string]openApi.Method)

	for _, method := range config.Methods {
		doc := openApi.Method{
			Summary:     path,
			Description: method.Summary,
			Responses:   make(map[string]openApi.Content),
		}

		if method.InputDto != nil {
			name := getDtoName(method.InputDto)

			// Add request body
			doc.RequestBody = &openApi.Content{
				Content: map[string]openApi.MediaType{
					"application/json": {
						Schema: openApi.SchemaType{
							Reference: "#/components/schemas/" + name,
						},
					},
				},
			}

			// Add schema (if applicable)
			if _, exists := schemata[name]; !exists {
				schemata[name] = util.ToOpenApiSchema(method.InputDto.Interface, method.InputDto.Groups)
			}
		}

		if method.IsAttachment {
			if method.Method == "POST" || method.Method == "PUT" || method.Method == "PATCH" {
				doc.RequestBody = &openApi.Content{
					Content: map[string]openApi.MediaType{
						"multipart/form-data": {
							Schema: openApi.SchemaType{
								Type: "object",
								Properties: map[string]openApi.SchemaType{
									"attachment": {
										Type:   "object",
										Format: "binary",
									},
								},
							},
						},
					},
				}
			}
		}

		for code, response := range method.Responses {
			responseContent := map[string]openApi.MediaType{}
			if response.Dto != nil {
				name := getDtoName(response.Dto)

				// Add response body
				responseContent["application/json"] = openApi.MediaType{
					Schema: openApi.SchemaType{
						Reference: "#/components/schemas/" + name,
					},
				}

				// Add schema (if applicable)
				if _, exists := schemata[name]; !exists {
					schemata[name] = util.ToOpenApiSchema(response.Dto.Interface, response.Dto.Groups)
				}
			}

			if response.IsAttachment {
				responseContent["application/octet-stream"] = openApi.MediaType{
					Schema: openApi.SchemaType{
						Type:   "string",
						Format: "binary",
					},
				}
			}

			doc.Responses[strconv.Itoa(code)] = openApi.Content{
				Description: response.Description,
				Content:     responseContent,
			}
		}

		methods[strings.ToLower(method.Method)] = doc
	}

	return methods
}

func routerConfigWalk(node *c.RouterConfig, base string, schemata map[string]openApi.SchemaType) map[string]map[string]openApi.Method {
	nodes := make(map[string]map[string]openApi.Method)
	path := util.PathConcat(base, node.Path)

	for _, sub := range node.PathValueSubstitutions {
		path = strings.Replace(path, "{"+sub.Original+"}", "{"+sub.New+"}", 1)
	}

	nodes[path] = routerConfigToMethodsElement(node, path, schemata)

	for _, child := range node.Children {
		maps.Copy(nodes, routerConfigWalk(&child, path, schemata))
	}

	return nodes
}

func muxToOpenApi(mux *c.MuxWrapper) openApi.OpenApi {
	api := openApi.OpenApi{
		Version: "3.1.0",
		Info: openApi.ApiInfo{
			Title:       "Indigo REST API",
			Description: "Documentation for interacting with OptionsIRCIL/indigo-backend over HTTP.",
			Version:     "0.0.0",
		},
		Components: openApi.Components{
			Schemas: map[string]openApi.SchemaType{},
		},
	}

	// Walk config
	api.Paths = routerConfigWalk(&mux.Routes, "", api.Components.Schemas)

	return api
}

// RunGenerateOpenApiSpec initializes some dummy services and dumps all configured routes from config.CreateMux.
// Utilized during testing to verify that all configured routes are documented in the OpenAPI spec.
// TODO Update description
func RunGenerateOpenApiSpec() int {
	// Create routes using MuxWrapper. Provide dummy services.
	mux := c.CreateMux(
		c.Services{
			Ldap: &s.LdapConnection{},
		},
	)

	output, _ := json.MarshalIndent(muxToOpenApi(&mux), "", "  ")
	fmt.Println(string(output[:]))

	return 0
}
