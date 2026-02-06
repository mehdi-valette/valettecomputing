package page

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io"
	"log"

	"valette.software/internal/blog"
	"valette.software/internal/i18n"
)

//go:embed template
var fsTemplate embed.FS

var templates *template.Template

type templateData struct {
	T i18n.Localizer
}

func Init() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "template/*.*")

	if err != nil {
		log.Fatal(err)
	}
}

func DisplayIndex(buf io.Writer) error {

	type data struct {
		templateData
		Articles []blog.RenderedPost
	}

	localizer, _ := i18n.GetLocale("en")

	return templates.ExecuteTemplate(buf, "index.html", templateData{T: localizer})
}

func DisplayArticlesSummary(buf io.Writer) error {
	articles, err := blog.ListPosts()

	if err != nil {
		return err
	}

	type data struct {
		templateData
		Articles []blog.RenderedPost
	}

	localizer, _ := i18n.GetLocale("fr")

	return templates.ExecuteTemplate(buf, "posts.html", data{
		templateData: templateData{T: localizer}, Articles: articles,
	})
}

func DisplayArticle(buf io.Writer, slug string) error {
	articleText := bytes.NewBuffer(nil)
	article, err := blog.Render(articleText, slug)

	if errors.Is(err, blog.ErrNotFound) {
		return templates.ExecuteTemplate(buf, "post.html", nil)
	}

	type templateData struct {
		Article blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "post.html", templateData{article})
}

func DisplayContactFormSuccess(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "contactformsuccess.html", nil)
}

func DisplayAgenda(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "agenda.html", nil)
}
