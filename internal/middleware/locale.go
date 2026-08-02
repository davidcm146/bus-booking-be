package middleware

import (
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/gin-gonic/gin"
)

// LocaleMiddleware parses the Accept-Language header and stores
// the resolved locale in the request context. The I18n service
// itself is globally available via i18n.MustInit — not stored in context.
func LocaleMiddleware(i18nService i18n.I18n) gin.HandlerFunc {
	return func(c *gin.Context) {
		locale := i18nService.ParseAcceptLanguage(c.GetHeader("Accept-Language"))
		ctx := i18n.WithLocale(c.Request.Context(), locale)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
