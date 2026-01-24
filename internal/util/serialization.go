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

type SerializationError struct {
	Msg string
}

func (s *SerializationError) Error() string {
	return s.Msg
}

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

// Deserialize takes JSON data from an io.Reader and transforms it into a type K. During this process,
// it utilizes the "groups" tag to optionally filter out disallowed properties and uses the validate
// library to validate all properties.
func Deserialize[K interface{}](content io.Reader, group string) (error, K) {
	var target K

	// Get proper type
	var maskedType reflect.Type
	if group == "-" {
		maskedType = reflect.TypeOf(target)
	} else {
		maskedType = subtype(reflect.TypeOf(target), group)
	}
	masked := reflect.New(maskedType).Interface()

	// Decode
	decoder := json.NewDecoder(content)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&masked); err != nil {
		return err, target
	}

	// Validate
	// TODO: Explore options, Utilize caching by building onto struct?
	validate := validator.New()
	if err := validate.Struct(masked); err != nil {
		return err, target
	}

	// Copy into target
	err := copier.CopyWithOption(&target, masked, copier.Option{DeepCopy: true})
	return err, target
}
