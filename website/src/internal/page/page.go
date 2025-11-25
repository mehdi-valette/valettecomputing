package page

import (
	"embed"
	"html/template"
	"io"
	"log"
)

//go:embed template/*.html
var fsTemplate embed.FS

var templates *template.Template

func Init() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "template/*.html")

	if err != nil {
		log.Fatal(err)
	}
}

func DisplayIndex(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "index.html", nil)
}

func DisplayContactFormSuccess(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "contactformsuccess.html", nil)
}
