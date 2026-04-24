package validator

import (
	"reflect"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func Phone(fieldLevel validator.FieldLevel) bool {
	field := fieldLevel.Field()
	if field.Kind() != reflect.String {
		return field.IsValid()
	}

	expr := regexp.MustCompile("^\\+?[0-9/()\\-. ]{7,18}(?:x ?[0-9]{1,4})?$")
	return expr.MatchString(field.String())
}
