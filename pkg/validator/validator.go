package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// ValidationError represents a single validation error
type ValidationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Value   string `json:"value,omitempty"`
	Message string `json:"message"`
}

// ValidateStruct validates a struct and returns validation errors
func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidationError
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			element.Message = getErrorMsg(err)
			errors = append(errors, element)
		}
	}

	return errors
}

// getErrorMsg returns a human-readable error message for a validation error
func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + " is required"
	case "email":
		return fe.Field() + " must be a valid email address"
	case "min":
		return fe.Field() + " must be at least " + fe.Param() + " characters"
	case "max":
		return fe.Field() + " must be at most " + fe.Param() + " characters"
	case "oneof":
		return fe.Field() + " must be one of: " + fe.Param()
	case "url":
		return fe.Field() + " must be a valid URL"
	case "numeric":
		return fe.Field() + " must be a number"
	case "gte":
		return fe.Field() + " must be greater than or equal to " + fe.Param()
	case "lte":
		return fe.Field() + " must be less than or equal to " + fe.Param()
	default:
		return fe.Field() + " is invalid"
	}
}

// GetValidator returns the validator instance if additional customization is needed
func GetValidator() *validator.Validate {
	return validate
}
