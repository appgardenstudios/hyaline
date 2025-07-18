package docs

import (
	"database/sql"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"log/slog"
	"slices"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/utils/merkletrie"
)

func ExtractChange(system *config.System, head string, headRef string, base string, baseRef string, documentationIDs []string, db *sql.DB) (err error) {
	// Process each docs source
	for _, d := range system.DocumentationSources {
		// Only extract if this documentation ID is passed in
		if !slices.Contains(documentationIDs, d.ID) {
			slog.Debug("docs.ExtractChange skipping non-included documentation source", "system", system.ID, "doc", d.ID)
			continue
		}

		// Only extract changed code for git sources
		if d.Extractor.Type != config.ExtractorTypeGit {
			slog.Debug("docs.ExtractChange skipping non-git documentation source", "system", system.ID, "doc", d.ID)
			continue
		}
		slog.Debug("docs.ExtractChange extracting docs", "system", system.ID, "doc", d.ID, "head", head, "headRef", headRef, "base", base, "baseRef", baseRef)

		// Get document path
		path := d.Extractor.Options.Path
		if path == "" {
			path = d.Extractor.Options.Repo
		}

		// Insert Documentation
		err = sqlite.InsertSystemDocumentation(sqlite.SystemDocumentation{
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
		r, err = repo.GetRepo(d.Extractor.Options)
		if err != nil {
			slog.Debug("code.ExtractChange could not get repo", "error", err)
			return
		}

		// Resolve head and base references
		resolvedHead, err := repo.ResolveRef(r, head, headRef)
		if err != nil {
			slog.Debug("docs.ExtractChange could not resolve head reference", "error", err)
			return err
		}

		resolvedBase, err := repo.ResolveRef(r, base, baseRef)
		if err != nil {
			slog.Debug("docs.ExtractChange could not resolve base reference", "error", err)
			return err
		}

		// Get our diffs
		diff, err := repo.GetDiff(r, *resolvedHead, *resolvedBase)
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
				if config.PathIsIncluded(change.To.Name, d.Extractor.Include, d.Extractor.Exclude) {
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
						extractedData, err = extractHTMLDocument(string(bytes), d.Options.Selector)
						if err != nil {
							slog.Debug("docs.ExtractCurrentFs could not extract html document", "error", err, "doc", change.To.Name)
							return err
						}
					default:
						extractedData = strings.TrimSpace(string(bytes))
					}
					extractedData = strings.ReplaceAll(extractedData, "\r", "")

					err = sqlite.InsertSystemDocument(sqlite.SystemDocument{
						ID:              change.To.Name,
						DocumentationID: d.ID,
						SystemID:        system.ID,
						Type:            d.Type.String(),
						Action:          sqlite.MapAction(action, change.From.Name, change.To.Name),
						OriginalID:      change.From.Name,
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
				if config.PathIsIncluded(change.To.Name, d.Extractor.Include, d.Extractor.Exclude) {
					slog.Debug("docs.ExtractChange inserting document", "document", change.From.Name, "action", action)
					err = sqlite.InsertSystemDocument(sqlite.SystemDocument{
						ID:              change.From.Name,
						DocumentationID: d.ID,
						SystemID:        system.ID,
						Type:            d.Type.String(),
						Action:          sqlite.ActionDelete,
						OriginalID:      "",
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
