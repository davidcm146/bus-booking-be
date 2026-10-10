package utils

import (
	"strings"
	"time"
	"unicode"

	"github.com/davidcm146/bus-booking-be/internal/shared/apperror"
	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
)

var businessLocation = time.FixedZone(constant.BusinessTimeZoneName, constant.BusinessTimeZoneOffset)

func DeferString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ParseLocalDay(value string) (time.Time, error) {
	day, err := time.ParseInLocation(constant.DateLayout, value, businessLocation)
	if err != nil {
		return time.Time{}, apperror.BadRequest(i18n.KeyErrorBadRequest)
	}
	return day, nil
}

func Later(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func ToCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

func Paginate(page, limit, defaultSize, maxSize int) (int, int) {
	if limit <= 0 {
		limit = defaultSize
	}
	limit = min(limit, maxSize)
	return limit, (max(page, 1) - 1) * limit
}

func ExtractBearerToken(header string) (string, bool) {
	scheme, token, ok := strings.Cut(header, " ")
	if !ok || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}
	return token, true
}
