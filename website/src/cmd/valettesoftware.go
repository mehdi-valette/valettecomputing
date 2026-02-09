package main

import (
	"log"
	"net/http"
	"os"

	"valette.software/internal/blog"
	"valette.software/internal/config"
	"valette.software/internal/i18n"
	"valette.software/internal/page"
	"valette.software/internal/router"
)

func main() {
	config.Init()
	page.Init()
	blog.Init()
	i18n.Init()

	root := router.Build()

	http.Handle("/", root)

	url := buildListenUrl()

	println("server listening on", url)
	err := http.ListenAndServe(url, nil)

	if err != nil {
		log.Fatal(err)
	}

	println("server closed")
}

func buildListenUrl() string {
	port := "80"

	for _, arg := range os.Args {
		if len(arg) > 7 && arg[0:7] == "--port=" {
			port = arg[7:]
		}
	}

	return "0.0.0.0:" + port
}
