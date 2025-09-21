package main

import (
	"embed"
	"errors"
	"html/template"
	"net/http"
)

//go:embed *.webp *.svg
var fsImage embed.FS

//go:embed *.html
var fsTemplate embed.FS

var templates *template.Template

func image() *http.ServeMux {
	server := http.NewServeMux()

	matrixCode, err1 := fsImage.ReadFile("matrix-code.webp")
	linuxLaptop, err2 := fsImage.ReadFile("linux-laptop.svg")
	deployment, err3 := fsImage.ReadFile("deployment.svg")
	shield, err4 := fsImage.ReadFile("shield.svg")

	if err := errors.Join(err1, err2, err3, err4); err != nil {
		panic(err.Error())
	}

	server.HandleFunc("GET /matrix-code", func(res http.ResponseWriter, req *http.Request) {
		res.Write(matrixCode)
	})

	server.HandleFunc("GET /linux-laptop", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("content-type", "image/svg+xml")
		res.Write(linuxLaptop)
	})

	server.HandleFunc("GET /deployment", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("content-type", "image/svg+xml")
		res.Write(deployment)
	})

	server.HandleFunc("GET /shield", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Add("content-type", "image/svg+xml")
		res.Write(shield)
	})

	return server
}

func main() {
	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "*.html")

	if err != nil {
		panic(err.Error())
	}

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})

	http.Handle("GET /image/", http.StripPrefix("/image", image()))

	println("server listening")
	http.ListenAndServe("0.0.0.0:80", nil)
	println("server closed")
}
