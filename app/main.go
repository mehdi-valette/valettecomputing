package main

import (
	"embed"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
)

//go:embed static
var fsStatic embed.FS

//go:embed *.html
var fsTemplate embed.FS

var templates *template.Template

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

	staticServer := http.FileServer(http.FS(fsStatic))
	http.Handle("GET /static/", staticServer)

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		templates.ExecuteTemplate(res, "index.html", nil)
	})

	http.HandleFunc("POST /contact", HandleContactFormRequest)

	url := buildListenUrl()

	println("server listening on", url)
	err = http.ListenAndServe(url, nil)

	if err != nil {
		log.Fatal(err)
	}

	println("server closed")
}
