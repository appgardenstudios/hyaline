package suggest

import (
	"fmt"
	"hyaline/internal/check"
	"hyaline/internal/sqlite"
	"strings"
)

// Eventually we should group this up and handle a document and section updates in the same call
func Change(systemID string, documentationSource string, document string, section []string, reasons []string, references []check.ChangeResultReference, pullRequests []*sqlite.PullRequest, issues []*sqlite.Issue) (string, error) {
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	// Diff(s)
	// TODO describe the structure of diffs?
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
	// TODO get existing document/section content
	existingContent := ""
	if isSection {
		// TODO Get section contents
		// TODO consider adding document info
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
	prompt.WriteString(fmt.Sprintf("Take into account the following reasons that this %s needs to be updated:\n", tagName))
	for _, reason := range reasons {
		prompt.WriteString(fmt.Sprintf("* %s\n", reason))
	}
	prompt.WriteString("\n")
	prompt.WriteString(fmt.Sprintf("Once you have determined what changes need to be made to the %s, make those changes to the %s and call the tool %s with the full contents of the %s. ", tagName, tagName, toolName, tagName))
	if existingContent != "" {
		prompt.WriteString(fmt.Sprintf("Match the voice and style of the existing %s content where possible. ", tagName))
	}

	fmt.Println(systemPrompt)
	fmt.Println(prompt.String())

	return "", nil
}
