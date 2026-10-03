package validator

import (
	"context"

	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
)

// GetMessage trả về message lỗi validate đã dịch, với Field và Param được điền qua template.
func GetMessage(ctx context.Context, tag, field, param string) string {
	key := "validator." + tag
	if messageKey, ok := customMessages[tag]; ok {
		key = string(messageKey)
	}

	data := map[string]interface{}{
		"Field": field,
		"Param": param,
	}

	msg := i18n.T(ctx, key, data)
	if msg == key {
		return i18n.T(ctx, "validator.fallback", map[string]interface{}{
			"Field": field,
			"Tag":   tag,
		})
	}

	return msg
}
