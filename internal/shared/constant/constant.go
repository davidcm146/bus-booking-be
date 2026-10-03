package constant

// --- MsgKey is a typed key for i18n message lookup ---
type MsgKey string

// --- Locale represents a supported language ---
type Locale string

const (
	LocaleEN Locale = "en"
	LocaleVI Locale = "vi"
)

// DefaultLocale is the fallback locale when none is detected.
const DefaultLocale = LocaleEN

// --- Validator tags ---
const (
	TagPhone          = "phone"
	TagStrongPassword = "strong_password"
)

// --- Validator constraints ---
const (
	MinPasswordLength = 8
	PhoneMinDigits    = 10
	PhoneMaxDigits    = 11
	PhoneRegexPattern = `^\+?[0-9]{10,11}$`
)

// --- API error codes ---
const (
	CodeInternalError   = "INTERNAL_ERROR"
	CodeValidationError = "VALIDATION_ERROR"
	CodeBadRequestError = "BAD_REQUEST"
	CodeNotFoundError   = "NOT_FOUND"
	CodeConflictError   = "CONFLICT"
	CodeUnauthorized    = "UNAUTHORIZED"
	CodeForbidden       = "FORBIDDEN"
)
