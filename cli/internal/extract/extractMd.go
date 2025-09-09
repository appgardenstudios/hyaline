package extract

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

func extractMd(id string, sourceID string, rawData []byte, db *sqlite.Queries) error {
	// Clean up raw data
	extractedData := strings.TrimSpace(string(rawData))
	extractedData = strings.ReplaceAll(extractedData, "\r", "")

	// Extract purpose
	// TODO get purpose from key
	purpose := extractMdDocumentPurpose(string(rawData), "purpose")

	// Insert document
	err := db.InsertDocument(context.Background(), sqlite.InsertDocumentParams{
		ID:            id,
		SourceID:      sourceID,
		Type:          config.DocTypeMarkdown.String(),
		Purpose:       purpose,
		RawData:       string(rawData),
		ExtractedData: extractedData,
	})
	if err != nil {
		slog.Debug("extract.extractMd could not insert document", "error", err)
		return err
	}

	// Extract/insert sections
	err = extractSections(id, sourceID, extractedData, db)
	if err != nil {
		slog.Debug("extract.extractMd could not extract sections", "error", err)
		return err
	}

	return nil
}

func extractMdDocumentPurpose(document string, purposeKey string) string {
	key := purposeKey + ":"
	parts := strings.Split(document, "\n")
	// Only support frontmatter extraction if the document starts with frontmatter
	if len(parts) >= 3 && strings.HasPrefix(parts[0], "---") {
		inFrontmatter := true
		for _, line := range parts[1:] {
			trimmedLine := strings.TrimSpace(line)
			if strings.HasPrefix(trimmedLine, "---") {
				inFrontmatter = false
			}
			if inFrontmatter && strings.HasPrefix(trimmedLine, key) {
				purpose := strings.TrimSpace(trimmedLine[len(key):])
				if len(purpose) > 0 {
					return purpose
				}
			}
		}
	}

	// Only support comment extraction if the document starts with html comment
	if len(parts) >= 1 && strings.HasPrefix(parts[0], "<!--") {
		// Handle single line comment
		// Ensuring that --> is fully after <!--
		loc := strings.Index(parts[0][4:], "-->")
		if loc > -1 {
			remainder := strings.TrimSpace(parts[0][4 : loc+4]) // +4 because loc starts and the end of <!--
			if strings.HasPrefix(remainder, key) {
				purpose := strings.TrimSpace(remainder[len(key):])
				if len(purpose) > 0 {
					return purpose
				}
			}
		} else {
			// Handle multi-line comment
			// See if purpose is on the first line
			remainder := strings.TrimSpace(parts[0][4:])
			if strings.HasPrefix(remainder, key) {
				purpose := strings.TrimSpace(remainder[len(key):])
				if len(purpose) > 0 {
					return purpose
				}
			}
			if len(parts) > 1 {
				inComment := true
				for _, line := range parts[1:] {
					trimmedLine := strings.TrimSpace(line)
					loc := strings.Index(trimmedLine, "-->")
					if loc > -1 {
						// See if purpose is on this line before the comment end
						if strings.HasPrefix(trimmedLine, key) {
							purpose := strings.TrimSpace(trimmedLine[len(key):loc])
							if len(purpose) > 0 {
								return purpose
							}
						}
						inComment = false
					}
					if inComment && strings.HasPrefix(trimmedLine, key) {
						purpose := strings.TrimSpace(trimmedLine[len(key):])
						if len(purpose) > 0 {
							return purpose
						}
					}
				}
			}

			// Else search each line
			// TODO
		}

	}

	return ""
}
