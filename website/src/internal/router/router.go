package router

import (
	"log"
	"net/http"
	"strings"

	"valette.software/internal/authentication"
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
		localizer, newPath := getLocale(req.URL.Path)
		sessionId := getSessionId(req)

		req.URL.Path = newPath

		ctxValue := reqcontext.ReqContext{
			Localizer:   localizer,
			Admin:       authentication.CheckSession(sessionId),
			CurrentPath: newPath,
		}

		newCtx := reqcontext.SetValue(req.Context(), ctxValue)

		router.ServeHTTP(res, req.WithContext(newCtx))
	})

	return root
}

func getLocale(path string) (i18n.Localizer, string) {
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

	newPath := strings.TrimPrefix(path, stripPrefix)

	if newPath == "" {
		newPath = "/"
	}

	return localizer, newPath
}

func getSessionId(req *http.Request) string {
	cookie, err := req.Cookie("session-id")

	if err != nil {
		return ""
	}

	return cookie.Value
}

func buildRouter() *http.ServeMux {
	router := http.NewServeMux()

	router.Handle("GET /static/", http.StripPrefix("/static/", static.Serve()))

	router.HandleFunc("GET /", indexPage)

	router.HandleFunc("GET /articles/", listPosts)

	router.HandleFunc("POST /login", login)

	router.HandleFunc("GET /logout", logout)

	router.HandleFunc("GET /articles/{name}", getPost)

	router.HandleFunc("GET /agenda", getAgenda)

	router.HandleFunc("POST /contact", contactform.HandleContactFormRequest)

	router.HandleFunc("GET /admin/", requireAdmin(adminPage))

	router.HandleFunc("GET /new-post", requireAdmin(newPostController))

	router.HandleFunc("GET /edit-posts/{id}", requireAdmin(getEditablePost))

	router.HandleFunc("POST /posts", requireAdmin(createPost))

	router.HandleFunc("PUT /posts/{id}", requireAdmin(updatePost))

	router.HandleFunc("DELETE /posts/{id}", requireAdmin(deletePost))

	return router
}

func requireAdmin(handler http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		reqCtx := reqcontext.GetValue(req.Context())

		if reqCtx.Admin {
			handler(res, req)
			return
		}

		printError(page.DisplayLoginForm(res))
	}
}

func printError(err error) {
	if err != nil {
		log.Print(err)
	}
}
