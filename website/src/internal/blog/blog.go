package blog

import (
	"errors"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

var htmlRenderer *html.Renderer
var mdExtensions parser.Extensions
var db *sql.DB
var timezoneCet *time.Location
var monthsFr = []string{"janvier", "février", "mars", "avril", "mai", "juin", "juillet", "août", "septembre", "octobre", "novembre", "décembre"}

func Init() {
	var err error

	db, err = sql.Open("sqlite", "../blog.db")

	if err != nil {
		log.Fatal("couldn't open the blog's database:", err)
	}

	mdExtensions = parser.CommonExtensions

	htmlRenderer = html.NewRenderer(html.RendererOptions{})

	timezoneCet, err = time.LoadLocation("Europe/Zurich")

	if err != nil {
		log.Fatal("couldn't load the timezone Europe/Zurich")
	}
}

func AddPost(newArticle NewPost) error {
	currentTimestamp := time.Now().UnixMilli() / 1000
	slug := makeSlug(newArticle.Title)

	_, err := db.Exec(
		"INSERT INTO post(title, author, timestamp, slug, summary, content) VALUES(?, ?, ?, ?, ?, ?)",
		newArticle.Title, newArticle.Author, currentTimestamp, slug, newArticle.Summary, newArticle.Content,
	)

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}

func ListPosts() ([]RenderedPost, error) {
	currentPost := RenderedPost{}
	allPosts := []RenderedPost{}

	results, err := db.Query("SELECT title, author, timestamp, summary, slug FROM post")

	if err != nil {
		return []RenderedPost{}, err
	}

	for results.Next() {
		err := results.Scan(
			&currentPost.Title,
			&currentPost.Author,
			&currentPost.Timestamp,
			&currentPost.Summary,
			&currentPost.Slug,
		)

		if err != nil {
			return []RenderedPost{}, err
		}

		currentPost.CalculateDates()

		allPosts = append(allPosts, currentPost)
	}

	return allPosts, nil
}

func Render(output io.Writer, slug string) (RenderedPost, error) {
	post := RenderedPost{}

	result := db.QueryRow("SELECT title, author, timestamp, summary, content FROM post WHERE slug = ?", slug)

	err := result.Scan(&post.Title, &post.Author, &post.Timestamp, &post.Summary, &post.Content)

	if errors.Is(err, sql.ErrNoRows) {
		return RenderedPost{}, ErrNotFound
	} else if err != nil {
		log.Print(err)
		return RenderedPost{}, err
	}

	post.CalculateDates()
	post.CalculateHtmlContent()

	return post, err
}

func makeSlug(text string) string {
	slug := strings.ToLower(text)

	charsToRemove := regexp.MustCompile(`[.,:;!?^'"]`)
	charsToMap := strings.NewReplacer(" ", "-", "ç", "c", "à", "a", "â", "a", "é", "e", "è", "e", "ê", "e", "ë", "e", "î", "i", "ï", "i", "ô", "o", "ö", "o", "ù", "u", "û", "u", "ü", "u", "ÿ", "y")
	unknownChars := regexp.MustCompile(`[^a-z0-9-]`)

	slug = string(charsToRemove.ReplaceAll([]byte(slug), []byte("")))
	slug = charsToMap.Replace(slug)
	slug = string(unknownChars.ReplaceAll([]byte(slug), []byte("0")))

	return slug
}
