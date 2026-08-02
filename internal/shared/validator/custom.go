package validator

import (
	"regexp"
	"unicode"

	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	govalidator "github.com/go-playground/validator/v10"
)

// Holds a custom validator function and its i18n message key.
type customEntry struct {
	fn         govalidator.Func
	messageKey constant.MsgKey
}

// Stores custom validators registered before Setup() is called.
var registry []struct {
	tag   string
	entry customEntry
}

// customMessages stores custom message keys keyed by validator tag.
// Populated during Setup() from the registry.
var customMessages = map[string]constant.MsgKey{}

// RegisterCustom queues a custom validator to be registered during Setup().
// tag is the binding tag name, fn is the validation function,
// and key is the i18n message key for the error template.
//
// Example:
//
//	validator.RegisterCustom("vn_phone", func(fl validator.FieldLevel) bool {
//	    phone := fl.Field().String()
//	    return strings.HasPrefix(phone, "+84") || strings.HasPrefix(phone, "0")
//	}, constant.MsgKeyPhone)
func RegisterCustom(tag string, fn govalidator.Func, key constant.MsgKey) {
	registry = append(registry, struct {
		tag   string
		entry customEntry
	}{
		tag:   tag,
		entry: customEntry{fn: fn, messageKey: key},
	})
}

// phoneRegex matches phone numbers: optional + prefix, 9-15 digits.
var phoneRegex = regexp.MustCompile(constant.PhoneRegexPattern)

// validatePhone checks if the field value is a valid phone number format.
func validatePhone(fl govalidator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// minimum 8 characters, at least one uppercase, one lowercase, and one digit.
func validatePassword(fl govalidator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < constant.MinPasswordLength {
		return false
	}

	var hasUpper, hasLower, hasDigit bool
	for _, ch := range password {
		switch {
		case unicode.IsUpper(ch):
			hasUpper = true
		case unicode.IsLower(ch):
			hasLower = true
		case unicode.IsDigit(ch):
			hasDigit = true
		}
	}

	return hasUpper && hasLower && hasDigit
}

// init registers built-in custom validators.
func init() {
	RegisterCustom(constant.TagPhone, validatePhone, constant.MsgKeyPhone)
	RegisterCustom(constant.TagStrongPassword, validatePassword, constant.MsgKeyStrongPassword)
}
