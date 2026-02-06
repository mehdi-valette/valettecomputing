package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"valette.software/internal/blog"
	"valette.software/internal/config"
	"valette.software/internal/contactform"
	"valette.software/internal/page"
	"valette.software/internal/static"
)

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
	config.ReadConfig("config.ini")
	page.Init()
	blog.Init()

	http.Handle("GET /static/", http.StripPrefix("/static/", static.Serve()))

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayIndex(res)

		if err != nil {
			log.Print("couldn't display the index page")
		}
	})

	http.HandleFunc("GET /articles/", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayArticlesSummary(res)

		if err != nil {
			log.Print("couldn't display the articles's summary page\n", err)
		}
	})

	http.HandleFunc("POST /articles", func(res http.ResponseWriter, req *http.Request) {
		newArticle := blog.NewPost{}

		decoder := json.NewDecoder(req.Body)
		err := decoder.Decode(&newArticle)

		if err != nil {
			res.WriteHeader(404)
			log.Print(err)
			return
		}

		err = blog.AddPost(newArticle)

		if err != nil {
			res.WriteHeader(500)
			log.Print(err)
		}
	})

	http.HandleFunc("GET /articles/{name}", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayArticle(res, req.PathValue("name"))

		if err != nil {
			log.Print("couldn't display the article page")
		}
	})

	http.HandleFunc("GET /agenda", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayAgenda(res)

		if err != nil {
			log.Print("couldn't display the agenda page", err)
		}
	})

	http.HandleFunc("POST /contact", contactform.HandleContactFormRequest)

	url := buildListenUrl()

	println("server listening on", url)
	err := http.ListenAndServe(url, nil)

	if err != nil {
		log.Fatal(err)
	}

	println("server closed")
}
