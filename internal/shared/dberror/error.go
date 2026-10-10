package dberror

import (
	"context"
	"errors"

	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/jackc/pgx/v5/pgconn"
)

type config struct {
	conflicts map[string]constant.MsgKey // constraint name -> message key
}

type Option func(*config)

func WithConflict(constraint string, key constant.MsgKey) Option {
	return func(c *config) { c.conflicts[constraint] = key }
}

func Translate(err error, fallback constant.MsgKey, opts ...Option) error {
	if err == nil {
		return nil
	}

	cfg := &config{conflicts: map[string]constant.MsgKey{}}
	for _, o := range opts {
		o(cfg)
	}

	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return apperror.Internal(fallback).WithCause(err)
	}

	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return apperror.Internal(fallback).WithCause(err)
	}

	switch pgErr.Code {
	case constant.PgCodeUniqueViolation:
		key, ok := cfg.conflicts[pgErr.ConstraintName]
		if !ok {
			key = i18n.KeyErrorResourceConflict
		}
		return apperror.Conflict(key).WithCause(err)

	case constant.PgCodeForeignKeyViolation:
		return apperror.BadRequest(i18n.KeyErrorInvalidReference).WithCause(err)

	case constant.PgCodeNotNullViolation, constant.PgCodeCheckViolation, constant.PgCodeInvalidText:
		return apperror.BadRequest(i18n.KeyErrorInvalidData).WithCause(err)

	default:
		return apperror.Internal(fallback).WithCause(err)
	}
}
