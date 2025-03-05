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
		documentationId := targetSystem.ID + "-" + d.ID
		err = sqlite.InsertCurrentDocumentation(sqlite.CurrentDocumentation{
			ID:       documentationId,
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
				DocumentationID: documentationId,
				SystemID:        targetSystem.ID,
				RelativePath:    relativePath,
				Format:          d.Type,
				RawData:         string(contents),
				ExtractedText:   extractMarkdownText(contents),
			}, db)
			if err != nil {
				return err
			}

			// Get and insert sections
			cleanContent := strings.ReplaceAll(string(contents), "\r", "")
			sections := getMarkdownSections(strings.Split(cleanContent, "\n"))
			err = insertSectionAndChildren(sections, 0, relativePath, documentationId, targetSystem.ID, d.Type, db)
			if err != nil {
				return err
			}
		}
	}

	return
}

func insertSectionAndChildren(s *section, order int, documentId string, documentationId string, systemId string, format string, db *sql.DB) error {
	// Insert this section
	parentSectionId := ""
	if s.Parent != nil {
		parentSectionId = documentId + "#" + s.Parent.Title
	}
	err := sqlite.InsertCurrentSection(sqlite.CurrentSection{
		ID:              documentId + "#" + s.Title,
		DocumentID:      documentId,
		DocumentationID: documentationId,
		SystemID:        systemId,
		ParentSectionID: parentSectionId,
		Order:           order,
		Title:           s.Title,
		Format:          format,
		RawData:         strings.TrimSpace(s.Content),
		ExtractedText:   extractMarkdownText([]byte(s.Content)),
	}, db)
	if err != nil {
		return err
	}

	// Insert children
	for i, child := range s.Children {
		err = insertSectionAndChildren(child, i, documentId, documentationId, systemId, format, db)
		if err != nil {
			return err
		}
	}

	return nil
}
