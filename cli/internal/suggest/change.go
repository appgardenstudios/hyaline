package suggest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hyaline/internal/check"
	"hyaline/internal/config"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"

	"github.com/invopop/jsonschema"
)

type changeUpdateSchema struct {
	Content string `json:"content" jsonschema:"title=content,description=The full content of the updated document or section"`
}
type changeNoUpdateNeededSchema struct {
}

// Eventually we should group this up and handle a document and section updates in the same call
func Change(systemID string, documentationSource string, document string, section []string, purpose string, reasons []string, references []check.ChangeResultReference, pullRequests []*sqlite.PullRequest, issues []*sqlite.Issue, cfg *config.LLM, currentDB *sql.DB) (suggestion string, err error) {
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	// Diff(s)
	prompt.WriteString("The contents of the <diffs> tag contain a list of unified diffs for files that were changed.\n")
	prompt.WriteString("<diffs>\n")
	for _, ref := range references {
		prompt.WriteString("  <diff>\n")
		prompt.WriteString(ref.Diff)
		prompt.WriteString("  </diff>\n")
	}
	prompt.WriteString("</diffs>\n")
	prompt.WriteString("\n")

	// Add Pull Request (if any)
	numPullRequests := len(pullRequests)
	for _, pr := range pullRequests {
		prompt.WriteString("<pull_request>\n")
		prompt.WriteString(fmt.Sprintf("  <pull_request_title>%s</pull_request_title>\n", pr.Title))
		prompt.WriteString("  <pull_request_content>\n")
		prompt.WriteString(pr.Body)
		prompt.WriteString("\n")
		prompt.WriteString("  </pull_request_content>\n")
		prompt.WriteString("</pull_request>\n")
		prompt.WriteString("\n\n")
	}

	// Add issue(s) (if any)
	numIssues := len(issues)
	for _, issue := range issues {
		prompt.WriteString("<issue>\n")
		prompt.WriteString(fmt.Sprintf("  <issue_title>%s</issue_title>\n", issue.Title))
		prompt.WriteString("  <issue_content>\n")
		prompt.WriteString(issue.Body)
		prompt.WriteString("\n")
		prompt.WriteString("  </issue_content>\n")
		prompt.WriteString("</issue>\n")
		prompt.WriteString("\n\n")
	}

	// Add document/section
	isSection := len(section) > 0
	existingContent := ""
	if isSection {
		// Get current section content
		sectionID := fmt.Sprintf("%s#%s", document, strings.Join(section, "#"))
		var originalSection *sqlite.Section
		originalSection, err = sqlite.GetSection(sectionID, document, documentationSource, systemID, currentDB)
		if err != nil {
			slog.Debug("suggest.Change could not retrieve section", "sectionID", sectionID, "error", err)
			return
		}
		if originalSection != nil {
			existingContent = originalSection.ExtractedData
		}
		// Add section
		title := section[len(section)-1]
		prompt.WriteString("<section>\n")
		prompt.WriteString(fmt.Sprintf("  <section_title>%s</section_title>\n", title))
		prompt.WriteString("  <section_content>\n")
		prompt.WriteString(existingContent)
		prompt.WriteString("\n")
		prompt.WriteString("  </section_content>\n")
		prompt.WriteString("</section>\n")
		prompt.WriteString("\n")
	} else {
		// Get current document content
		var originalDocument *sqlite.Document
		originalDocument, err = sqlite.GetDocument(document, documentationSource, systemID, currentDB)
		if err != nil {
			slog.Debug("suggest.Change could not retrieve document", "document", document, "error", err)
			return
		}
		if originalDocument != nil {
			existingContent = originalDocument.ExtractedData
		}
		// Add document
		prompt.WriteString("<document>\n")
		prompt.WriteString(fmt.Sprintf("  <document_name>%s</document_name>\n", document))
		prompt.WriteString("  <document_content>\n")
		prompt.WriteString(existingContent)
		prompt.WriteString("\n")
		prompt.WriteString("  </document_content>\n")
		prompt.WriteString("</document>\n")
		prompt.WriteString("\n")
	}

	// Add prompt
	var tagName string
	var toolName string
	const noUpdateNeededName = "no_update_needed"
	if isSection {
		tagName = "section"
		toolName = "update_section"
	} else {
		tagName = "document"
		toolName = "update_document"
	}
	prompt.WriteString("Given the set of file changes detailed in <diffs>, ")
	if numPullRequests > 0 {
		prompt.WriteString("and the contents of related pull request(s) in <pull_request>, ")
	}
	if numIssues > 0 {
		prompt.WriteString("and the contents of related issue(s) in <issue>, ")
	}
	prompt.WriteString(fmt.Sprintf("determine what changes need to be made to the %s contained in <%s>. ", tagName, tagName))
	prompt.WriteString("Be concise and accurate. ")
	if purpose != "" {
		prompt.WriteString(fmt.Sprintf("The purpose of this %s is \"%s\". ", tagName, purpose))
	}
	prompt.WriteString(fmt.Sprintf("Take into account the following reasons that this %s needs to be updated:\n", tagName))
	for _, reason := range reasons {
		prompt.WriteString(fmt.Sprintf("* %s\n", reason))
	}
	prompt.WriteString("\n")
	prompt.WriteString(fmt.Sprintf("Once you have determined what changes need to be made to the %s, make those changes to the %s and call the tool %s with the full contents of the %s. ", tagName, tagName, toolName, tagName))
	if existingContent != "" {
		prompt.WriteString(fmt.Sprintf("Match the voice and style of the existing %s content where possible. ", tagName))
	}
	prompt.WriteString(fmt.Sprintf("If no changes need to be made call the tool %s.", noUpdateNeededName))

	// Add tool(s)
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	tools := []*llm.Tool{
		{
			Name:        toolName,
			Description: fmt.Sprintf("Update the contents of the %s", tagName),
			Schema:      reflector.Reflect(&changeUpdateSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("suggest.Change - checkLLM made an update")
				// Parse the input
				var update changeUpdateSchema
				err := json.Unmarshal([]byte(input), &update)
				if err != nil {
					slog.Debug("suggest.Change - checkLLM could not parse tool call input, invalid json", "tool", toolName, "input", input, "error", err)
					return true, "", err
				}
				// Record the suggestion
				suggestion = update.Content

				return true, "", nil
			},
		},
		{
			Name:        noUpdateNeededName,
			Description: fmt.Sprintf("Identify that the %s does not need to be updated.", tagName),
			Schema:      reflector.Reflect(&changeNoUpdateNeededSchema{}),
			Callback: func(params string) (bool, string, error) {
				slog.Debug("suggest.Change - checkLLM determined no updates needed")

				// Do nothing so that we pass back a blank suggestion with no error

				// Return with done = true so we stop
				return true, "", nil
			},
		},
	}

	// Call LLM
	userPrompt := prompt.String()
	// fmt.Println(userPrompt)
	slog.Debug("suggest.Change calling the llm")
	// slog.Debug("suggest.Change calling the llm", "systemPrompt", systemPrompt, "userPrompt", userPrompt)
	_, err = llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("suggest.Change encountered an error when calling the llm", "error", err)
		return
	}

	return
}
