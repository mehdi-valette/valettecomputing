package blog

import "html/template"

type NewArticle struct {
	Title     string `json:"title"`
	Timestamp int64  `json:"timestamp"`
	Summary   string `json:"summary"`
	Content   string `json:"content"`
}

type Article struct {
	ArticleId int64
	Slug      string
	Title     string
	Timestamp int64
	Summary   string
	Content   string
}

type RenderedArticle struct {
	Article
	Html template.HTML
}
