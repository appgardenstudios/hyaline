package docs

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
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
			Path:     d.FsOptions.Path,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not insert docs", "error", err, "doc", d.ID)
			return err
		}

		// Extract based on the extractor
		switch d.Extractor {
		case config.ExtractorFs:
			err = ExtractCurrentFs(system.ID, &d, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract docs using fs extractor", "error", err, "doc", d.ID)
				return
			}
		case config.ExtractorGit:
			err = ExtractCurrentGit(system.ID, &d, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract docs using git extractor", "error", err, "doc", d.ID)
				return
			}
		default:
			slog.Debug("docs.ExtractCurrent unknown extractor", "extractor", d.Extractor.String(), "doc", d.ID)
			return errors.New("Unknown Extractor '" + d.Extractor.String() + "' for doc " + d.ID)
		}
	}

	return
}

func ExtractCurrentFs(systemID string, d *config.Doc, db *sql.DB) (err error) {
	// Get our absolute path
	absPath, err := filepath.Abs(d.FsOptions.Path)
	if err != nil {
		slog.Debug("docs.ExtractCurrent could not determine absolute docs path", "error", err, "path", d.FsOptions.Path)
		return err
	}
	absPath += string(os.PathSeparator)
	slog.Debug("docs.ExtractCurrent extracting docs from path", "absPath", absPath)

	// Get our root FS
	// We use a root FS so symlinks and relative paths don't escape our path
	// https://pkg.go.dev/os@go1.24.1#Root
	root, err := os.OpenRoot(absPath)
	if err != nil {
		slog.Debug("code.ExtractCurrentFs could not open fs root", "error", err, "path", d.FsOptions.Path)
		return err
	}
	fsRoot := root.FS()

	// Our set of files (as a map so we don't get dupes)
	docs := map[string]struct{}{}

	// Loop through our includes and get files
	for _, include := range d.Include {
		slog.Debug("docs.ExtractCurrent extracting docs using include", "include", include, "doc", d.ID)

		// Get matched docs
		matches, err := doublestar.Glob(fsRoot, include)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not find docs files with include", "include", include, "error", err)
			return err
		}

		// Loop through docs and add those that aren't in our excludes
		for _, doc := range matches {
			// See if we have a match for at least one of our excludes
			match := false
			for _, exclude := range d.Exclude {
				match = doublestar.MatchUnvalidated(exclude, doc)
				if match {
					slog.Debug("docs.ExtractCurrent doc excluded", "doc", doc, "exclude", exclude)
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
		// Get file rawData
		rawData, err := fs.ReadFile(fsRoot, doc)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not read doc file", "error", err, "doc", doc)
			return err
		}
		// Calculate our relative path to the document path
		relativePath := strings.TrimPrefix(doc, absPath)

		// Extract and clean data (trim whitespace and remove carriage returns)
		var extractedData string
		switch d.Type {
		case config.DocTypeHTML:
			extractedData, err = extractHTMLDocument(string(rawData), d.HTML.Selector)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract html document", "error", err, "doc", doc)
				return err
			}
		default:
			extractedData = strings.TrimSpace(string(rawData))
		}
		extractedData = strings.ReplaceAll(extractedData, "\r", "")

		// Insert our document
		err = sqlite.InsertDocument(sqlite.Document{
			ID:              relativePath,
			DocumentationID: d.ID,
			SystemID:        systemID,
			Type:            d.Type.String(),
			Action:          "",
			RawData:         string(rawData),
			ExtractedData:   extractedData,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not insert document", "error", err)
			return err
		}

		// Get and insert sections
		sections := getMarkdownSections(strings.Split(extractedData, "\n"))
		err = insertMarkdownSectionAndChildren(sections, 0, relativePath, d.ID, systemID, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not insert section", "error", err)
			return err
		}
	}

	return
}

func ExtractCurrentGit(systemID string, d *config.Doc, db *sql.DB) (err error) {
	return
}

func insertMarkdownSectionAndChildren(s *section, order int, documentId string, documentationId string, systemId string, db *sql.DB) error {
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
		Name:            s.Name,
		ParentID:        parentId,
		PeerOrder:       order,
		ExtractedData:   strings.TrimSpace(s.Content),
	}, db)
	if err != nil {
		slog.Debug("docs.insertSectionAndChildren could not insert section", "error", err)
		return err
	}

	// Insert children
	for i, child := range s.Children {
		err = insertMarkdownSectionAndChildren(child, i, documentId, documentationId, systemId, db)
		if err != nil {
			return err
		}
	}

	return nil
}
