package docs

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

func ExtractChange(system *config.System, head string, base string, db *sql.DB) (err error) {
	// Process each docs source
	for _, d := range system.Docs {
		// Only extract changed code for git sources
		if d.Extractor != config.ExtractorGit {
			slog.Debug("code.ExtractChange skipping non-git code source", "system", system.ID, "doc", d.ID)
			continue
		}
		slog.Debug("docs.ExtractChange extracting docs", "system", system.ID, "doc", d.ID, "head", head, "base", base)

		// Get document path
		path := d.GitOptions.Path
		if path == "" {
			path = d.GitOptions.Repo
		}

		// Insert Documentation
		err = sqlite.InsertDocumentation(sqlite.Documentation{
			ID:       d.ID,
			SystemID: system.ID,
			Type:     d.Type.String(),
			Path:     path,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractChange could not insert documentation", "error", err, "doc", d.ID)
			return err
		}

		// Initialize go-git repo (on disk or in mem)
		var r *git.Repository
		r, err = repo.GetRepo(d.GitOptions)
		if err != nil {
			slog.Debug("code.ExtractChange could not get repo", "error", err)
			return
		}

		// Get our diffs
		diff, err := repo.GetDiff(r, head, base)
		if err != nil {
			slog.Debug("docs.ExtractChange could not get diff", "error", err)
			return err
		}

		// Load any files in the diff that match our preset
		for _, change := range diff {
			slog.Debug("docs.ExtractChange processing diff", "diff", change.String())
			action, err := change.Action()
			if err != nil {
				slog.Debug("docs.ExtractChange could not retrieve action for diff", "error", err, "diff", change)
				return err
			}
			_, to, err := change.Files()
			if err != nil {
				slog.Debug("docs.ExtractChange could not retrieve files for diff", "error", err, "diff", change)
				return err
			}
			switch action {
			case merkletrie.Insert:
				fallthrough
			case merkletrie.Modify:
				if isIncluded(change.To.Name, d.Include, d.Exclude) {
					slog.Debug("docs.ExtractChange inserting document", "document", change.To.Name, "action", action)
					bytes, err := repo.GetBlobBytes(to.Blob)
					if err != nil {
						slog.Debug("docs.ExtractChange could not retrieve blob from diff", "error", err)
						return err
					}

					// Extract and clean data (trim whitespace and remove carriage returns)
					var extractedData string
					switch d.Type {
					case config.DocTypeHTML:
						extractedData, err = extractHTMLDocument(string(bytes), d.HTML.Selector)
						if err != nil {
							slog.Debug("docs.ExtractCurrentFs could not extract html document", "error", err, "doc", change.To.Name)
							return err
						}
					default:
						extractedData = strings.TrimSpace(string(bytes))
					}
					extractedData = strings.ReplaceAll(extractedData, "\r", "")

					err = sqlite.InsertDocument(sqlite.Document{
						ID:              change.To.Name,
						DocumentationID: d.ID,
						SystemID:        system.ID,
						Type:            d.Type.String(),
						Action:          action.String(),
						RawData:         string(bytes),
						ExtractedData:   extractedData,
					}, db)
					if err != nil {
						slog.Debug("docs.ExtractChange could not insert document", "error", err)
						return err
					}

					// Get and insert sections
					sections := getMarkdownSections(strings.Split(extractedData, "\n"))
					err = insertMarkdownSectionAndChildren(sections, 0, change.To.Name, d.ID, system.ID, db)
					if err != nil {
						slog.Debug("docs.ExtractCurrentFs could not insert section", "error", err)
						return err
					}
				}
			case merkletrie.Delete:
				if isIncluded(change.To.Name, d.Include, d.Exclude) {
					slog.Debug("docs.ExtractChange inserting document", "document", change.From.Name, "action", action)
					err = sqlite.InsertDocument(sqlite.Document{
						ID:              change.From.Name,
						DocumentationID: d.ID,
						SystemID:        system.ID,
						Type:            d.Type.String(),
						Action:          action.String(),
						RawData:         "",
					}, db)
					if err != nil {
						slog.Debug("docs.ExtractChange could not insert document", "error", err)
						return err
					}
				}
			}
		}
	}

	return
}

func isIncluded(name string, includes []string, excludes []string) bool {
	for _, include := range includes {
		if doublestar.MatchUnvalidated(include, name) {
			for _, exclude := range excludes {
				if doublestar.MatchUnvalidated(exclude, name) {
					return false
				}
			}
			return true
		}
	}

	return false
}

// func insertChangeSectionAndChildren(s *section, order int, documentId string, documentationId string, systemId string, format string, db *sql.DB) error {
// 	// Insert this section
// 	parentSectionId := ""
// 	if s.Parent != nil {
// 		parentSectionId = documentId + "#" + s.Parent.Name
// 	}
// 	err := sqlite.InsertChangeSection(sqlite.ChangeSection{
// 		ID:              documentId + "#" + s.Name,
// 		DocumentID:      documentId,
// 		DocumentationID: documentationId,
// 		SystemID:        systemId,
// 		ParentSectionID: parentSectionId,
// 		Order:           order,
// 		Title:           s.Name,
// 		Format:          format,
// 		RawData:         strings.TrimSpace(s.Content),
// 		ExtractedText:   extractMarkdownText([]byte(s.Content)),
// 	}, db)
// 	if err != nil {
// 		slog.Debug("docs.insertChangeSectionAndChildren could not insert section", "error", err)
// 		return err
// 	}

// 	// Insert children
// 	for i, child := range s.Children {
// 		err = insertChangeSectionAndChildren(child, i, documentId, documentationId, systemId, format, db)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
