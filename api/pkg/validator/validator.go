package validator

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"unicode"
)

type CustomValidator struct {
	validator *validator.Validate
}

type CustomValidatorInterface interface {
	Validate(echo.Context, interface{}) interface{}
}

func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}
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

func errorWrapper(err error) interface{} {
	var (
		invalidValidationError *validator.InvalidValidationError
		httpErr                *echo.HTTPError
		validationErrors       []string
	)

	if errors.As(err, &invalidValidationError) {
		return map[string]interface{}{
			"error": "validation error",
		}
	}

	if errors.As(err, &httpErr) {
		internalErr := httpErr.Internal.(*json.UnmarshalTypeError)
		return errors.New(fmt.Sprintf(
			"field %v expected %v got %v",
			internalErr.Field,
			internalErr.Type,
			internalErr.Value,
		))
	}

	for _, err := range err.(validator.ValidationErrors) {
		validationErrors = append(validationErrors, formatError(err))
	}
	return map[string]interface{}{
		"error": validationErrors,
	}
}

func (v *CustomValidator) Validate(c echo.Context, data interface{}) interface{} {
	if err := c.Bind(&data); err != nil {
		return errorWrapper(err)
	}
	if err := v.validator.Struct(data); err != nil {
		return errorWrapper(err)
	}
	return nil
}