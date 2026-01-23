package util

import (
	"encoding/json"
	"fmt"
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

func subtype(source interface{}, group string) reflect.Type {
	t := reflect.TypeOf(source)
	var fields []reflect.StructField

	// Extract any fields holding the desired group
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tag, present := field.Tag.Lookup("groups")
		if present {
			if slices.Contains(strings.Split(tag, ","), group) {
				fields = append(fields, field)
			}
		}
	}
	// TODO: Cascade?

	// Ensure the array exists
	if fields == nil {
		return t
	}

	// Create new struct
	return reflect.StructOf(fields)
}

func Deserialize[K interface{}](content io.Reader, group string) (error, K) {
	// Get proper type
	var maskedType reflect.Type
	var target K
	if group == "-" {
		maskedType = reflect.TypeOf(target)
	} else {
		maskedType = subtype(target, group)
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
		fmt.Println(err)
		return err, target
	}

	// Copy into target
	err := copier.CopyWithOption(&target, masked, copier.Option{DeepCopy: true})
	return err, target
}
