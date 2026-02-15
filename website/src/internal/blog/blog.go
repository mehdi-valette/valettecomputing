package blog

import (
	"errors"
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
var monthsEn = []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}

func Init() {
	var err error

	db, err = sql.Open("sqlite", "/var/lib/valettesoftware/blog.db")

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

func AddPost(newPost NewPost) (RenderedPost, error) {
	slug := makeSlug(newPost.Title)

	if newPost.Timestamp == 0 {
		newPost.Timestamp = time.Now().Unix()
	}

	result, err := db.Exec(
		"INSERT INTO post(title, language, author, timestamp, slug, summary, content) VALUES(?, ?, ?, ?, ?, ?, ?)",
		newPost.Title, newPost.Language, newPost.Author, newPost.Timestamp, slug, newPost.Summary, newPost.Content,
	)

	if err != nil {
		return RenderedPost{}, err
	}

	newId, err := result.LastInsertId()

	if err != nil {
		return RenderedPost{}, err
	}

	renderedPost := newPost.ToRenderedPost(newId, slug)

	return renderedPost, nil
}

func UpdatePost(post RenderedPost) (RenderedPost, error) {
	_, err := db.Exec(
		"UPDATE post SET title = ?, language = ?, author = ?, timestamp = ?, slug = ?, summary = ?, content = ? WHERE post_id = ?",
		post.Title, post.Language, post.Author, post.Timestamp, post.Slug, post.Summary, post.Content, post.ArticleId,
	)

	if err != nil {
		return RenderedPost{}, err
	}

	post.CalculateDates()
	post.CalculateHtmlContent()

	return post, nil
}

func ListPosts(lang string) ([]RenderedPost, error) {
	currentPost := RenderedPost{}
	allPosts := []RenderedPost{}
	var err error
	var results *sql.Rows

	if lang != "" {
		results, err = db.Query("SELECT post_id, title, author, language, timestamp, summary, slug FROM post WHERE language = ? ORDER BY timestamp DESC", lang)
	} else {
		results, err = db.Query("SELECT post_id, title, author, language, timestamp, summary, slug FROM post ORDER BY timestamp DESC")
	}

	if err != nil {
		return []RenderedPost{}, err
	}

	for results.Next() {
		err := results.Scan(
			&currentPost.ArticleId,
			&currentPost.Title,
			&currentPost.Author,
			&currentPost.Language,
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

func GetPostBySlug(slug string) (RenderedPost, error) {
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

func GetPostById(id int64) (RenderedPost, error) {
	post := RenderedPost{}

	result := db.QueryRow("SELECT post_id, language, slug, title, author, timestamp, summary, content FROM post WHERE post_id = ?", id)

	err := result.Scan(&post.ArticleId, &post.Language, &post.Slug, &post.Title, &post.Author, &post.Timestamp, &post.Summary, &post.Content)

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

func DeletePostById(id int64) error {
	_, err := db.Exec("DELETE FROM post WHERE post_id = ?", id)

	return err
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
