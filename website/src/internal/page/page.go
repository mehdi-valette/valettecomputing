package page

import (
	"bytes"
	"embed"
	"errors"
	"html/template"
	"io"
	"log"

	"valette.software/internal/blog"
	"valette.software/internal/reqcontext"
)

//go:embed template
var fsTemplate embed.FS

var templates *template.Template

type templateData struct {
	Ctx reqcontext.ReqContext
}

func Init() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "template/*.*")

	if err != nil {
		log.Fatal(err)
	}
}

func DisplayIndex(buf io.Writer, reqCtx reqcontext.ReqContext) error {

	type data struct {
		templateData
		Articles []blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "index.html", templateData{Ctx: reqCtx})
}

func DisplayArticlesSummary(buf io.Writer, reqCtx reqcontext.ReqContext) error {
	articles, err := blog.ListPosts(reqCtx.Localizer.Lang())

	if err != nil {
		return err
	}

	type data struct {
		templateData
		Articles []blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "posts.html", data{
		templateData: templateData{Ctx: reqCtx}, Articles: articles,
	})
}

func DisplayArticle(buf io.Writer, reqCtx reqcontext.ReqContext, slug string) error {
	articleText := bytes.NewBuffer(nil)
	article, err := blog.GetPostBySlug(articleText, slug)

	if errors.Is(err, blog.ErrNotFound) {
		return templates.ExecuteTemplate(buf, "post.html", nil)
	}

	type data struct {
		templateData
		Article blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "post.html", data{templateData: templateData{Ctx: reqCtx}, Article: article})
}

func DisplayContactFormSuccess(buf io.Writer, reqCtx reqcontext.ReqContext) error {
	return templates.ExecuteTemplate(buf, "contactformsuccess.html", templateData{Ctx: reqCtx})
}

func DisplayAgenda(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "agenda.html", nil)
}

func DisplayAdmin(buf io.Writer) error {
	articles, err := blog.ListPosts("")

	if err != nil {
		return err
	}

	type listItem struct {
		Status string
		Post   blog.RenderedPost
	}

	ps := make([]listItem, 0, len(articles))

	for _, p := range articles {
		ps = append(ps, listItem{Status: "update", Post: p})
	}

	type data struct {
		Posts []listItem
	}

	return templates.ExecuteTemplate(buf, "admin.html", data{ps})
}

func DisplayPostEdition(buf io.Writer, post blog.RenderedPost) error {
	type data struct {
		Post blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "post-edit.html", data{Post: post})
}

func DisplayPostNew(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "post-edit.html", nil)
}

func DisplayLoginForm(buf io.Writer) error {
	return templates.ExecuteTemplate(buf, "admin-login.html", nil)
}

func DisplayPostListItem(buf io.Writer, post blog.RenderedPost, status string) error {
	type data struct {
		Status string
		Post   blog.RenderedPost
	}

	return templates.ExecuteTemplate(buf, "post-edit-list-item.html", data{Status: status, Post: post})
}
