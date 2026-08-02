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
)

// --- Generic error message keys ---
const (
	MsgKeyResourceNotFound MsgKey = "error.resource_not_found"
	MsgKeyBadRequest       MsgKey = "error.bad_request"
	MsgKeyInternalError    MsgKey = "error.internal"
	MsgKeyResourceConflict MsgKey = "error.resource_conflict"
	MsgKeyUnauthorized     MsgKey = "error.unauthorized"
	MsgKeyForbidden        MsgKey = "error.forbidden"
	MsgKeyValidationFailed MsgKey = "error.validation_failed"
	MsgKeyInvalidIDFormat  MsgKey = "error.invalid_id_format"
)

// --- Auth message keys ---
const (
	MsgKeyMissingAuthHeader   MsgKey = "auth.missing_header"
	MsgKeyInvalidAuthFormat   MsgKey = "auth.invalid_format"
	MsgKeyInvalidToken        MsgKey = "auth.invalid_token"
	MsgKeyInvalidCredentials  MsgKey = "auth.invalid_credentials"
	MsgKeyFailedHashPassword  MsgKey = "auth.failed_hash_password"
	MsgKeyFailedRegister      MsgKey = "auth.failed_register"
	MsgKeyFailedGenerateToken MsgKey = "auth.failed_generate_token"
)

// --- User message keys ---
const (
	MsgKeyUserNotFound      MsgKey = "user.not_found"
	MsgKeyUserEmailConflict MsgKey = "user.email_conflict"
	MsgKeyUserPhoneConflict MsgKey = "user.phone_conflict"
	MsgKeyFailedCreateUser  MsgKey = "user.failed_create"
	MsgKeyFailedFindUser    MsgKey = "user.failed_find"
	MsgKeyFailedFetchUsers  MsgKey = "user.failed_fetch"
	MsgKeyFailedUpdateUser  MsgKey = "user.failed_update"
	MsgKeyFailedDeleteUser  MsgKey = "user.failed_delete"
)

// --- Validator custom message keys ---
const (
	MsgKeyPhone          MsgKey = "validator.phone"
	MsgKeyStrongPassword MsgKey = "validator.strong_password"
	MsgKeyFallback       MsgKey = "validator.fallback"
)
