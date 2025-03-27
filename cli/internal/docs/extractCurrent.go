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
		slog.Debug("docs.ExtractCurrent extracting docs", "system", system.ID, "docs", d.ID)
		// Insert Documentation
		err = sqlite.InsertDocumentation(sqlite.Documentation{
			ID:       d.ID,
			SystemID: system.ID,
			Type:     d.Type.String(),
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

		// Our set of files (as a map so we don't get dupes)
		docs := map[string]struct{}{}

		// Loop through our includes and get files
		for _, include := range d.Include {
			slog.Debug("docs.ExtractCurrent extracting docs using include", "include", include, "doc", d.ID)

			// Construct our includePattern and get matches
			includePattern := filepath.Join(absPath, include)
			matches, err := zglob.Glob(includePattern)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not find docs files with glob", "glob", includePattern, "error", err)
				return err
			}

			// Loop through docs and add those that aren't in our excludes
			for _, doc := range matches {
				// See if we have a match for at least one of our excludes
				match := false
				for _, exclude := range d.Exclude {
					excludePattern := filepath.Join(absPath, exclude)
					match, err = zglob.Match(excludePattern, doc)
					if err != nil {
						slog.Debug("docs.ExtractCurrent could not match exclude", "excludePattern", excludePattern, "doc", doc, "error", err)
						return err
					}
					if match {
						slog.Debug("docs.ExtractCurrent doc excluded", "doc", doc, "excludePattern", excludePattern)
						break
					}
				}
				if !match {
					docs[doc] = struct{}{}
				}
			}
		}

		// Insert docs
		for doc := range docs {
			// Get file contents
			contents, err := os.ReadFile(doc)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not read doc file", "error", err, "doc", doc)
				return err
			}
			// Calculate our relative path to the document path
			relativePath := strings.TrimPrefix(doc, absPath)

			// Extract data
			var extractedData string
			switch d.Type {
			case config.DocTypeHTML:
				extractedData = "TODO" // wire this up
			default:
				extractedData = strings.TrimSpace(string(contents))
			}

			// Insert our document
			err = sqlite.InsertDocument(sqlite.Document{
				ID:              relativePath,
				DocumentationID: d.ID,
				SystemID:        system.ID,
				Type:            d.Type.String(),
				Action:          "",
				RawData:         string(contents),
				ExtractedData:   extractedData,
			}, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not insert document", "error", err)
				return err
			}

			// Get and insert sections
			switch d.Type {
			case config.DocTypeHTML:
				// TODO
			case config.DocTypeMarkdown:
				cleanContent := strings.ReplaceAll(string(contents), "\r", "")
				sections := getMarkdownSections(strings.Split(cleanContent, "\n"))
				err = insertMarkdownSectionAndChildren(sections, 0, relativePath, d.ID, system.ID, d.Type, db)
				if err != nil {
					slog.Debug("docs.ExtractCurrent could not insert section", "error", err)
					return err
				}
			}
		}
	}

	return
}

func insertMarkdownSectionAndChildren(s *markdownSection, order int, documentId string, documentationId string, systemId string, docType config.DocType, db *sql.DB) error {
	// Insert this section
	parentId := ""
	if s.Parent != nil {
		parentId = documentId + s.Parent.FullName
	}
	err := sqlite.InsertSection(sqlite.Section{
		ID:              documentId + s.FullName,
		DocumentID:      documentId,
		DocumentationID: documentationId,
		SystemID:        systemId,
		Type:            docType.String(),
		Name:            s.Name,
		ParentID:        parentId,
		PeerOrder:       order,
		RawData:         s.Content,
		ExtractedData:   strings.TrimSpace(s.Content),
	}, db)
	if err != nil {
		slog.Debug("docs.insertSectionAndChildren could not insert section", "error", err)
		return err
	}

	// Insert children
	for i, child := range s.Children {
		err = insertMarkdownSectionAndChildren(child, i, documentId, documentationId, systemId, docType, db)
		if err != nil {
			return err
		}
	}

	return nil
}
