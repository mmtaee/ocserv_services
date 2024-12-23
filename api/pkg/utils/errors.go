package utils

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"unicode"
)

func formatSnakeCase(s string) string {
	var result []rune
	for i, char := range s {
		if unicode.IsUpper(char) && i > 0 {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(char))
	}
	return string(result)
}

func formatError(err validator.FieldError) string {
	field := formatSnakeCase(err.Field())
	switch err.Tag() {
	case "required":
		return field + " is required"
	case "min":
		return field + " must be at least " + err.Param() + " characters long"
	case "max":
		return field + " must be at most " + err.Param() + " characters long"
	default:
		return "Invalid input"
	}
}

func InvalidBodyError(err error) map[string]interface{} {
	var invalidValidationError *validator.InvalidValidationError
	if errors.As(err, &invalidValidationError) {
		return map[string]interface{}{
			"error": "validation error",
		}
	}

	var validationErrors []string
	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, formatError(err))
	}
	return map[string]interface{}{
		"error": validationErrors,
	}
}
