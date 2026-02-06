package i18n

import (
	"embed"
	"errors"
	"io/fs"

	"github.com/leonelquinteros/gotext"
)

//go:embed "locales"
var embededFs embed.FS

var localesFs fs.FS

var locales map[string]*gotext.Locale

var ErrLocaleNotFound = errors.New("locale not found")

func Init() {
	var err error

	localesFs, err = fs.Sub(embededFs, "locales")

	if err != nil {
		panic("cannot load the locales")
	}

	locales = map[string]*gotext.Locale{
		"fr": gotext.NewLocaleFS("fr", localesFs),
		"en": gotext.NewLocaleFS("en", localesFs),
	}

	for _, locale := range locales {
		locale.AddDomain("main")
	}
}

type Localizer interface {
	Get(key string, args ...any) string
}

func GetLocale(lang string) (Localizer, error) {
	locale, ok := locales[lang]

	if !ok {
		return &gotext.Locale{}, ErrLocaleNotFound
	}

	return locale, nil
}
