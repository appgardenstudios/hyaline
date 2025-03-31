package docs

import (
	"errors"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func extractHTMLDocument(rawHTML string, selector string) (markdown string, err error) {
	// Parse the raw HTML, which will clean it up a bit and ensure it is well formatted
	rootNode, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return
	}

	// Extract the documentation section using a query (default to body)
	query := "body"
	if selector != "" {
		query = selector
	}
	sel, err := cascadia.Parse(query)
	if err != nil {
		return
	}
	docNode := cascadia.Query(rootNode, sel)
	if docNode == nil {
		err = errors.New("could not find documentation in html document using: " + query)
		return
	}

	// Render the node to an html string
	var b strings.Builder
	err = html.Render(&b, docNode)
	if err != nil {
		return
	}
	cleanHTML := b.String()

	// Convert the html to markdown
	markdown, err = htmltomarkdown.ConvertString(cleanHTML)
	if err != nil {
		return
	}

	return
}
