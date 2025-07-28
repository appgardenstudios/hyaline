package extract

import (
	"context"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/andybalholm/cascadia"
	"golang.org/x/net/html"
)

func extractHtml(id string, sourceID string, rawData []byte, options *config.DocumentationSourceOptions, db *sqlite.Queries) error {
	// Parse the raw HTML, which will clean it up a bit and ensure it is well formatted
	rootNode, err := html.Parse(strings.NewReader(string(rawData)))
	if err != nil {
		slog.Debug("extract.extractHtml could parse html", "error", err)
		return err
	}

	// Extract the documentation section using a query (default to body)
	query := "body"
	if options.Selector != "" {
		query = options.Selector
	}
	sel, err := cascadia.Parse(query)
	if err != nil {
		slog.Debug("extract.extractHtml could parse selector", "selector", query, "error", err)
		return err
	}
	docNode := cascadia.Query(rootNode, sel)
	if docNode == nil {
		err = errors.New("could not find documentation in html document using: " + query)
		slog.Debug("extract.extractHtml could not find documentation in html using selector", "selector", query, "error", err)
		return err
	}

	// Render the node to an html string
	var b strings.Builder
	err = html.Render(&b, docNode)
	if err != nil {
		slog.Debug("extract.extractHtml could not render documentation", "error", err)
		return err
	}
	cleanHTML := b.String()

	// Convert the html to markdown
	markdown, err := htmltomarkdown.ConvertString(cleanHTML)
	if err != nil {
		slog.Debug("extract.extractHtml could not convert html to markdown", "error", err)
		return err
	}

	// Insert document
	err = db.InsertDocument(context.Background(), sqlite.InsertDocumentParams{
		ID:            id,
		SourceID:      sourceID,
		Type:          config.DocTypeHTML.String(),
		Purpose:       "",
		RawData:       string(rawData),
		ExtractedData: markdown,
	})
	if err != nil {
		slog.Debug("extract.extractHtml could not insert Document", "error", err)
		return err
	}

	// Extract/insert sections
	err = extractSections(id, sourceID, markdown, db)
	if err != nil {
		slog.Debug("extract.extractMd could not extract sections", "error", err)
		return err
	}

	return nil
}
