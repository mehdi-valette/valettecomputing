package page

import (
	"bytes"
	"embed"
	"html/template"
	"io"
	"log"

	"valette.software/internal/blog"
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

func DisplayArticle(buf io.Writer, articleName string) error {
	articleText := bytes.NewBuffer(nil)
	err := blog.Render(articleText, articleName)

	if err != nil {
		return err
	}

	return templates.ExecuteTemplate(buf, "article.html", struct{ Article template.HTML }{Article: template.HTML(articleText.String())})
}

func DisplayContactFormSuccess(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "contactformsuccess.html", nil)
}

func DisplayAgenda(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "agenda.html", nil)
}
