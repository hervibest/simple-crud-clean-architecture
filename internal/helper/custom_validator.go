package helper

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type CustomValidator interface {
	Validate(any) []ValidationError
}

type customValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator(viper *viper.Viper) CustomValidator {
	validate := validator.New()
	return &customValidator{Validator: validate}
}

type ValidationError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func (cv *customValidator) Validate(payload any) []ValidationError {
	var validationErrors []ValidationError

	err := cv.Validator.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, ValidationError{
				Field:   err.Field(),
				Rule:    err.Tag(),
				Message: getErrorMessage(err),
			})
		}
	}

	return validationErrors
}

func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", err.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", err.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", err.Field(), err.Param())
	default:
		return fmt.Sprintf("%s is invalid", err.Field())
	}
}
