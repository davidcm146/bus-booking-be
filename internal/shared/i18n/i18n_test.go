package i18n

import (
	"context"
	"testing"

	"github.com/davidcm146/bus-booking-be/internal/shared/constant"
	"github.com/stretchr/testify/require"
)

func newTestI18n(t *testing.T) I18n {
	t.Helper()
	svc, err := New()
	require.NoError(t, err)
	return svc
}

func TestParseAcceptLanguage(t *testing.T) {
	svc := newTestI18n(t)

	tests := []struct {
		name   string
		header string
		want   constant.Locale
	}{
		{"vi", "vi", constant.LocaleVI},
		{"en", "en", constant.LocaleEN},
		{"vi-VN with q", "vi-VN, en;q=0.9", constant.LocaleVI},
		{"en-US with q", "en-US, vi;q=0.8", constant.LocaleEN},
		{"unknown falls back", "fr-FR", constant.DefaultLocale},
		{"empty falls back", "", constant.DefaultLocale},
		{"case-insensitive", "  VI ", constant.LocaleVI},
		{"first match wins", "zh-CN,ja,en", constant.LocaleEN},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, svc.ParseAcceptLanguage(tt.header))
		})
	}
}

func TestTWithLocale(t *testing.T) {
	svc := newTestI18n(t)

	t.Run("vi and en differ for a known key", func(t *testing.T) {
		viMsg := svc.TWithLocale(constant.LocaleVI, constant.MsgKeyUserNotFound)
		enMsg := svc.TWithLocale(constant.LocaleEN, constant.MsgKeyUserNotFound)
		require.NotEqual(t, viMsg, enMsg)
		require.NotEqual(t, string(constant.MsgKeyUserNotFound), viMsg)
	})

	t.Run("unknown key falls back to the key string", func(t *testing.T) {
		msg := svc.TWithLocale(constant.LocaleEN, constant.MsgKey("nonexistent.key.xyz"))
		require.Equal(t, "nonexistent.key.xyz", msg)
	})
}

func TestT_UsesLocaleFromContext(t *testing.T) {
	svc := newTestI18n(t)

	viCtx := WithLocale(context.Background(), constant.LocaleVI)
	enCtx := WithLocale(context.Background(), constant.LocaleEN)

	viMsg := svc.T(viCtx, constant.MsgKeyUserNotFound)
	enMsg := svc.T(enCtx, constant.MsgKeyUserNotFound)

	require.NotEqual(t, viMsg, enMsg)
}

func TestT_DefaultLocaleWhenUnset(t *testing.T) {
	svc := newTestI18n(t)

	msg := svc.T(context.Background(), constant.MsgKeyUserNotFound)
	enMsg := svc.TWithLocale(constant.LocaleEN, constant.MsgKeyUserNotFound)

	require.Equal(t, enMsg, msg)
}

func TestHasKey(t *testing.T) {
	svc := newTestI18n(t)

	require.True(t, svc.HasKey(constant.MsgKeyUserNotFound))
	require.False(t, svc.HasKey(constant.MsgKey("nonexistent.key.xyz")))
}

func TestSetup_PanicsOnDoubleInit(t *testing.T) {
	require.Panics(t, func() {
		svc, _ := New()
		Setup(svc)
		Setup(svc)
	})
}
