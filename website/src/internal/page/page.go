package page

import (
	"bytes"
	"embed"
	"errors"
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

func DisplayArticlesSummary(buf io.Writer) error {
	articles, err := blog.ListArticles()

	if err != nil {
		return err
	}

	type templateData struct {
		Articles []blog.Article
	}

	return templates.ExecuteTemplate(buf, "articles.html", templateData{articles})
}

func DisplayArticle(buf io.Writer, slug string) error {
	articleText := bytes.NewBuffer(nil)
	article, err := blog.Render(articleText, slug)

	if errors.Is(err, blog.ErrNotFound) {
		return templates.ExecuteTemplate(buf, "article.html", nil)
	}

	type templateData struct {
		Article blog.RenderedArticle
	}

	return templates.ExecuteTemplate(buf, "article.html", templateData{article})
}

func DisplayContactFormSuccess(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "contactformsuccess.html", nil)
}

func DisplayAgenda(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "agenda.html", nil)
}
