package router

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"valette.software/internal/authentication"
	"valette.software/internal/blog"
	"valette.software/internal/page"
	"valette.software/internal/reqcontext"
)

func getAgenda(res http.ResponseWriter, req *http.Request) {
	printError(page.DisplayAgenda(res))
}

func indexPage(res http.ResponseWriter, req *http.Request) {
	reqCtx := reqcontext.GetValue(req.Context())
	printError(page.DisplayIndex(res, reqCtx))
}

func getPost(res http.ResponseWriter, req *http.Request) {
	reqCtx := reqcontext.GetValue(req.Context())

	printError(page.DisplayArticle(res, reqCtx, req.PathValue("name")))
}

func adminPage(res http.ResponseWriter, req *http.Request) {
	printError(page.DisplayAdmin(res))
}

func listPosts(res http.ResponseWriter, req *http.Request) {
	reqCtx := reqcontext.GetValue(req.Context())

	printError(page.DisplayArticlesSummary(res, reqCtx))
}

func login(res http.ResponseWriter, req *http.Request) {
	sessionId, err := authentication.Authenticate(req.FormValue("password"))

	if err != nil {
		printError(page.DisplayLoginForm(res))
		return
	}

	sessionCookie := http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Name:     "session-id",
		Value:    sessionId,
	}

	http.SetCookie(res, &sessionCookie)

	req.Method = "GET"
	http.Redirect(res, req, "/admin", 303)

	printError(page.DisplayAdmin(res))
}

func logout(res http.ResponseWriter, req *http.Request) {
	authentication.Logout()
	http.Redirect(res, req, "/", http.StatusTemporaryRedirect)
}

func newPostController(res http.ResponseWriter, req *http.Request) {
	printError(page.DisplayPostNew(res))
}

func getEditablePost(res http.ResponseWriter, req *http.Request) {
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

	printError(page.DisplayPostEdition(res, post))
}

func createPost(res http.ResponseWriter, req *http.Request) {
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

	printError(page.DisplayPostListItem(res, renderedPost, "new"))
	printError(page.DisplayPostEdition(res, renderedPost))
}

func updatePost(res http.ResponseWriter, req *http.Request) {
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

	printError(page.DisplayPostListItem(res, renderedPost, "update"))
	printError(page.DisplayPostEdition(res, renderedPost))
}

func deletePost(res http.ResponseWriter, req *http.Request) {
	if req.FormValue("confirm-delete") != "confirm" {
		return
	}

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

	type data struct {
		Status string
		Post   blog.RenderedPost
	}

	printError(page.DisplayPostListItem(res, blog.RenderedPost{Post: blog.Post{ArticleId: id}}, "delete"))
	printError(page.DisplayPostNew(res))
}
