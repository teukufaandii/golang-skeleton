package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()

	// Use JSON field names in error messages
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validators
	_ = validate.RegisterValidation("password", validatePassword)
	_ = validate.RegisterValidation("player_id", validatePlayerID)
}

// Validate validates a struct and returns formatted errors
func Validate(data interface{}) error {
	err := validate.Struct(data)
	if err == nil {
		return nil
	}

	var errors []string
	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, formatValidationError(err))
	}

	return fmt.Errorf("%s", strings.Join(errors, "; "))
}

// formatValidationError formats a single validation error
func formatValidationError(err validator.FieldError) string {
	field := err.Field()

	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, err.Param())
	case "max":
		return fmt.Sprintf("%s must not exceed %s characters", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, err.Param())
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "password":
		return fmt.Sprintf("%s must contain at least one uppercase, lowercase, number, and special character", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// validatePassword checks password strength
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var hasUpper, hasLower, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case 'A' <= char && char <= 'Z':
			hasUpper = true
		case 'a' <= char && char <= 'z':
			hasLower = true
		case '0' <= char && char <= '9':
			hasNumber = true
		case strings.ContainsRune("!@#$%^&*()_+-=[]{}|;:',.<>?/~`", char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasNumber && hasSpecial
}

// validatePlayerID validates game player ID format
func validatePlayerID(fl validator.FieldLevel) bool {
	playerID := fl.Field().String()
	// Basic validation - alphanumeric only
	for _, char := range playerID {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9')) {
			return false
		}
	}
	return len(playerID) >= 3 && len(playerID) <= 100
}
