package static

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
)

//go:embed resource
var fsStatic embed.FS

func Serve() http.Handler {
	fmt.Println("serve...")

	test, err := fs.Sub(fsStatic, "resource")

	if err != nil {
		log.Fatal("the resources couldn't be embeded properly")
	}

	return http.FileServer(http.FS(test))
}
