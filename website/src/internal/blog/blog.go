package blog

import (
	"embed"
	"errors"
	"html/template"
	"io"
	"io/fs"
	"log"
	"regexp"
	"strings"
	"time"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed articles
var fsBlog embed.FS

var resources fs.FS
var htmlRenderer *html.Renderer
var mdExtensions parser.Extensions
var db *sql.DB

func Init() {
	var err error

	db, err = sql.Open("sqlite", "../articles.db")

	if err != nil {
		log.Fatal("couldn't open the articles' database:", err)
	}

	resources, err = fs.Sub(fsBlog, "articles")
	mdExtensions = parser.CommonExtensions

	if err != nil {
		log.Fatal("the blog articles cannot be found")
	}

	htmlRenderer = html.NewRenderer(html.RendererOptions{})
}

func AddArticle(newArticle NewArticle) error {
	currentTimestamp := time.Now().UnixMilli() / 1000
	slug := makeSlug(newArticle.Title)

	_, err := db.Exec(
		"INSERT INTO article(title, timestamp, slug, summary, content) VALUES(?, ?, ?, ?, ?)",
		newArticle.Title, currentTimestamp, slug, newArticle.Summary, newArticle.Content,
	)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func ListArticles() ([]Article, error) {
	currentArticle := Article{}
	allArticles := []Article{}

	results, err := db.Query("SELECT title, timestamp, summary, slug FROM article")

	if err != nil {
		return []Article{}, err
	}

	for results.Next() {
		err := results.Scan(
			&currentArticle.Title,
			&currentArticle.Timestamp,
			&currentArticle.Summary,
			&currentArticle.Slug,
		)

		if err != nil {
			return []Article{}, err
		}

		allArticles = append(allArticles, currentArticle)
	}

	return allArticles, nil
}

func Render(output io.Writer, slug string) (RenderedArticle, error) {
	article := RenderedArticle{}

	result := db.QueryRow("SELECT title, timestamp, summary, content FROM article WHERE slug = ?", slug)

	err := result.Scan(&article.Title, &article.Timestamp, &article.Summary, &article.Content)

	if errors.Is(err, sql.ErrNoRows) {
		return RenderedArticle{}, ErrNotFound
	} else if err != nil {
		log.Print(err)
		return RenderedArticle{}, err
	}

	htmlFile := markdown.Render(
		parser.NewWithExtensions(parser.CommonExtensions).Parse([]byte(article.Content)),
		htmlRenderer,
	)

	article.Html = template.HTML(htmlFile)

	return article, err
}

func makeSlug(text string) string {
	slug := strings.ToLower(text)
	slug = strings.ReplaceAll(slug, " ", "-")
	reg := regexp.MustCompile(`[^a-z0-9-]`)

	slug = string(reg.ReplaceAll([]byte(slug), []byte("0")))

	return slug
}
