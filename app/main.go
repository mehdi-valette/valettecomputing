package main

import (
	"embed"
	"html/template"
	"net/http"
)

//go:embed *.html
var fsTemplate embed.FS

var templates *template.Template

func main() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "*.html")

	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})

	println("server listening")
	http.ListenAndServe("0.0.0.0:80", nil)
	println("server closed")
}
