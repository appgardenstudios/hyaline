package docs

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/mattn/go-zglob"
)

func ExtractCurrent(system *config.System, db *sql.DB) (err error) {
	// Process each docs source
	for _, d := range system.Docs {
		slog.Debug("docs.ExtractCurrent extracting docs", "system", system, "docs", d.ID)
		// Insert Documentation
		documentationId := system.ID + "-" + d.ID
		err = sqlite.InsertCurrentDocumentation(sqlite.CurrentDocumentation{
			ID:       documentationId,
			SystemID: system.ID,
			Type:     d.Type,
			Path:     d.Path,
		}, db)

		// Get our absolute path
		absPath, err := filepath.Abs(d.Path)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not determine absolute docs path", "error", err, "path", d.Path)
			return err
		}
		absPath += string(os.PathSeparator)
		slog.Debug("docs.ExtractCurrent extracting docs from path", "absPath", absPath)

		// Get files from our fully qualified glob path
		glob := filepath.Join(absPath, d.Glob)
		files, err := zglob.Glob(glob)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not find doc files with glob", "error", err)
			return err
		}
		slog.Debug("docs.ExtractCurrent found the following doc file matches using glob", "glob", glob, "matches", files)

		// Insert documents/sections
		for _, file := range files {
			contents, err := os.ReadFile(file)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not read doc file", "error", err, "file", file)
				return err
			}
			relativePath := strings.TrimPrefix(file, absPath)
			err = sqlite.InsertCurrentDocument(sqlite.CurrentDocument{
				ID:              relativePath,
				DocumentationID: documentationId,
				SystemID:        system.ID,
				RelativePath:    relativePath,
				Format:          d.Type,
				RawData:         string(contents),
				ExtractedText:   extractMarkdownText(contents),
			}, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not insert document", "error", err)
				return err
			}

			// Get and insert sections
			cleanContent := strings.ReplaceAll(string(contents), "\r", "")
			sections := getMarkdownSections(strings.Split(cleanContent, "\n"))
			err = insertSectionAndChildren(sections, 0, relativePath, documentationId, system.ID, d.Type, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not insert section", "error", err)
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
		slog.Debug("docs.insertSectionAndChildren could not insert section", "error", err)
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
