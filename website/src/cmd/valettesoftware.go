package main

import (
	"errors"
	"log"
	"net/http"
	"os"

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

	http.Handle("GET /static/", http.StripPrefix("/static/", static.Serve()))

	http.HandleFunc("GET /", func(res http.ResponseWriter, req *http.Request) {
		err := page.DisplayIndex(res)

		if err != nil {
			log.Print("couldn't display the index page")
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
