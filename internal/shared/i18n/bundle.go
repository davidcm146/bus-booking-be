package i18n

import (
	"embed"
	"sync"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

//go:embed locales/*.yaml
var localeFS embed.FS

var (
	bundle *i18n.Bundle
	once   sync.Once
)

func Init() {
	once.Do(func() {
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

		bundle.LoadMessageFileFS(localeFS, "locales/en.yaml")
		bundle.LoadMessageFileFS(localeFS, "locales/vi.yaml")
	})
}

func GetBundle() *i18n.Bundle {
	return bundle
}
