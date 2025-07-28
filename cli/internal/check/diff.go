package check

import (
	"encoding/json"
	"fmt"
	"hyaline/internal/code"
	"hyaline/internal/diff"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/llm"
	"log/slog"
	"strings"

	"github.com/invopop/jsonschema"
)

type Result struct {
	Source   string
	Document string
	Section  []string
	Reasons  []string
}

const checkNeedsUpdateName = "needs_update"
const checkNoUpdateNeededName = "no_update_needed"

type checkNeedsUpdateSchema struct {
	Entries []checkNeedsUpdateSchemaEntry `json:"entries" jsonschema:"title=The list of entries,description=The list of documents and/or sections that need to be updated along with the reason for each update"`
}

type checkNeedsUpdateSchemaEntry struct {
	ID     string `json:"id" jsonschema:"title=The document/section ID,description=The ID of the document and/or section that needs to be updated,example=app.1"` // TODO update examples
	Reason string `json:"reason" jsonschema:"title=The reason,description=The reason the document and/or section needs to be updated,example=This section needs to be updated because the change modifies a file that is mentioned in the reference to this section"`
}

type checkNoUpdateNeededSchema struct {
}

func Diff(files []code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue) (results []Result, err error) {
	resultMap := make(map[string][]string)

	systemPrompt := "You are a senior technical writer who writes clear and accurate documentation."
	tools := getCheckTools(func(id string, reason string) {
		entry, ok := resultMap[id]
		if ok {
			entry = append(entry, reason)
			resultMap[id] = entry
		} else {
			resultMap[id] = []string{reason}
		}
	})

	for _, file := range files {
		var prompt string
		prompt, err = formatCheckPrompt(file, documents, pr, issues)
		if err != nil {
			slog.Debug("check.Diff could not format prompt", "error", err)
			return
		}

		// TODO call llm
		slog.Debug("check.Diff calling llm", "file", file.Filename, "systemPrompt", systemPrompt, "prompt", prompt, "tools", len(tools))
	}

	// Process resultMap into results
	// TODO

	// Sort
	// TODO

	return
}

func formatCheckPrompt(file code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue) (string, error) {
	// Calculate the diff
	edits := diff.Strings(string(file.OriginalContents), string(file.Contents))
	textDiff, err := diff.ToUnified("a/"+file.OriginalFilename, "b/"+file.Filename, string(file.OriginalContents), edits, 3)
	if err != nil {
		slog.Debug("check.Diff could not generate diff", "file", file.Filename, "error", err)
		return "", err
	}

	var prompt strings.Builder

	// Add documentation context
	// https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/long-context-tips#example-quote-extraction
	formattedDocuments := formatCheckPromptDocuments(documents)
	prompt.WriteString("The documentation for this system is given in the <documents> tag, which contains a list of documents and the sections contained within each document.")
	prompt.WriteString("\n\n")
	prompt.WriteString(formattedDocuments)
	prompt.WriteString("\n\n")

	// Add Pull Request (if any)
	// Note: When we support more than just pull requests this will need to be updated
	if pr != nil {
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
	// Note: When we support more than just pull requests this will need to be updated
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

	// Add type-of-change specific information
	switch file.Action {
	case code.ActionInsert:
		// Add <file>
		prompt.WriteString("<file>\n")
		prompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.Filename))
		prompt.WriteString("  <file_content>\n")
		prompt.WriteString(string(file.Contents))
		prompt.WriteString("\n")
		prompt.WriteString("  </file_content>\n")
		prompt.WriteString("</file>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was created, ", file.Filename))
		prompt.WriteString("and that the contents of the created file are in <file>, ")
	case code.ActionModify:
		// Add <diff>
		prompt.WriteString("<diff>\n")
		prompt.WriteString(textDiff)
		prompt.WriteString("</diff>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was modified, ", file.Filename))
		prompt.WriteString("and that a patch representing the changes to that file is in <diff>, ")
	case code.ActionRename:
		// Add <diff> optionally
		if textDiff != "" {
			prompt.WriteString("<diff>\n")
			prompt.WriteString(textDiff)
			prompt.WriteString("</diff>\n")
			prompt.WriteString("\n")
		}
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was renamed to %s, ", file.OriginalFilename, file.Filename))
		if textDiff != "" {
			prompt.WriteString("and that a patch representing the changes to the renamed file is in <diff>, ")
		}
	case code.ActionDelete:
		// Add <file>
		prompt.WriteString("<file>\n")
		prompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.OriginalFilename))
		prompt.WriteString("  <file_content>\n")
		prompt.WriteString(string(file.OriginalContents))
		prompt.WriteString("\n")
		prompt.WriteString("  </file_content>\n")
		prompt.WriteString("</file>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was deleted, ", file.OriginalFilename))
		prompt.WriteString("and that the contents of the deleted file are in <file>, ")
	default:
		// Do nothing and return
		slog.Warn("check.Change encountered an unknown action", "file", file.Filename, "action", file.Action)
		return "", fmt.Errorf("unknown action encountered when creating prompt: %s", file.Action)
	}

	// Add prompt instructions for pull requests and/or issue(s)
	// Note: When we support more than just pull requests this will need to be updated
	if pr != nil {
		prompt.WriteString("and that the contents of related pull request(s) are in <pull_request>, ")
	}
	if numIssues > 0 {
		prompt.WriteString("and that the contents of related issue(s) are in <issue>, ")
	}

	// Add instructions
	prompt.WriteString("look at the documentation provided in <documents> and determine which documents, if any, should be updated based on this change.\n")
	prompt.WriteString(fmt.Sprintf("Then, call the provided %s tool with a list of ids of the documents and/or sections that should be updated along with the reason they should be updated.\n", checkNeedsUpdateName))
	prompt.WriteString(fmt.Sprintf("If there are no documents that need to be updated call the %s tool instead.", checkNoUpdateNeededName))

	return prompt.String(), nil
}

func formatCheckPromptDocuments(documents []*docs.FilteredDoc) string {
	var str strings.Builder
	indent := 0

	str.WriteString("<documents>\n")

	indent += 2

	for _, document := range documents {
		// <document>
		str.WriteString(fmt.Sprintf("%s<document id=\"%s\">\n", strings.Repeat(" ", indent), document.Document.ID))

		indent += 2

		// <document_name>{{NAME}}<document_name>
		str.WriteString(fmt.Sprintf("%s<document_name>%s</document_name>\n", strings.Repeat(" ", indent), document.Document.ID))

		// <document_purpose>{{PURPOSE}}</document_purpose>
		if document.Document.Purpose != "" {
			str.WriteString(fmt.Sprintf("%s<document_purpose>%s</document_purpose>\n", strings.Repeat(" ", indent), document.Document.Purpose))
		}

		// <sections>
		if len(document.Sections) > 0 {
			str.WriteString(formatCheckPromptSections(document.Sections, indent))
		}

		indent -= 2

		// </document>
		str.WriteString(fmt.Sprintf("%s</document>\n", strings.Repeat(" ", indent)))
	}

	indent -= 2

	str.WriteString("<documents>\n")

	return str.String()
}

func formatCheckPromptSections(sections []docs.FilteredSection, indent int) string {
	var str strings.Builder

	// <sections>
	str.WriteString(fmt.Sprintf("%s<sections>\n", strings.Repeat(" ", indent)))

	indent += 2

	for _, section := range sections {
		// <section id="">
		str.WriteString(fmt.Sprintf("%s<section id=\"%s\">\n", strings.Repeat(" ", indent), section.Section.ID))

		indent += 2

		// <section_name>{{NAME}}<section_name>
		str.WriteString(fmt.Sprintf("%s<section_name>%s</section_name>\n", strings.Repeat(" ", indent), section.Section.Name))

		// <section_purpose>{{PURPOSE}}</section_purpose>
		if section.Section.Purpose != "" {
			str.WriteString(fmt.Sprintf("%s<section_purpose>%s</section_purpose>\n", strings.Repeat(" ", indent), section.Section.Purpose))
		}

		indent -= 2

		// </section>
		str.WriteString(fmt.Sprintf("%s<section>\n", strings.Repeat(" ", indent)))
	}

	indent -= 2

	// </sections>
	str.WriteString(fmt.Sprintf("%s</sections>\n", strings.Repeat(" ", indent)))

	return str.String()
}

func getCheckTools(cb func(id string, reason string)) (tools []*llm.Tool) {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	tools = []*llm.Tool{
		{
			Name:        checkNeedsUpdateName,
			Description: "Identify a set of documents and/or sections that need to be updated for this change",
			Schema:      reflector.Reflect(&checkNeedsUpdateSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("check.Diff - determined updates are needed")
				// Parse the input
				var needsUpdate checkNeedsUpdateSchema
				err := json.Unmarshal([]byte(input), &needsUpdate)
				if err != nil {
					slog.Debug("check.Diff - could not parse tool call input, invalid json", "tool", checkNeedsUpdateName, "input", input, "error", err)
					return true, "", err
				}

				// Loop through and handle each document/section identified by the llm
				for _, update := range needsUpdate.Entries {
					cb(update.ID, update.Reason)
				}

				// Return with done = true so we stop
				return true, "", nil
			},
		},
		{
			Name:        checkNoUpdateNeededName,
			Description: "Identify that there are no documents that need to be updated for this change",
			Schema:      reflector.Reflect(&checkNoUpdateNeededSchema{}),
			Callback: func(params string) (bool, string, error) {
				slog.Debug("check.Diff - determined no updates needed")
				// Return with done = true so we stop
				return true, "", nil
			},
		},
	}

	return
}
