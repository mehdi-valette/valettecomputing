package blog

import (
	"bytes"
	"embed"
	"io"
	"io/fs"
	"log"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

//go:embed articles
var fsBlog embed.FS

var resources fs.FS
var htmlRenderer *html.Renderer
var mdExtensions parser.Extensions

func Init() {
	var err error

	resources, err = fs.Sub(fsBlog, "articles")
	mdExtensions = parser.CommonExtensions

	if err != nil {
		log.Fatal("the blog articles cannot be found")
	}

	htmlRenderer = html.NewRenderer(html.RendererOptions{})
}

func Render(output io.Writer, name string) error {
	fileHandler, err := resources.Open(name + ".md")

	if err != nil {
		return err
	}

	mdFile := bytes.NewBuffer(nil)
	_, err = io.Copy(mdFile, fileHandler)

	if err != nil {
		return err
	}

	htmlFile := markdown.Render(
		parser.NewWithExtensions(mdExtensions).Parse(mdFile.Bytes()),
		htmlRenderer,
	)

	_, err = output.Write(htmlFile)

	return err
}
