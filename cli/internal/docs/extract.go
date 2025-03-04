package docs

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

func ExtractCurrent(system string, cfg *config.Config, db *sql.DB) (err error) {
	// Find our target system (error if not found)
	var targetSystem *config.System
	for _, s := range cfg.Systems {
		if s.ID == system {
			targetSystem = &s
		}
	}
	if targetSystem == nil {
		// TODO better error message here
		return errors.New("system not found")
	}

	// Process each docs source
	for _, d := range targetSystem.Docs {
		// Insert Documentation
		err = sqlite.InsertCurrentDocumentation(sqlite.CurrentDocumentation{
			ID:       d.ID,
			SystemID: targetSystem.ID,
			Type:     d.Type,
			Path:     d.Path,
		}, db)

		// Get our absolute path
		absPath, err := filepath.Abs(d.Path)
		if err != nil {
			return err
		}
		absPath += string(os.PathSeparator)

		// Get files from our fully qualified glob path
		glob := filepath.Join(absPath, d.Glob)
		files, err := zglob.Glob(glob)
		if err != nil {
			return err
		}

		// Insert documents/sections
		for _, file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				return err
			}
			relativePath := strings.TrimPrefix(file, absPath)
			err = sqlite.InsertCurrentDocument(sqlite.CurrentDocument{
				ID:              relativePath,
				DocumentationID: d.ID,
				SystemID:        targetSystem.ID,
				RelativePath:    relativePath,
				Format:          d.Type,
				RawData:         string(contents),
				ExtractedText:   "TODO", // Use https://github.com/gomarkdown/markdown https://github.com/gomarkdown/markdown/blob/master/md/md_renderer.go
			}, db)
			if err != nil {
				return err
			}

			// TODO get and insert sections
		}
	}

	return
}
