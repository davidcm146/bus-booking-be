package validator

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
)

// Return localized error message for the given validator tag.
func GetMessage(ctx context.Context, tag, field, param string) string {
	// Check custom validator messages first
	if key, found := customMessages[tag]; found {
		return i18n.T(ctx, key, field)
	}

	// Standard validator tag → construct key dynamically (e.g. "validator.required")
	key := constant.MsgKey("validator." + tag)

	// Check if this key has a translation; if not, use the fallback
	if !i18n.HasKey(key) {
		return i18n.T(ctx, constant.MsgKeyFallback, field, tag)
	}

	if param != "" {
		return i18n.T(ctx, key, field, param)
	}
	return i18n.T(ctx, key, field)
}
