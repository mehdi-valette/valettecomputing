package router

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"valette.software/internal/blog"
	"valette.software/internal/contactform"
	"valette.software/internal/i18n"
	"valette.software/internal/page"
	"valette.software/internal/reqcontext"
	"valette.software/internal/static"
)

func Build() *http.ServeMux {

	root := http.NewServeMux()
	router := buildRouter()

	root.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		ctxValue := reqcontext.ReqContext{}

		// add the locale in the context and remove it from the path
		req.URL.Path = AddLocale(&ctxValue, req.URL.Path)

		if req.URL.Path == "" {
			req.URL.Path = "/"
		}

		// add the current path in the context
		ctxValue.CurrentPath = req.URL.Path

		newCtx := reqcontext.SetValue(req.Context(), ctxValue)

		router.ServeHTTP(res, req.WithContext(newCtx))
	})

	return root
}

func AddLocale(ctx *reqcontext.ReqContext, path string) string {
	stripPrefix := ""
	pathFragments := strings.Split(path, "/")

	// [lang] is at the beginning of the path
	// the path begins at index 0 or 1 depending on the leading slash
	var firstPathParam string

	if path[0] == '/' && len(pathFragments) > 1 {
		firstPathParam = pathFragments[1]
	} else if len(pathFragments) > 0 {
		firstPathParam = pathFragments[0]
	} else {
		firstPathParam = ""
	}

	var localizer i18n.Localizer
	var err error

	switch firstPathParam {
	case "en":
		localizer, err = i18n.GetLocale("en")
		stripPrefix = "/en"
	case "fr":
		localizer, err = i18n.GetLocale("fr")
		stripPrefix = "/fr"
	default:
		localizer, err = i18n.GetLocale("fr")
		stripPrefix = ""
	}

	if err != nil {
		log.Print("locale not found\n", err)
	}

	ctx.Localizer = localizer

	return strings.TrimPrefix(path, stripPrefix)
}

func buildRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /static/", http.StripPrefix("/static/", static.Serve()))

	router.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		reqCtx := reqcontext.GetValue(req.Context())

		err := page.DisplayIndex(res, reqCtx)

		if err != nil {
			log.Print("couldn't display the index page\n", err)
		}
	})

	router.HandleFunc("GET /admin/", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayAdmin(res)

		if err != nil {
			log.Print("Couldn't display the admin page\n", err)
		}
	})

	router.HandleFunc("GET /articles/", func(res http.ResponseWriter, req *http.Request) {
		reqCtx := reqcontext.GetValue(req.Context())

		err := page.DisplayArticlesSummary(res, reqCtx)

		if err != nil {
			log.Print("couldn't display the articles's summary page\n", err)
		}
	})

	router.HandleFunc("GET /new-post", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayPostNew(res)

		if err != nil {
			res.WriteHeader(500)
			log.Print(err)
		}
	})

	router.HandleFunc("GET /edit-posts/{id}", func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(
			req.PathValue("id"),
			10,
			64,
		)

		if err != nil {
			res.WriteHeader(400)
			res.Write([]byte("post's id must be a number"))
			return
		}

		post, err := blog.GetPostById(id)

		if err != nil {
			res.WriteHeader(500)
			return
		}

		page.DisplayPostEdition(res, post)
	})

	router.HandleFunc("POST /posts", func(res http.ResponseWriter, req *http.Request) {
		date, err := time.Parse("2006-01-02", req.FormValue("date"))

		if err != nil {
			date = time.Now()
		}

		newPost := blog.NewPost{
			Author:    req.FormValue("author"),
			Language:  req.FormValue("language"),
			Title:     req.FormValue("title"),
			Timestamp: date.Unix(),
			Summary:   req.FormValue("summary"),
			Content:   req.FormValue("content"),
		}

		renderedPost, err := blog.AddPost(newPost)

		if err != nil {
			res.WriteHeader(500)
			log.Print(err)
		}

		err = page.DisplayPostEdition(res, renderedPost)

		if err != nil {
			log.Print("couldn't display the post after creation\n", err)
		}
	})

	router.HandleFunc("PUT /posts/{id}", func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(req.FormValue("id"), 10, 64)

		if err != nil {
			res.WriteHeader(400)
			res.Write([]byte("the post's ID must be an integer"))
			return
		}

		date, err := time.Parse("2006-01-02", req.FormValue("date"))

		if err != nil {
			date = time.Now()
		}

		newPost := blog.RenderedPost{
			Post: blog.Post{
				ArticleId: id,
				Slug:      req.FormValue("slug"),
				Author:    req.FormValue("author"),
				Language:  req.FormValue("language"),
				Title:     req.FormValue("title"),
				Timestamp: date.Unix(),
				Summary:   req.FormValue("summary"),
				Content:   req.FormValue("content"),
			},
		}

		renderedPost, err := blog.UpdatePost(newPost)

		if err != nil {
			res.Write([]byte(err.Error()))
			log.Print(err)
			return
		}

		err = page.DisplayPostEdition(res, renderedPost)

		if err != nil {
			log.Print("couldn't display the post after update\n", err)
		}
	})

	router.HandleFunc("DELETE /posts/{id}", func(res http.ResponseWriter, req *http.Request) {
		id, err := strconv.ParseInt(req.FormValue("id"), 10, 64)

		if err != nil {
			res.WriteHeader(500)
			return
		}

		err = blog.DeletePostById(id)

		if err != nil {
			res.WriteHeader(500)
			log.Print(err)
		}

		err = page.DisplayPostNew(res)

		if err != nil {
			log.Print(err)
		}
	})

	router.HandleFunc("GET /articles/{name}", func(res http.ResponseWriter, req *http.Request) {
		reqCtx := reqcontext.GetValue(req.Context())

		err := page.DisplayArticle(res, reqCtx, req.PathValue("name"))

		if err != nil {
			log.Print("couldn't display the article page\n", err)
		}
	})

	router.HandleFunc("GET /agenda", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayAgenda(res)

		if err != nil {
			log.Print("couldn't display the agenda page", err)
		}
	})

	router.HandleFunc("POST /contact", contactform.HandleContactFormRequest)

	return router
}
