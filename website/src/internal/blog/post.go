package blog

import (
	"fmt"
	"html/template"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

type NewPost struct {
	Author    string `json:"author"`
	Language  string `json:"language"`
	Timestamp int64  `json:"timestamp"`
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
}

type Post struct {
	ArticleId int64
	Language  string
	Slug      string
	Author    string
	Title     string
	Timestamp int64
	Summary   string
	Content   string
}

type RenderedPost struct {
	Post
	Html      template.HTML
	DateHuman string
	DateIso   string
}

func (post *NewPost) ToRenderedPost(id int64, slug string) RenderedPost {
	renderedPost := RenderedPost{
		Post: Post{
			ArticleId: id,
			Language:  post.Language,
			Author:    post.Author,
			Title:     post.Title,
			Timestamp: post.Timestamp,
			Summary:   post.Summary,
			Content:   post.Content,
			Slug:      slug,
		},
	}

	renderedPost.CalculateDates()
	renderedPost.CalculateHtmlContent()

	return renderedPost
}

func (post *RenderedPost) CalculateDates() {
	datetime := time.Unix(post.Timestamp, 0).In(timezoneCet)

	var month string

	if post.Language == "en" {
		month = monthsEn[datetime.Month()-1]
		post.DateHuman = fmt.Sprintf("%s %d, %d", month, datetime.Day(), datetime.Year())
	} else {
		month = monthsFr[datetime.Month()-1]
		post.DateHuman = fmt.Sprintf("%d %s %d", datetime.Day(), month, datetime.Year())
	}

	post.DateIso = datetime.UTC().Format("2006-01-02T15:04:05Z")
}

func (post *RenderedPost) CalculateHtmlContent() {
	htmlFile := markdown.Render(
		parser.NewWithExtensions(parser.CommonExtensions).Parse([]byte(post.Content)),
		htmlRenderer,
	)

	post.Html = template.HTML(htmlFile)
}
