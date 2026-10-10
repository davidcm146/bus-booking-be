// internal/shared/i18n/localizer.go
package i18n

import (
	"context"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
)

type ctxKey string

const localizerCtxKey ctxKey = "i18n_localizer"

func NewLocalizer(lang string) *goi18n.Localizer {
	return goi18n.NewLocalizer(GetBundle(), lang, "en")
}

func WithLocalizer(ctx context.Context, loc *goi18n.Localizer) context.Context {
	return context.WithValue(ctx, localizerCtxKey, loc)
}

func FromContext(ctx context.Context) *goi18n.Localizer {
	loc, ok := ctx.Value(localizerCtxKey).(*goi18n.Localizer)
	if !ok || loc == nil {
		return NewLocalizer("en") // Fallback
	}
	return loc
}

func T(ctx context.Context, key string, data map[string]interface{}) string {
	loc := FromContext(ctx)
	msg, err := loc.Localize(&goi18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: data,
	})
	if err != nil {
		return key
	}
	return msg
}
