package middleware

import (
	"github.com/davidcm146/bus-booking-be/internal/shared/i18n"
	"github.com/gin-gonic/gin"
)

func I18n() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			lang = c.GetHeader("Accept-Language")
		}
		loc := i18n.NewLocalizer(lang)

		ctx := i18n.WithLocalizer(c.Request.Context(), loc)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
