package apperror

import (
	"net/http"

	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
)

type AppError struct {
	Code       string          `json:"code"`
	MessageKey constant.MsgKey `json:"-"`
	HTTPStatus int             `json:"status"`
	cause      error
}

// WithCause returns a copy of the error with the underlying cause attached.
// Returning a copy keeps shared sentinels (ErrNotFound, ...) immutable.
func (e *AppError) WithCause(err error) *AppError {
	clone := *e
	clone.cause = err
	return &clone
}

// Unwrap lets errors.Is / errors.As reach the underlying cause.
func (e *AppError) Unwrap() error {
	return e.cause
}

// Error returns the message key, plus the cause when present (for logging).
func (e *AppError) Error() string {
	if e.cause != nil {
		return string(e.MessageKey) + ": " + e.cause.Error()
	}
	return string(e.MessageKey)
}

// New creates a new AppError with an i18n message key.
func New(code string, key constant.MsgKey, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		MessageKey: key,
		HTTPStatus: httpStatus,
	}
}

// --- Predefined sentinel errors ---

var (
	ErrNotFound     = New(constant.CodeNotFoundError, i18n.KeyErrorResourceNotFound, http.StatusNotFound)
	ErrBadRequest   = New(constant.CodeBadRequestError, i18n.KeyErrorBadRequest, http.StatusBadRequest)
	ErrInternal     = New(constant.CodeInternalError, i18n.KeyErrorInternal, http.StatusInternalServerError)
	ErrConflict     = New(constant.CodeConflictError, i18n.KeyErrorResourceConflict, http.StatusConflict)
	ErrUnauthorized = New(constant.CodeUnauthorized, i18n.KeyErrorUnauthorized, http.StatusUnauthorized)
	ErrForbidden    = New(constant.CodeForbidden, i18n.KeyErrorForbidden, http.StatusForbidden)
)

// --- Convenience constructors that customize the message key ---

func NotFound(key constant.MsgKey) *AppError {
	return New(ErrNotFound.Code, key, ErrNotFound.HTTPStatus)
}

func BadRequest(key constant.MsgKey) *AppError {
	return New(ErrBadRequest.Code, key, ErrBadRequest.HTTPStatus)
}

func Internal(key constant.MsgKey) *AppError {
	return New(ErrInternal.Code, key, ErrInternal.HTTPStatus)
}

func Conflict(key constant.MsgKey) *AppError {
	return New(ErrConflict.Code, key, ErrConflict.HTTPStatus)
}

func Unauthorized(key constant.MsgKey) *AppError {
	return New(ErrUnauthorized.Code, key, ErrUnauthorized.HTTPStatus)
}
