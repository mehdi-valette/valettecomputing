package main

import (
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed *.webp *.svg icon/banner/*
var fsImage embed.FS

//go:embed *.html
var fsTemplate embed.FS

var templates *template.Template

func image() *http.ServeMux {
	server := http.NewServeMux()

	banner, errBanner := fsImage.ReadDir("icon/banner")
	matrixCode, err1 := fsImage.ReadFile("matrix-code.webp")
	linuxLaptop, err2 := fsImage.ReadFile("linux-laptop.svg")
	deployment, err3 := fsImage.ReadFile("deployment.svg")
	shield, err4 := fsImage.ReadFile("shield.svg")

	if err := errors.Join(err1, err2, err3, err4, errBanner); err != nil {
		log.Fatal(err)
	}

	server.HandleFunc("GET /icon/banner/{img}", func(res http.ResponseWriter, req *http.Request) {
		for _, image := range banner {
			if image.Name() != req.PathValue("img") {
				continue
			}

			result, err := fsImage.ReadFile("icon/banner/" + image.Name())

			if err != nil {
				log.Fatal(err)
			}

			res.Header().Add("content-type", "image/svg+xml")
			res.Write(result)
			break
		}
	})

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

func buildListenUrl() string {
	port := "80"

	for argIndex, arg := range os.Args {
		if arg == "--port" {
			if argIndex >= len(os.Args)-1 {
				log.Fatal(errors.New("the argument --port requires a value"))
			}

			port = os.Args[argIndex+1]
		}
	}

	return "0.0.0.0:" + port
}

func main() {
	ReadConfig("config.ini")

	var err error
	templates, err = template.New("").ParseFS(fsTemplate, "*.html")

	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})

	http.Handle("GET /image/", http.StripPrefix("/image", image()))

	http.HandleFunc("POST /contact", HandleContactFormRequest)

	url := buildListenUrl()

	println("server listening on", url)
	err = http.ListenAndServe(url, nil)

	if err != nil {
		log.Fatal(err)
	}

	println("server closed")
}
