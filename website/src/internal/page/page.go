package page

import (
	"embed"
	"html/template"
	"io"
	"log"
)

//go:embed template
var fsTemplate embed.FS

var templates *template.Template

func Init() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "template/*.*")

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

func DisplayAgenda(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "agenda.html", nil)
}
