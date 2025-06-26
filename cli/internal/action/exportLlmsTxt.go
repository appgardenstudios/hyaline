package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type ExportLlmsTxtArgs struct {
	Current     string
	Output      string
	DocumentURI string
	Full        bool
}

func ExportLlmsTxt(args *ExportLlmsTxtArgs) error {
	slog.Info("Exporting documentation to llms.txt format")
	slog.Debug("action.ExportLlmsTxt Args", slog.Group("args",
		"current", args.Current,
		"output", args.Output,
		"documentURI", args.DocumentURI,
		"full", args.Full,
	))

	// Open current database
	absCurrentPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt could not get absolute path for current database", "current", args.Current, "error", err)
		return err
	}

	// Check that input database exists
	_, err = os.Stat(absCurrentPath)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt current database does not exist", "absCurrentPath", absCurrentPath, "error", err)
		return errors.New("current database file does not exist")
	}

	db, err := sql.Open("sqlite", absCurrentPath)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt could not open current database", "absCurrentPath", absCurrentPath, "error", err)
		return err
	}
	defer db.Close()

	// Check that output file does not exist
	absOutputPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt could not get absolute path for output", "output", args.Output, "error", err)
		return err
	}

	_, err = os.Stat(absOutputPath)
	if err == nil {
		slog.Debug("action.ExportLlmsTxt output file already exists", "absOutputPath", absOutputPath)
		return errors.New("output file already exists")
	}

	// Retrieve documentation data
	systems, err := sqlite.GetAllSystem(db)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt could not retrieve systems", "error", err)
		return err
	}

	if len(systems) == 0 {
		slog.Debug("action.ExportLlmsTxt no systems found in database")
		return errors.New("no systems found in database")
	}

	// Generate llms.txt content
	var content strings.Builder

	// Write header - use first system name as title
	if len(systems) == 1 {
		content.WriteString(fmt.Sprintf("# %s Documentation\n\n", systems[0].ID))
		content.WriteString(fmt.Sprintf("> Documentation extracted from the %s system\n\n", systems[0].ID))
	} else {
		content.WriteString("# Hyaline Documentation Export\n\n")
		content.WriteString("> Documentation extracted from multiple systems\n\n")
	}

	// Process each system
	for _, system := range systems {
		err = processSystemForExport(system, db, &content, args.DocumentURI, args.Full)
		if err != nil {
			slog.Debug("action.ExportLlmsTxt could not process system", "systemID", system.ID, "error", err)
			return err
		}
	}

	// Write output file
	err = os.WriteFile(absOutputPath, []byte(content.String()), 0644)
	if err != nil {
		slog.Debug("action.ExportLlmsTxt could not write output file", "absOutputPath", absOutputPath, "error", err)
		return err
	}

	slog.Info("Export complete", "outputPath", absOutputPath)
	return nil
}

func processSystemForExport(system *sqlite.System, db *sql.DB, content *strings.Builder, documentURIFilter string, full bool) error {
	// Get all documents for the system
	documents, err := sqlite.GetAllSystemDocumentsForSystem(system.ID, db)
	if err != nil {
		return err
	}

	if len(documents) == 0 {
		slog.Debug("processSystemForExport no documents found for system", "systemID", system.ID)
		return nil
	}

	// Filter documents by URI if specified
	if documentURIFilter != "" {
		documents = filterDocumentsByURI(documents, documentURIFilter)
		if len(documents) == 0 {
			slog.Debug("processSystemForExport no documents match URI filter", "systemID", system.ID, "filter", documentURIFilter)
			return nil
		}
	}

	// Group documents by documentation source
	docsBySource := make(map[string][]*sqlite.SystemDocument)
	for _, doc := range documents {
		docsBySource[doc.DocumentationID] = append(docsBySource[doc.DocumentationID], doc)
	}

	// Sort documentation sources for consistent output
	var sources []string
	for source := range docsBySource {
		sources = append(sources, source)
	}
	sort.Strings(sources)

	// Write system header if multiple systems
	if len(sources) > 1 || system.ID != "" {
		content.WriteString(fmt.Sprintf("## System: %s\n\n", system.ID))
	}

	// Process each documentation source
	for _, source := range sources {
		docs := docsBySource[source]
		sort.Slice(docs, func(i, j int) bool {
			return docs[i].ID < docs[j].ID
		})

		// Write documentation source header
		content.WriteString(fmt.Sprintf("### %s\n\n", strings.Title(source)))

		for _, doc := range docs {
			if full {
				err = writeFullDocument(doc, db, content)
			} else {
				err = writeLinkDocument(doc, content)
			}
			if err != nil {
				return err
			}
		}

		content.WriteString("\n")
	}

	return nil
}

func filterDocumentsByURI(documents []*sqlite.SystemDocument, uriFilter string) []*sqlite.SystemDocument {
	var filtered []*sqlite.SystemDocument
	
	for _, doc := range documents {
		// Create document URI in the format: documentationID/documentID
		docURI := fmt.Sprintf("%s/%s", doc.DocumentationID, doc.ID)
		
		// Support partial matching
		if strings.Contains(docURI, uriFilter) || strings.HasPrefix(docURI, uriFilter) {
			filtered = append(filtered, doc)
		}
	}
	
	return filtered
}

func writeLinkDocument(doc *sqlite.SystemDocument, content *strings.Builder) error {
	// For link format, we create a reference to the document
	docPath := fmt.Sprintf("%s/%s", doc.DocumentationID, doc.ID)
	
	// Extract title from document content if possible
	title := extractDocumentTitle(doc.ExtractedData)
	if title == "" {
		title = doc.ID
	}
	
	content.WriteString(fmt.Sprintf("- [%s](%s): %s documentation\n", title, docPath, doc.Type))
	return nil
}

func writeFullDocument(doc *sqlite.SystemDocument, db *sql.DB, content *strings.Builder) error {
	// Write document header
	content.WriteString(fmt.Sprintf("#### %s\n\n", doc.ID))
	
	// Write document content
	if doc.ExtractedData != "" {
		content.WriteString(doc.ExtractedData)
		content.WriteString("\n\n")
	}
	
	// Get and write sections in order
	sections, err := sqlite.GetAllSystemSectionsForDocument(doc.ID, doc.DocumentationID, doc.SystemID, db)
	if err != nil {
		return err
	}
	
	for _, section := range sections {
		if section.ExtractedData != "" {
			// Calculate heading level based on section depth
			headingLevel := calculateSectionDepth(section.ID) + 4 // Start at #### since doc is already ####
			headingPrefix := strings.Repeat("#", headingLevel)
			
			content.WriteString(fmt.Sprintf("%s %s\n\n", headingPrefix, section.Name))
			content.WriteString(section.ExtractedData)
			content.WriteString("\n\n")
		}
	}
	
	return nil
}

func extractDocumentTitle(content string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ")
		}
	}
	return ""
}

func calculateSectionDepth(sectionID string) int {
	// Count the number of # characters in the section ID to determine depth
	// Section IDs follow pattern: "document.md#Section#Subsection#Sub-subsection"
	parts := strings.Split(sectionID, "#")
	return len(parts) - 1 // Subtract 1 because first part is document name
}