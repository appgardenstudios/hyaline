package docs

import (
	"errors"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/antchfx/htmlquery"
)

func extractHTMLDocument(rawHTML string) (markdown string, err error) {
	// Parse the raw rootNode, which will clean it up a bit and ensure it is well formatted
	rootNode, err := htmlquery.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return
	}

	// Extract the documentation
	// TODO use selector from config
	docNode := htmlquery.FindOne(rootNode, "//body")
	if docNode == nil {
		err = errors.New("could not find body tag in html document")
		return
	}
	cleanHTML := htmlquery.OutputHTML(docNode, true)

	markdown, err = htmltomarkdown.ConvertString(cleanHTML)
	if err != nil {
		return
	}

	return
}
