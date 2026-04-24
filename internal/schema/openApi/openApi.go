package openApi

type OpenApi struct {
	Version    string                       `json:"openapi"`
	Info       ApiInfo                      `json:"info"`
	Paths      map[string]map[string]Method `json:"paths"`
	Components Components                   `json:"components"`
}

type SchemaType struct {
	Type       []string              `json:"type,omitempty"`
	Properties map[string]SchemaType `json:"properties,omitempty"`
	Items      *SchemaType           `json:"items,omitempty"`
	Examples   []string              `json:"examples,omitempty"`
	Format     string                `json:"format,omitempty"`
	Required   []string              `json:"required,omitempty"`
	Default    string                `json:"default,omitempty"`
	Reference  string                `json:"$ref,omitempty"`
}

type Components struct {
	Schemas map[string]SchemaType `json:"schemas"`
}

type ApiInfo struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Version     string `json:"version"`
}

type Content struct {
	Description string `json:"description,omitempty"`
	// The key of this property is a mime-type
	Content map[string]MediaType `json:"content,omitempty"`
}

type MediaType struct {
	Schema SchemaType `json:"schema"`
}

type Method struct {
	Summary     string             `json:"summary"`
	Description string             `json:"description,omitempty"`
	RequestBody *Content           `json:"requestBody,omitempty"`
	Responses   map[string]Content `json:"responses,omitempty"`
}
