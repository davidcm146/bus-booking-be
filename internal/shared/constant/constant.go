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

// --- Business time zone (Vietnam has no DST, so a fixed UTC+7 offset is exact) ---
const (
	BusinessTimeZoneName   = "Asia/Ho_Chi_Minh"
	BusinessTimeZoneOffset = 7 * 60 * 60
	DateLayout             = "2006-01-02"
)

// --- Trip search ---
const (
	DefaultTripPageSize = 25
	MaxTripPageSize     = 100
)

const (
	PgCodeNotNullViolation    = "23502"
	PgCodeForeignKeyViolation = "23503"
	PgCodeUniqueViolation     = "23505"
	PgCodeCheckViolation      = "23514"
	PgCodeInvalidText         = "22P02"
)
