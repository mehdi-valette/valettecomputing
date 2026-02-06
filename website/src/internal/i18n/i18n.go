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

var locales map[string]*Locale

var ErrLocaleNotFound = errors.New("locale not found")

type Localizer interface {
	Get(key string, args ...any) string
	Lang() string
	Link(path string) string
}

type Locale struct {
	locale *gotext.Locale
}

var _ Localizer = &Locale{}

func Init() {
	var err error

	localesFs, err = fs.Sub(embededFs, "locales")

	if err != nil {
		panic("cannot load the locales")
	}

	locales = map[string]*Locale{
		"fr": {locale: gotext.NewLocaleFS("fr", localesFs)},
		"en": {locale: gotext.NewLocaleFS("en", localesFs)},
	}

	for _, locale := range locales {
		locale.locale.AddDomain("main")
	}
}

func GetLocale(lang string) (Localizer, error) {
	locale, ok := locales[lang]

	if !ok {
		return &Locale{}, ErrLocaleNotFound
	}

	return locale, nil
}

func (l *Locale) Get(str string, args ...any) string {
	return l.locale.Get(str, args...)
}

func (l *Locale) Link(path string) string {
	if path[0] != '/' {
		path = "/" + path
	}

	return "/" + l.locale.GetLanguage() + path
}

func (l *Locale) Lang() string {
	return l.locale.GetLanguage()
}
