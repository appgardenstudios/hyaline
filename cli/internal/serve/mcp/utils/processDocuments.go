package utils

import (
	"fmt"
	"log/slog"
	"net/url"
	"path"
	"strings"

	"hyaline/internal/config"

	giturls "github.com/whilp/git-urls"
)

// Results holds processing results
type Results struct {
	Total  int
	Result strings.Builder
}

// ProcessDocuments processes documents based on URI and returns results
func ProcessDocuments(data *DocumentationData, documentURI *DocumentURI, includeContent bool) *Results {
	results := &Results{}

	// Open documents tag
	results.Result.WriteString("<documents>\n")

	// Process all sources
	for _, source := range data.Sources {
		// Filter by source ID if specified
		if documentURI.SourceID != "" && source.ID != documentURI.SourceID {
			continue
		}

		// Process documents for this source
		for _, document := range source.Documents {
			// Filter by document path if specified
			if documentURI.DocumentPath != "" {
				// Check for exact match or prefix match
				if document.ID != documentURI.DocumentPath && !strings.HasPrefix(document.ID, documentURI.DocumentPath+"/") {
					continue
				}
			}

			// Filter by tags
			if !documentURI.MatchesTags(document.Tags) {
				continue
			}

			// Process this document
			processDocument(results, &source, &document, includeContent)
		}
	}

	// Close documents tag
	results.Result.WriteString("</documents>")

	return results
}

// processDocument processes a single document and adds it to results
func processDocument(results *Results, source *Source, document *Document, includeContent bool) {
	uri := &DocumentURI{
		SourceID:     source.ID,
		DocumentPath: document.ID,
	}

	results.Result.WriteString("          <document>\n")
	fmt.Fprintf(&results.Result, "            <uri>%s</uri>\n", uri.String())

	// Generate and add source URL
	sourceURL := generateSourceURL(source.Crawler, source.Root, document.ID)
	fmt.Fprintf(&results.Result, "            <source>%s</source>\n", sourceURL)

	// Add purpose if present
	if document.Purpose != "" {
		fmt.Fprintf(&results.Result, "            <purpose>%s</purpose>\n", document.Purpose)
	}

	// Add tags if present
	if len(document.Tags) > 0 {
		results.Result.WriteString("            <tags>\n")
		writeTags(&results.Result, document.Tags, "            ")
		results.Result.WriteString("            </tags>\n")
	}

	if includeContent {
		results.Result.WriteString("            <document_content>\n")
		// Use ExtractedData if available, otherwise RawData
		content := document.ExtractedData
		if content == "" {
			content = document.RawData
		}
		results.Result.WriteString(content)
		results.Result.WriteString("\n            </document_content>\n")

		// Include sections with content
		if len(document.Sections) > 0 {
			results.Result.WriteString("            <sections>\n")
			for _, section := range document.Sections {
				processSectionWithContent(results, &section)
			}
			results.Result.WriteString("            </sections>\n")
		}
	} else {
		// Just list sections without content
		results.Result.WriteString("            <sections>\n")
		processSectionsForDocument(results, document)
		results.Result.WriteString("            </sections>\n")
	}

	results.Result.WriteString("          </document>\n")

	slog.Debug("mcp.utils.processDocument",
		"uri", uri.String(),
		"includeContent", includeContent,
		"hasExtractedData", document.ExtractedData != "",
		"hasRawData", document.RawData != "",
		"sectionsCount", len(document.Sections),
	)
}

// processSectionsForDocument processes sections for a document (list mode)
func processSectionsForDocument(results *Results, document *Document) {
	for _, section := range document.Sections {
		// Skip sections with empty names
		if section.Name == "" {
			continue
		}

		fmt.Fprintf(&results.Result, "              <section>\n")
		fmt.Fprintf(&results.Result, "                <name>%s</name>\n", section.Name)

		// Add purpose if present
		if section.Purpose != "" {
			fmt.Fprintf(&results.Result, "                <purpose>%s</purpose>\n", section.Purpose)
		}

		// Add tags if present
		if len(section.Tags) > 0 {
			results.Result.WriteString("                <tags>\n")
			writeTags(&results.Result, section.Tags, "                ")
			results.Result.WriteString("                </tags>\n")
		}

		results.Result.WriteString("              </section>\n")
	}
}

// processSectionWithContent processes a section with its content
func processSectionWithContent(results *Results, section *Section) {
	fmt.Fprintf(&results.Result, "              <section>\n")
	fmt.Fprintf(&results.Result, "                <name>%s</name>\n", section.Name)

	// Add purpose if present
	if section.Purpose != "" {
		fmt.Fprintf(&results.Result, "                <purpose>%s</purpose>\n", section.Purpose)
	}

	// Add tags if present
	if len(section.Tags) > 0 {
		results.Result.WriteString("                <tags>\n")
		writeTags(&results.Result, section.Tags, "                ")
		results.Result.WriteString("                </tags>\n")
	}

	// Add content
	if section.ExtractedData != "" {
		results.Result.WriteString("                <content>\n")
		results.Result.WriteString(section.ExtractedData)
		results.Result.WriteString("\n                </content>\n")
	}

	results.Result.WriteString("              </section>\n")
}

// generateSourceURL generates a source URL based on crawler type, root, and document ID
func generateSourceURL(crawlerType string, root string, documentID string) string {
	switch crawlerType {
	case string(config.ExtractorTypeGit):
		// For git crawler, root should be a git URL
		if gitURL, err := giturls.Parse(root); err == nil {
			// Convert git URL to web URL using proper URL construction
			sourceURL := &url.URL{
				Scheme: "https",
				Host:   gitURL.Host,
				Path:   strings.TrimSuffix(gitURL.Path, ".git"),
			}
			sourceURL.Path = path.Join(sourceURL.Path, "blob/HEAD", documentID)
			return sourceURL.String()
		}
		// If parsing fails, fall through to default

	case string(config.ExtractorTypeHttp):
		// For http crawler, root is already a web URL
		if sourceURL, err := url.Parse(root); err == nil && sourceURL.Scheme != "" {
			// Valid URL with scheme - use proper URL path joining
			sourceURL.Path = path.Join(sourceURL.Path, documentID)
			return sourceURL.String()
		}
		// If parsing fails, fall through to default

	case string(config.ExtractorTypeFs):
		// For fs crawler, root is a file system path
		// Just join the paths
		return path.Join(root, documentID)
	}

	// Default: just join as paths
	return path.Join(root, documentID)
}

// writeTags writes tags in a deterministic order (sorted by key)
func writeTags(builder *strings.Builder, tags Tags, indent string) {
	if len(tags) == 0 {
		return
	}

	// Write tags in sorted order using Tags.Keys() method
	for _, key := range tags.Keys() {
		values := tags[key]
		builder.WriteString(indent + "  <tag>\n")
		fmt.Fprintf(builder, indent+"    <key>%s</key>\n", key)
		for _, value := range values {
			fmt.Fprintf(builder, indent+"    <value>%s</value>\n", value)
		}
		builder.WriteString(indent + "  </tag>\n")
	}
}
