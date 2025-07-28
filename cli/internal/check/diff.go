package check

import (
	"encoding/json"
	"fmt"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/diff"
	"hyaline/internal/docs"
	"hyaline/internal/github"
	"hyaline/internal/llm"
	"log/slog"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
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

type updateResultMapCallback func(id string, reason string)

func Diff(files []code.FilteredFile, documents []*docs.FilteredDoc, pr *github.PullRequest, issues []*github.Issue, checkCfg *config.Check, llmCfg *config.LLM) (results []Result, err error) {
	resultMap := make(map[string][]string)

	updateResultMap := func(id string, reason string) {
		entry, ok := resultMap[id]
		if ok {
			entry = append(entry, reason)
			resultMap[id] = entry
		} else {
			resultMap[id] = []string{reason}
		}
	}

	// LLM system prompt and tools
	systemPrompt := "You are a senior technical writer who writes clear and accurate documentation."
	tools := getCheckTools(updateResultMap)

	// Check each file in the diff
	for _, file := range files {
		// See if there are any updateIfs that apply
		checkNewUpdateIfs(&file, documents, checkCfg, updateResultMap)

		// Ask LLM for documentation that should be updated for this diff
		var prompt string
		prompt, err = formatCheckPrompt(file, documents, pr, issues)
		if err != nil {
			slog.Debug("check.Diff could not format prompt", "error", err)
			return
		}
		slog.Debug("check.Diff calling llm", "file", file.Filename, "systemPrompt", systemPrompt, "prompt", prompt, "tools", len(tools))
		_, err = llm.CallLLM(systemPrompt, prompt, tools, llmCfg)
		if err != nil {
			slog.Debug("check.Change encountered an error when calling the llm", "error", err)
			return
		}
	}

	// Process resultMap into results
	for id, reasons := range resultMap {
		source, remainder, _ := strings.Cut(id, "/")
		document, sections, _ := strings.Cut(remainder, "#")
		section := []string{}
		if sections != "" {
			section = strings.Split(sections, "/")
		}
		results = append(results, Result{
			Source:   source,
			Document: document,
			Section:  section,
			Reasons:  reasons,
		})
	}

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
		str.WriteString(fmt.Sprintf("%s<document id=\"%s/%s\">\n", strings.Repeat(" ", indent), document.Document.SourceID, document.Document.ID))

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
		str.WriteString(fmt.Sprintf("%s<section id=\"%s/%s#%s\">\n", strings.Repeat(" ", indent), section.Section.SourceID, section.Section.DocumentID, section.Section.ID))

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

func getCheckTools(cb updateResultMapCallback) (tools []*llm.Tool) {
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

func checkNewUpdateIfs(file *code.FilteredFile, documents []*docs.FilteredDoc, cfg *config.Check, cb updateResultMapCallback) {
	// Check touched
	for _, entry := range cfg.Options.UpdateIf.Touched {
		if (file.Filename != "" && doublestar.MatchUnvalidated(entry.Code.Path, file.Filename)) ||
			(file.OriginalFilename != "" && doublestar.MatchUnvalidated(entry.Code.Path, file.OriginalFilename)) {
			if file.Action == code.ActionRename {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, fmt.Sprintf("touched (%s was renamed to %s)", file.OriginalFilename, file.Filename))
			} else {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, "touched")
			}
		}
	}

	// Check other updateIfs based on the action
	switch file.Action {
	case code.ActionInsert:
		for _, entry := range cfg.Options.UpdateIf.Added {
			if doublestar.MatchUnvalidated(entry.Code.Path, file.Filename) {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, "added")
			}
		}
	case code.ActionModify:
		for _, entry := range cfg.Options.UpdateIf.Modified {
			if doublestar.MatchUnvalidated(entry.Code.Path, file.Filename) {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, "modified")
			}
		}
	case code.ActionRename:
		for _, entry := range cfg.Options.UpdateIf.Renamed {
			if doublestar.MatchUnvalidated(entry.Code.Path, file.Filename) ||
				doublestar.MatchUnvalidated(entry.Code.Path, file.OriginalFilename) {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, "renamed")
			}
		}
	case code.ActionDelete:
		for _, entry := range cfg.Options.UpdateIf.Deleted {
			if doublestar.MatchUnvalidated(entry.Code.Path, file.OriginalFilename) {
				checkNewUpdateIfDocuments(entry.Code.Path, documents, entry.Documentation, cb, "deleted")
			}
		}
	}
}

func checkNewUpdateIfDocuments(glob string, documents []*docs.FilteredDoc, filter config.CheckDocumentationFilter, cb updateResultMapCallback, action string) {
	for _, document := range documents {
		if docs.DocumentMatches(document.Document.ID, document.Document.SourceID, document.Tags, &filter) {
			cb(document.Document.SourceID+"/"+document.Document.ID,
				fmt.Sprintf("Update this document if any files matching %s were %s.", glob, action))
		} else {
			// Only check sections if document does not match to avoid pulling in a document
			// and all of its sections (we just need the document in that case)
			checkNewUpdateIfSections(glob, document.Sections, filter, cb, action)
		}
	}
}

func checkNewUpdateIfSections(glob string, sections []docs.FilteredSection, filter config.CheckDocumentationFilter, cb updateResultMapCallback, action string) {
	for _, section := range sections {
		if docs.SectionMatches(section.Section.ID, section.Section.DocumentID, section.Section.SourceID, section.Tags, &filter) {
			cb(section.Section.SourceID+"/"+section.Section.DocumentID+"#"+section.Section.ID,
				fmt.Sprintf("Update this document if any files matching %s were %s.", glob, action))
		}
		checkNewUpdateIfSections(glob, section.Sections, filter, cb, action)
	}
}
