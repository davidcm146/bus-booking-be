package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	i18nv2 "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

//go:embed locales/*.json
var localeFS embed.FS

// I18n defines the interface for language translation and locale management.
type I18n interface {
	T(ctx context.Context, key constant.MsgKey, args ...any) string
	TWithLocale(locale constant.Locale, key constant.MsgKey, args ...any) string
	ParseAcceptLanguage(header string) constant.Locale
	HasKey(key constant.MsgKey) bool
}

type service struct {
	bundle     *i18nv2.Bundle
	localizers map[constant.Locale]*i18nv2.Localizer
}

// --- Global instance (set once at startup via MustInit) ---

var (
	global     I18n
	globalOnce sync.Once
)

// MustInit sets the global I18n instance. Must be called exactly once during
// application bootstrap before any handler serves traffic.
// Panics if called more than once to catch wiring bugs early.
func Setup(instance I18n) {
	initialized := false
	globalOnce.Do(func() {
		global = instance
		initialized = true
	})
	if !initialized {
		panic("i18n: MustInit called more than once")
	}
}

// Global returns the global I18n instance. Panics if MustInit was not called.
func Global() I18n {
	if global == nil {
		panic("i18n: MustInit has not been called")
	}
	return global
}

// New creates and initializes a new I18n instance.
// Pure constructor with no side-effects — call MustInit separately to register globally.
func New() (I18n, error) {
	bundle := i18nv2.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	_, err := bundle.LoadMessageFileFS(localeFS, "locales/en.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load en.json: %w", err)
	}

	_, err = bundle.LoadMessageFileFS(localeFS, "locales/vi.json")
	if err != nil {
		return nil, fmt.Errorf("failed to load vi.json: %w", err)
	}

	s := &service{
		bundle: bundle,
		localizers: map[constant.Locale]*i18nv2.Localizer{
			constant.LocaleEN: i18nv2.NewLocalizer(bundle, string(constant.LocaleEN)),
			constant.LocaleVI: i18nv2.NewLocalizer(bundle, string(constant.LocaleVI), string(constant.LocaleEN)),
		},
	}

	return s, nil
}

// --- Context helpers: locale value only ---

type localeCtxKey struct{}

// WithLocale stores a Locale in the context.
func WithLocale(ctx context.Context, locale constant.Locale) context.Context {
	return context.WithValue(ctx, localeCtxKey{}, locale)
}

// LocaleFromContext extracts the Locale from context. Returns DefaultLocale if unset.
func LocaleFromContext(ctx context.Context) constant.Locale {
	if locale, ok := ctx.Value(localeCtxKey{}).(constant.Locale); ok {
		return locale
	}
	return constant.DefaultLocale
}

// --- Package-level convenience functions (use the global instance) ---

// T translates a key using the global I18n instance and the locale from context.
func T(ctx context.Context, key constant.MsgKey, args ...any) string {
	return Global().T(ctx, key, args...)
}

// HasKey checks if a key exists in the global I18n service.
func HasKey(key constant.MsgKey) bool {
	return Global().HasKey(key)
}

// --- Service methods ---

// ParseAcceptLanguage extracts the best-matching Locale from an Accept-Language header value.
func (s *service) ParseAcceptLanguage(header string) constant.Locale {
	header = strings.ToLower(strings.TrimSpace(header))

	for _, part := range strings.Split(header, ",") {
		lang := strings.TrimSpace(strings.SplitN(part, ";", 2)[0])
		if strings.HasPrefix(lang, "vi") {
			return constant.LocaleVI
		}
		if strings.HasPrefix(lang, "en") {
			return constant.LocaleEN
		}
	}

	return constant.DefaultLocale
}

// T translates a message key using the locale from context.
func (s *service) T(ctx context.Context, key constant.MsgKey, args ...any) string {
	locale := LocaleFromContext(ctx)
	return s.TWithLocale(locale, key, args...)
}

// TWithLocale translates a message key for a specific locale.
func (s *service) TWithLocale(locale constant.Locale, key constant.MsgKey, args ...any) string {
	localizer, ok := s.localizers[locale]
	if !ok {
		localizer = s.localizers[constant.DefaultLocale]
	}

	msg, err := localizer.Localize(&i18nv2.LocalizeConfig{
		MessageID: string(key),
	})
	if err != nil || msg == "" {
		msg = string(key)
	}

	if len(args) > 0 {
		return fmt.Sprintf(msg, args...)
	}
	return msg
}

// HasKey checks if a translation exists for the given key in the bundle.
func (s *service) HasKey(key constant.MsgKey) bool {
	enLocalizer := s.localizers[constant.LocaleEN]
	_, err := enLocalizer.Localize(&i18nv2.LocalizeConfig{
		MessageID: string(key),
	})
	return err == nil
}
