package extract

import (
	"context"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"

	"gopkg.in/yaml.v3"
)

func extractMd(id string, sourceID string, rawData []byte, options *config.ExtractorOptions, db *sqlite.Queries) error {
	// Clean up raw data
	extractedData := strings.TrimSpace(string(rawData))
	extractedData = strings.ReplaceAll(extractedData, "\r", "")

	// Extract purpose (if not disabled)
	purpose := ""
	purposeKey := options.PurposeKey
	if !options.DisablePurposeExtraction {
		if purposeKey == "" {
			purposeKey = "purpose"
		}
		purpose = extractMdDocumentPurpose(string(rawData), purposeKey)
	}
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
	err = extractSections(id, sourceID, extractedData, !options.DisablePurposeExtraction, purposeKey, db)
	if err != nil {
		slog.Debug("extract.extractMd could not extract sections", "error", err)
		return err
	}

	return nil
}

// Extract purpose statement from a document. If not found return <blank>
func extractMdDocumentPurpose(document string, purposeKey string) string {
	lines := strings.Split(document, "\n")
	metadata := ""
	// Only support frontMatter extraction if the document starts with frontMatter
	if len(lines) > 0 && strings.HasPrefix(lines[0], "---") {
		metadata = extractFrontMatter(lines)
	}

	// Only support HTML comment extraction if the document starts with an HTML comment
	if len(lines) > 0 && strings.HasPrefix(lines[0], "<!--") {
		metadata = extractHTMLComment(lines)
	}

	return extractPurpose(metadata, purposeKey)
}

func extractFrontMatter(lines []string) string {
	// Return early if we know we don't have frontmatter
	if len(lines) < 3 || !strings.HasPrefix(lines[0], "---") {
		return ""
	}

	// Get and return contents
	contents := []string{}
	for _, line := range lines[1:] {
		if strings.HasPrefix(line, "---") {
			break
		}
		contents = append(contents, line)
	}

	return strings.TrimSpace(strings.Join(contents, "\n"))
}

func extractHTMLComment(lines []string) string {
	// Return early if we don't have a comment
	if len(lines) < 1 || !strings.HasPrefix(lines[0], "<!--") {
		return ""
	}

	contents := []string{}
	// Handle remainder of first line
	loc := strings.Index(lines[0][4:], "-->")
	if loc > -1 {
		return strings.TrimSpace(lines[0][4 : loc+4])
	} else {
		contents = append(contents, lines[0][4:])
	}

	// Go until we see the end of the comment
	for _, line := range lines[1:] {
		loc = strings.Index(line, "-->")
		if loc > -1 {
			contents = append(contents, line[:loc])
			break
		} else {
			contents = append(contents, line)
		}
	}

	return strings.TrimSpace(strings.Join(contents, "\n"))
}

// Extract purpose (identified by key) from a yaml formatted metadata string
func extractPurpose(metadata string, key string) string {
	purposeStruct := map[string]string{}
	err := yaml.Unmarshal([]byte(metadata), &purposeStruct)
	if err != nil {
		slog.Debug("extract.extractPurpose could not parse yaml metadata", "error", err)
		return ""
	}

	return purposeStruct[key]
}
