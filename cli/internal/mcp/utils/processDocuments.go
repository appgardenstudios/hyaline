package utils

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	giturls "github.com/whilp/git-urls"
)

// Results holds the result of document processing
type Results struct {
	Total  int
	Result strings.Builder
}

// ProcessDocuments processes documents matching the filter
// If includeContent is true, includes document content; otherwise just metadata
func ProcessDocuments(documentationData *DocumentationData, filter *DocumentURI, includeContent bool) *Results {
	results := &Results{}

	results.Result.WriteString("<systems>\n")

	for _, system := range documentationData.Systems {
		if filter == nil || filter.SystemID == "" || filter.SystemID == system.System.ID {
			processSystem(results, &system, filter, includeContent)
		}
	}

	results.Result.WriteString("</systems>")

	return results
}

func processSystem(results *Results, system *System, filter *DocumentURI, includeContent bool) {
	systemURI := &DocumentURI{
		SystemID: system.System.ID,
	}
	fmt.Fprintf(&results.Result, "  <system id=\"%s\">\n", systemURI.String())
	results.Result.WriteString("    <documentation>\n")

	for _, documentation := range system.Documentation {
		if filter == nil || filter.DocumentationID == "" || filter.DocumentationID == documentation.Documentation.ID {
			processDocumentationSource(results, &documentation, filter, includeContent)
		}
	}

	results.Result.WriteString("    </documentation>\n")
	results.Result.WriteString("  </system>\n")
}

func processDocumentationSource(results *Results, documentation *Documentation, filter *DocumentURI, includeContent bool) {
	documentationURI := &DocumentURI{
		SystemID:        documentation.Documentation.SystemID,
		DocumentationID: documentation.Documentation.ID,
	}
	fmt.Fprintf(&results.Result, "      <documentation_source id=\"%s\">\n", documentationURI.String())
	results.Result.WriteString("        <documents>\n")

	for _, document := range documentation.Documents {
		if filter == nil || filter.DocumentPath == "" || strings.HasPrefix(document.Document.ID, filter.DocumentPath) {
			processDocument(results, documentation, &document, includeContent)
			results.Total++
		}
	}

	results.Result.WriteString("        </documents>\n")
	results.Result.WriteString("      </documentation_source>\n")
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

func processDocument(results *Results, documentation *Documentation, document *Document, includeContent bool) {
	documentURI := &DocumentURI{
		SystemID:        document.Document.SystemID,
		DocumentationID: document.Document.DocumentationID,
		DocumentPath:    document.Document.ID,
	}
	fullURI := documentURI.String()

	fmt.Fprintf(&results.Result, "          <document id=\"%s\">\n", fullURI)

	// Generate source URL using the passed documentation
	sourceURL := generateSourceURL(documentation.Documentation.Path, document.Document.ID)
	fmt.Fprintf(&results.Result, "            <source>%s</source>\n", sourceURL)

	if includeContent {
		results.Result.WriteString("            <document_content>\n")

		// Use ExtractedData if available, otherwise RawData
		content := document.Document.ExtractedData
		if content == "" {
			content = document.Document.RawData
		}

		results.Result.WriteString(content)
		results.Result.WriteString("\n            </document_content>\n")
	} else {
		results.Result.WriteString("            <sections>\n")

		processSectionsForDocument(results, document)

		results.Result.WriteString("            </sections>\n")
	}

	results.Result.WriteString("          </document>\n")
}

func processSectionsForDocument(results *Results, document *Document) {
	for _, section := range document.Sections {
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
