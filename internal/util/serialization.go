package util

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
	"myoptions.info/indigo/backend/internal/schema/openApi"
)

func pascalToCamel(s string) string {
	if len(s) < 2 {
		return s
	}

	return strings.ToLower(string(s[0])) + s[1:]
}

func intersects[T comparable](a []T, b []T) bool {
	for _, j := range a {
		for _, k := range b {
			if j == k {
				return true
			}
		}
	}

	return false
}

func subtype(t reflect.Type, groups []string) reflect.Type {
	// If the passed item is a slice, we need to unwrap it to get the contained type.
	if t.Kind() == reflect.Slice {
		return reflect.SliceOf(subtype(t.Elem(), groups))
	}

	// We can only support if the slice is a struct type.
	if t.Kind() != reflect.Struct {
		return t
	}

	var fields []reflect.StructField

	// Extract any fields holding the desired group
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, present := field.Tag.Lookup("groups")
		if present {
			props := strings.Split(tag, ",")
			if intersects(props, groups) {
				// Ensure JSON key exists, else, create one in camelCase
				_, hasJson := field.Tag.Lookup("json")
				if !hasJson {
					field.Tag = reflect.StructTag(string(field.Tag) + ` json:"` + pascalToCamel(field.Name) + `"`)
				}

				// Recursively subtype the property
				field.Type = subtype(field.Type, groups)

				fields = append(fields, field)
			}
		}
	}

	// Ensure the array exists
	if fields == nil {
		return t
	}

	return reflect.StructOf(fields)
}

func toEmptyMask(baseType reflect.Type, groups []string) interface{} {
	mask := baseType
	if groups != nil && len(groups) != 0 {
		mask = subtype(baseType, groups)
	}

	return reflect.New(mask).Interface()
}

// Deserialize takes JSON data from an io.Reader and transforms it into a type K. During this process,
// it utilizes the "groups" tag to optionally filter out disallowed properties and uses the validate
// library to validate all properties.
func Deserialize[K interface{}](content io.Reader, deserializationGroups []string) (error, K) {
	var target K

	// Trim down to mask
	mask := toEmptyMask(reflect.TypeFor[K](), deserializationGroups)

	// Decode
	// TODO: Integrate with json.go
	decoder := json.NewDecoder(content)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&mask); err != nil {
		return err, target
	}

	// Validate
	// TODO: Explore options, Utilize caching by building onto struct?
	validate := validator.New()
	if err := validate.Struct(mask); err != nil {
		return err, target
	}

	// Copy into target
	err := copier.CopyWithOption(&target, mask, copier.Option{DeepCopy: true})
	return err, target
}

// Serialize takes any interface and a group name. It'll then filter out properties not
// tagged with the provided group, unless group is equal to "-". This filtered interface
// is then passed through json.Marshal and returned.
func Serialize(content interface{}, groups []string) ([]byte, error) {
	// Get masked struct
	mask := toEmptyMask(reflect.TypeOf(content), groups)

	// Populate masked
	err := copier.CopyWithOption(mask, content, copier.Option{DeepCopy: true})
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(mask)
}

func maskToOpenApiSchema(reflection reflect.Type) openApi.SchemaType {
	if reflection.Kind() == reflect.Pointer {
		reflection = reflection.Elem()
	}

	// Edge cases for objects that don't serialize nicely or need additional context
	switch reflection.PkgPath() + "/" + reflection.Name() {
	case "github.com/google/uuid/UUID":
		return openApi.SchemaType{
			Type:    "string",
			Format:  "UUIDv4",
			Example: "00000000-0000-0000-0000-000000000000",
		}
	case "time/Time":
		return openApi.SchemaType{
			Type:    "string",
			Format:  "RFC3339",
			Example: "2006-01-02T15:04:05Z07:00",
		}
	case "myoptions.info/indigo/backend/model/Date":
		return openApi.SchemaType{
			Type:    "string",
			Format:  "YYYY-MM-DD",
			Example: "1970-01-01",
		}
	}

	switch reflection.Kind() {
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return openApi.SchemaType{Type: "integer"}
	case reflect.Float32, reflect.Float64:
		return openApi.SchemaType{Type: "number"}
	case reflect.String:
		return openApi.SchemaType{Type: "string"}
	case reflect.Bool:
		return openApi.SchemaType{Type: "boolean"}
	case reflect.Slice:
		return openApi.SchemaType{Type: "array"}
	case reflect.Interface, reflect.Struct:
		// We'll need to iterate over each property
		properties := make(map[string]openApi.SchemaType)

		for i := range reflection.NumField() {
			field := reflection.Field(i)
			properties[pascalToCamel(field.Name)] = maskToOpenApiSchema(field.Type)
		}

		return openApi.SchemaType{
			Type:       "object",
			Properties: properties,
		}
	default:
	}

	return openApi.SchemaType{Type: "null"}
}

func ToOpenApiSchema(content interface{}, groups []string) openApi.SchemaType {
	// Get masked struct
	mask := toEmptyMask(reflect.TypeOf(content), groups)

	return maskToOpenApiSchema(reflect.TypeOf(mask))
}
