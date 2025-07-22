package extract

import (
	"hyaline/internal/config"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/bmatcuk/doublestar/v4"
)

func crawlFs(cfg *config.ExtractCrawler, cb extractorCallback) error {
	// Get our absolute path
	absPath, err := filepath.Abs(cfg.Options.Path)
	if err != nil {
		slog.Debug("extract.crawlFs could not determine absolute docs path", "error", err, "path", cfg.Options.Path)
		return err
	}
	slog.Info("Crawling documentation using fs", "absPath", absPath)
	slog.Debug("extract.crawlFs crawling docs from path", "absPath", absPath)

	// Get our root FS
	// We use a root FS so symlinks and relative paths don't escape our path
	// https://pkg.go.dev/os@go1.24.1#Root
	root, err := os.OpenRoot(absPath)
	if err != nil {
		slog.Debug("extract.crawlFs could not open fs root", "error", err, "path", cfg.Options.Path)
		return err
	}
	fsRoot := root.FS()

	// Our set of files (as a map so we don't get dupes)
	docs := map[string]struct{}{}

	// Loop through our includes and get files
	for _, include := range cfg.Include {
		slog.Debug("extract.crawlFs crawling docs using include", "include", include)

		// Get matched docs
		matches, err := doublestar.Glob(fsRoot, include)
		if err != nil {
			slog.Debug("extract.crawlFs could not find docs files with include", "include", include, "error", err)
			return err
		}

		// Loop through docs and add those that aren't in our excludes
		for _, doc := range matches {
			// See if we have a excludeMatch for at least one of our excludes
			excludeMatch := false
			for _, exclude := range cfg.Exclude {
				excludeMatch = doublestar.MatchUnvalidated(exclude, doc)
				if excludeMatch {
					slog.Debug("extract.crawlFs doc excluded", "doc", doc, "exclude", exclude)
					break
				}
			}
			if !excludeMatch {
				docs[doc] = struct{}{}
			}
		}
	}

	// Process docs found
	for doc := range docs {
		// Get file rawData
		rawData, err := fs.ReadFile(fsRoot, doc)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs could not read doc file", "doc", doc, "error", err)
			return err
		}

		// Call extractor callback
		err = cb(doc, rawData)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs encountered callback error", "doc", doc, "error", err)
			return err
		}
	}

	return nil
}
