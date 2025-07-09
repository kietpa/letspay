package helper

import (
	"fmt"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (v *ValidationError) Error() string {
	return "Field " + v.Field + ": " + v.Message
}

func ValidateStruct(data any, validate validator.Validate) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(data)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Message = getErrorMessage(err)
			errors = append(errors, element)
		}
	}

	return errors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "gte":
		return fmt.Sprintf("Must be greater than or equal to %s", err.Param())
	case "lte":
		return fmt.Sprintf("Must be less than or equal to %s", err.Param())
	case "min":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Must be at least %s characters long", err.Param())
		}
		return fmt.Sprintf("Must have at least %s items", err.Param())
	case "max":
		if err.Kind() == reflect.String {
			return fmt.Sprintf("Must be no more than %s characters long", err.Param())
		}
		return fmt.Sprintf("Must have no more than %s items", err.Param())
	default:
		return fmt.Sprintf("Field validation failed on '%s' tag", err.Tag())
	}
}
