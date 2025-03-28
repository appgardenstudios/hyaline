package docs

import (
	"errors"
	"fmt"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/antchfx/htmlquery"
)

func extractHTMLDocument(rawHTML string, selector string) (markdown string, err error) {
	// Parse the raw HTML, which will clean it up a bit and ensure it is well formatted
	rootNode, err := htmlquery.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return
	}

	// Extract the documentation
	expr := "//body"
	if selector != "" {
		expr = fmt.Sprintf("//%s", selector)
	}
	docNode, err := htmlquery.Query(rootNode, expr)
	if err != nil {
		return
	}
	if docNode == nil {
		err = errors.New("could not find body tag in html document")
		return
	}
	cleanHTML := htmlquery.OutputHTML(docNode, true)

	// Convert it to markdown
	markdown, err = htmltomarkdown.ConvertString(cleanHTML)
	if err != nil {
		return
	}

	return
}
