package validator

import (
	"context"
	"unicode"

	"github.com/gin-gonic/gin/binding"
	govalidator "github.com/go-playground/validator/v10"
)

// Return field-error in structured format.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Setup registers all custom validators with gin's validator engine.
// Call this once during application bootstrap, before routes are served.
func Setup() {
	v, ok := binding.Validator.Engine().(*govalidator.Validate)
	if !ok {
		return
	}

	for _, r := range registry {
		v.RegisterValidation(r.tag, r.entry.fn)
		customMessages[r.tag] = r.entry.messageKey
	}
}

// Convert a binding/validation error into a structured slice of FieldError with localized messages.
func FormatError(ctx context.Context, err error) []FieldError {
	errs, ok := err.(govalidator.ValidationErrors)
	if !ok {
		return []FieldError{
			{Field: "_", Message: err.Error()},
		}
	}

	fieldErrors := make([]FieldError, 0, len(errs))
	for _, fe := range errs {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   toCamelCase(fe.Field()),
			Message: GetMessage(ctx, fe.Tag(), fe.Field(), fe.Param()),
		})
	}
	return fieldErrors
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}

	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
