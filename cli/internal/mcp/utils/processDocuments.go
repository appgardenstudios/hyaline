package utils

import (
	"fmt"
	"hyaline/internal/sqlite"
	"net/url"
	"path"
	"sort"
	"strings"

	giturls "github.com/whilp/git-urls"
)

// Results holds the result of document processing
type Results struct {
	Total  int
	Result strings.Builder
}

// ProcessDocuments processes documents matching the filter and builds output
// If includeContent is true, includes document content; otherwise just metadata
func ProcessDocuments(mcpData *MCPData, filter *DocumentURI, includeContent bool) *Results {
	results := &Results{}

	// Sort systems alphabetically by ID
	systemIDs := make([]string, 0, len(mcpData.Systems))
	for systemID := range mcpData.Systems {
		systemIDs = append(systemIDs, systemID)
	}
	sort.Strings(systemIDs)

	for _, systemID := range systemIDs {
		system := mcpData.Systems[systemID]
		if filter == nil || filter.SystemID == "" || filter.SystemID == system.ID {
			processSystem(results, mcpData, system, filter, includeContent)
		}
	}

	return results
}

func processSystem(results *Results, mcpData *MCPData, system *sqlite.System, filter *DocumentURI, includeContent bool) {
	systemURI := &DocumentURI{
		SystemID: system.ID,
	}
	fmt.Fprintf(&results.Result, "  <system id=\"%s\">\n", systemURI.String())
	results.Result.WriteString("    <documentation>\n")

	processDocumentation(results, mcpData, system, filter, includeContent)

	results.Result.WriteString("    </documentation>\n")
	results.Result.WriteString("  </system>\n")
}

func processDocumentation(results *Results, mcpData *MCPData, system *sqlite.System, filter *DocumentURI, includeContent bool) {
	if documentationSources, exists := mcpData.Documentation[system.ID]; exists {
		// Sort documentation source IDs alphabetically
		documentationIDS := make([]string, 0, len(documentationSources))
		for documentationID := range documentationSources {
			documentationIDS = append(documentationIDS, documentationID)
		}
		sort.Strings(documentationIDS)

		for _, documentationID := range documentationIDS {
			documentation := documentationSources[documentationID]
			if filter == nil || filter.DocumentationID == "" || filter.DocumentationID == documentation.ID {
				processDocumentationSource(results, mcpData, documentation, filter, includeContent)
			}
		}
	}
}

func processDocumentationSource(results *Results, mcpData *MCPData, documentation *sqlite.SystemDocumentation, filter *DocumentURI, includeContent bool) {
	documentationURI := &DocumentURI{
		SystemID:        documentation.SystemID,
		DocumentationID: documentation.ID,
	}
	fmt.Fprintf(&results.Result, "      <documentation_source id=\"%s\">\n", documentationURI.String())
	results.Result.WriteString("        <documents>\n")

	processDocuments(results, mcpData, documentation, filter, includeContent)

	results.Result.WriteString("        </documents>\n")
	results.Result.WriteString("      </documentation_source>\n")
}

func processDocuments(results *Results, mcpData *MCPData, documentation *sqlite.SystemDocumentation, filter *DocumentURI, includeContent bool) {
	if documents, exists := mcpData.Documents[documentation.SystemID][documentation.ID]; exists {
		// Sort document paths alphabetically
		documentPaths := make([]string, 0, len(documents))
		for documentPath := range documents {
			documentPaths = append(documentPaths, documentPath)
		}
		sort.Strings(documentPaths)

		for _, documentPath := range documentPaths {
			document := documents[documentPath]
			if filter == nil || filter.DocumentPath == "" || strings.HasPrefix(document.ID, filter.DocumentPath) {
				processDocument(results, mcpData, document, includeContent)
				results.Total++
			}
		}
	}
}

// generateSourceURL generates a source URL from a SystemDocumentation Path and document ID
func generateSourceURL(basePath string, documentID string) string {
	// Only try to parse as a git URL if basePath ends with .git
	if strings.HasSuffix(basePath, ".git") {
		if gitURL, err := giturls.Parse(basePath); err == nil {
			// Convert git URL to web URL using proper URL construction
			sourceURL := &url.URL{
				Scheme: "https",
				Host:   gitURL.Host,
				Path:   strings.TrimSuffix(gitURL.Path, ".git"),
			}
			sourceURL.Path = path.Join(sourceURL.Path, "blob/HEAD", documentID)
			return sourceURL.String()
		}
	}

	// Try to parse as URL
	if sourceURL, err := url.Parse(basePath); err == nil && sourceURL.Scheme != "" {
		// Valid URL with scheme - use proper URL path joining
		sourceURL.Path = path.Join(sourceURL.Path, documentID)
		return sourceURL.String()
	}

	// Handle file paths - prepend file:// and use proper path joining
	sourceURL := &url.URL{
		Scheme: "file",
		Path:   path.Join(basePath, documentID),
	}
	return sourceURL.String()
}

func processDocument(results *Results, mcpData *MCPData, document *sqlite.SystemDocument, includeContent bool) {
	documentURI := &DocumentURI{
		SystemID:        document.SystemID,
		DocumentationID: document.DocumentationID,
		DocumentPath:    document.ID,
	}
	fullURI := documentURI.String()

	if includeContent {
		// Format for get_documents: include content
		fmt.Fprintf(&results.Result, "          <document id=\"%s\">\n", fullURI)

		// Get the SystemDocumentation to obtain the Path
		if docData, exists := mcpData.Documentation[document.SystemID][document.DocumentationID]; exists {
			sourceURL := generateSourceURL(docData.Path, document.ID)
			fmt.Fprintf(&results.Result, "            <source>%s</source>\n", sourceURL)
		}

		results.Result.WriteString("            <document_content>\n")

		// Use ExtractedData if available, otherwise RawData
		content := document.ExtractedData
		if content == "" {
			content = document.RawData
		}

		results.Result.WriteString(content)
		results.Result.WriteString("\n            </document_content>\n          </document>\n")
	} else {
		// Format for list_documents: metadata only
		fmt.Fprintf(&results.Result, "          <document id=\"%s\">\n", fullURI)
		results.Result.WriteString("            <sections>\n")

		processSectionsForDocument(results, mcpData, document)

		results.Result.WriteString("            </sections>\n")
		results.Result.WriteString("          </document>\n")
	}
}

func processSectionsForDocument(results *Results, mcpData *MCPData, document *sqlite.SystemDocument) {
	if sections, exists := mcpData.Sections[document.SystemID][document.DocumentationID][document.ID]; exists && len(sections) > 0 {
		for _, section := range sections {
			// Skip sections with empty names
			// TODO: Should be able to remove this once https://github.com/appgardenstudios/hyaline/issues/157 is done
			if section.Name == "" {
				continue
			}

			fmt.Fprintf(&results.Result, "              <section>\n")
			fmt.Fprintf(&results.Result, "                <name>%s</name>\n", section.Name)
			results.Result.WriteString("              </section>\n")
		}
	}
}
