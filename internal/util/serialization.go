package util

import (
	"encoding/json"
	"io"
	"reflect"
	"slices"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/copier"
)

func subtype(t reflect.Type, group string) reflect.Type {
	// If the passed item is a slice, we need to unwrap it to get the contained type.
	if t.Kind() == reflect.Slice {
		return reflect.SliceOf(subtype(t.Elem(), group))
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
			if slices.Contains(props, group) {
				// Recursively subtype the property
				field.Type = subtype(field.Type, group)

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

func toEmptyMask(baseType reflect.Type, group string) interface{} {
	mask := baseType
	if group != "-" {
		mask = subtype(baseType, group)
	}

	return reflect.New(mask).Interface()
}

// Deserialize takes JSON data from an io.Reader and transforms it into a type K. During this process,
// it utilizes the "groups" tag to optionally filter out disallowed properties and uses the validate
// library to validate all properties.
func Deserialize[K interface{}](content io.Reader, group string) (error, K) {
	var target K

	// Trim down to mask
	mask := toEmptyMask(reflect.TypeFor[K](), group)

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
func Serialize(content interface{}, group string) ([]byte, error) {
	// Get masked struct
	mask := toEmptyMask(reflect.TypeOf(content), group)

	// Populate masked
	err := copier.CopyWithOption(mask, content, copier.Option{DeepCopy: true})
	if err != nil {
		return []byte{}, err
	}

	return json.Marshal(mask)
}
