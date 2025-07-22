package extract

import (
	"context"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"net/url"
)

type extractorCallback func(id string, data []byte) error

func Documentation(cfg *config.Extract, db *sqlite.Queries) (err error) {
	slog.Info("Extracting documentation from source", "source", cfg.Source.ID)
	ctx := context.Background()
	count := 0

	// Determine root
	root, err := getRoot(cfg)
	if err != nil {
		slog.Debug("extract.Documentation could not determine source root", "error", err)
		return
	}

	// Insert source
	err = db.InsertSource(ctx, sqlite.InsertSourceParams{
		ID:          cfg.Source.ID,
		Description: cfg.Source.Description,
		Crawler:     cfg.Crawler.Type.String(),
		Root:        root,
	})
	if err != nil {
		slog.Debug("extract.Documentation could not insert source", "error", err)
		return
	}

	// Initialize extractor callback
	extractor := func(id string, rawData []byte) error {
		count++

		// Find and call the first extractor that matches
		for _, e := range cfg.Extractors {
			if config.PathIsIncluded(id, e.Include, e.Exclude) {
				switch e.Type {
				case config.DocTypeMarkdown:
					return extractMd(id, cfg.Source.ID, rawData, db)
				case config.DocTypeHTML:
					return extractHtml(id, cfg.Source.ID, rawData, &e.Options, db)
				}
			}
		}

		// Return an error if we don't have an extractor that matches
		slog.Debug("extract.Documentation could not find extractor", "document", id)
		return fmt.Errorf("extract.Documentation could not find extractor for %s", id)
	}

	// Crawl
	switch cfg.Crawler.Type {
	case config.ExtractorTypeFs:
		err = crawlFs(&cfg.Crawler, extractor)
	case config.ExtractorTypeGit:
		err = crawlGit(&cfg.Crawler, extractor)
	case config.ExtractorTypeHttp:
		err = crawlHttp(&cfg.Crawler, extractor)
	}
	if err != nil {
		slog.Debug("extract.Documentation could not crawl", "error", err)
		return
	}

	// Add metadata
	err = addMetadata(cfg.Source.ID, cfg.Metadata, db)
	if err != nil {
		slog.Debug("extract.Documentation could not add metadata", "error", err)
		return
	}

	slog.Info("Extracted documentation", "count", count)
	return
}

func getRoot(cfg *config.Extract) (root string, err error) {
	root = cfg.Source.Root
	if root == "" {
		switch cfg.Crawler.Type {
		case config.ExtractorTypeFs:
			root = cfg.Crawler.Options.Path
		case config.ExtractorTypeGit:
			root = cfg.Crawler.Options.Repo
			if root == "" {
				root = cfg.Crawler.Options.Path
			}
		case config.ExtractorTypeHttp:
			var u *url.URL
			u, err = url.Parse(cfg.Crawler.Options.BaseURL)
			if err != nil {
				return
			}
			root = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		}
	}

	return
}
