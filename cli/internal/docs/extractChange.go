package docs

import (
	"database/sql"
	"hyaline/internal/config"
)

func ExtractChange(system *config.System, head string, base string, db *sql.DB) (err error) {
	// // Process each docs source
	// for _, d := range system.Docs {
	// 	slog.Debug("docs.ExtractChange extracting docs", "system", system, "docs", d.ID, "head", head, "base", base)
	// 	// Insert Documentation
	// 	documentationId := system.ID + "-" + d.ID
	// 	err = sqlite.InsertChangeDocumentation(sqlite.ChangeDocumentation{
	// 		ID:       documentationId,
	// 		SystemID: system.ID,
	// 		Type:     d.Type,
	// 		Path:     d.Path,
	// 	}, db)

	// 	// Get our diffs
	// 	diffs, err := repo.GetDiff(d.Path, head, base)
	// 	if err != nil {
	// 		slog.Debug("docs.ExtractChange could not get diff", "error", err)
	// 		return err
	// 	}

	// 	// Get our glob
	// 	glob, err := zglob.New(d.Glob)
	// 	if err != nil {
	// 		slog.Debug("docs.ExtractChange could not instantiate preset glob", "error", err, "glob", d.Glob)
	// 		return err
	// 	}

	// 	// Load any files in the diff that match our preset
	// 	for _, diff := range diffs {
	// 		slog.Debug("docs.ExtractChange processing diff", "diff", diff.String())
	// 		action, err := diff.Action()
	// 		if err != nil {
	// 			slog.Debug("docs.ExtractChange could not retrieve action for diff", "error", err, "diff", diff)
	// 			return err
	// 		}
	// 		_, to, err := diff.Files()
	// 		if err != nil {
	// 			slog.Debug("docs.ExtractChange could not retrieve files for diff", "error", err, "diff", diff)
	// 			return err
	// 		}
	// 		switch action {
	// 		case merkletrie.Insert:
	// 			fallthrough
	// 		case merkletrie.Modify:
	// 			if glob.Match(diff.To.Name) {
	// 				slog.Debug("docs.ExtractChange inserting document", "document", diff.To.Name, "action", action)
	// 				bytes, err := repo.GetBlobBytes(to.Blob)
	// 				if err != nil {
	// 					slog.Debug("docs.ExtractChange could not retrieve blob from diff", "error", err)
	// 					return err
	// 				}
	// 				err = sqlite.InsertChangeDocument(sqlite.ChangeDocument{
	// 					ID:              diff.To.Name,
	// 					DocumentationID: documentationId,
	// 					SystemID:        system.ID,
	// 					RelativePath:    diff.To.Name,
	// 					Format:          d.Type,
	// 					Action:          action.String(),
	// 					RawData:         string(bytes),
	// 					ExtractedText:   extractMarkdownText(bytes),
	// 				}, db)
	// 				if err != nil {
	// 					slog.Debug("docs.ExtractChange could not insert document", "error", err)
	// 					return err
	// 				}

	// 				// Get and insert sections
	// 				cleanContent := strings.ReplaceAll(string(bytes), "\r", "")
	// 				sections := getMarkdownSections(strings.Split(cleanContent, "\n"))
	// 				err = insertChangeSectionAndChildren(sections, 0, to.Name, documentationId, system.ID, d.Type, db)
	// 				if err != nil {
	// 					slog.Debug("docs.ExtractCurrent could not insert section", "error", err)
	// 					return err
	// 				}
	// 			}
	// 		case merkletrie.Delete:
	// 			if glob.Match(diff.From.Name) {
	// 				slog.Debug("docs.ExtractChange inserting document", "document", diff.From.Name, "action", action)
	// 				err = sqlite.InsertChangeDocument(sqlite.ChangeDocument{
	// 					ID:              diff.From.Name,
	// 					DocumentationID: documentationId,
	// 					SystemID:        system.ID,
	// 					RelativePath:    diff.From.Name,
	// 					Format:          d.Type,
	// 					Action:          action.String(),
	// 					RawData:         "",
	// 				}, db)
	// 				if err != nil {
	// 					slog.Debug("docs.ExtractChange could not insert document", "error", err)
	// 					return err
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	return
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
