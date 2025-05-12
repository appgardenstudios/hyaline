package check

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/diff"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/invopop/jsonschema"
)

type ChangeResult struct {
	DocumentationSource string
	Document            string
	Section             []string
	Reasons             []string
	References          []ChangeResultReference
}

type ChangeResultReference struct {
	CodeID string
	FileID string
	Diff   string
}

func Change(file *sqlite.File, codeSource config.CodeSource, desiredDocsMap map[string][]config.Document, pullRequests []*sqlite.PullRequest, issues []*sqlite.Issue, currentDB *sql.DB, changeDB *sql.DB, cfg *config.LLM) (results []ChangeResult, err error) {
	slog.Debug("check.Change checking file", "file", file.ID)

	// Get original ID and contents so we can calculate a diff
	originalID := file.ID
	originalContents := ""
	if file.Action == sqlite.ActionModify || file.Action == sqlite.ActionDelete {
		// Get original contents
		var original *sqlite.File
		original, err = sqlite.GetFile(file.ID, file.CodeID, file.SystemID, currentDB)
		if err != nil || original == nil {
			slog.Debug("check.Change could not get original file from modification", "file", file.ID, "error", err)
			return
		}
		originalContents = original.RawData
	}
	if file.Action == sqlite.ActionRename {
		// Get original contents
		var original *sqlite.File
		original, err = sqlite.GetFile(file.OriginalID, file.CodeID, file.SystemID, currentDB)
		if err != nil || original == nil {
			slog.Debug("check.Change could not get original file from rename", "file", file.ID, "error", err)
			return
		}
		originalID = original.ID
		originalContents = original.RawData
	}

	// Calculate the diff
	edits := diff.Strings(originalContents, file.RawData)
	textDiff, err := diff.ToUnified("a/"+originalID, "b/"+file.ID, originalContents, edits, 3)
	if err != nil {
		slog.Debug("check.Change could not generate diff", "file", file.ID, "error", err)
		return
	}

	// Check LLM
	llmResults, err := checkLLM(file, codeSource, textDiff, desiredDocsMap, pullRequests, issues, currentDB, changeDB, cfg)
	if err != nil {
		return
	}
	results = append(results, llmResults...)

	// Check updateIfs
	results = append(results, checkUpdateIfs(file.ID, file.OriginalID, file.Action, desiredDocsMap)...)

	// Loop through and add reference
	for idx := range results {
		results[idx].References = []ChangeResultReference{{
			CodeID: file.CodeID,
			FileID: file.ID,
			Diff:   textDiff,
		}}
	}

	return
}

const checkLLMNeedsUpdateName = "needs_update"
const checkLLMNoUpdateNeededName = "no_update_needed"

type checkLLMNeedsUpdateSchema struct {
	Entries []checkLLMNeedsUpdateSchemaEntry `json:"entries" jsonschema:"title=The list of entries,description=The list of documents and/or sections that need to be updated along with the reason for each update"`
}

type checkLLMNeedsUpdateSchemaEntry struct {
	ID     string `json:"id" jsonschema:"title=The document/section ID,description=The ID of the document and/or section that needs to be updated,example=app.1"`
	Reason string `json:"reason" jsonschema:"title=The reason,description=The reason the document and/or section needs to be updated,example=This section needs to be updated because the change modifies a file that is mentioned in the reference to this section"`
}

type checkLLMNoUpdateNeededSchema struct {
}

func checkLLM(file *sqlite.File, codeSource config.CodeSource, textDiff string, desiredDocsMap map[string][]config.Document, pullRequests []*sqlite.PullRequest, issues []*sqlite.Issue, currentDB *sql.DB, changeDB *sql.DB, cfg *config.LLM) (results []ChangeResult, err error) {
	slog.Debug("check.checkLLM checking file", "file", file.ID)

	// Generate the system and user prompt
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	// Add documentation context
	// https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/long-context-tips#example-quote-extraction
	documents, documentMap := formatDocuments(desiredDocsMap)
	prompt.WriteString("The documentation for this system is given in the <documents> tag, which contains a list of documents and the sections contained within each document.")
	prompt.WriteString("\n\n")
	prompt.WriteString(documents)
	prompt.WriteString("\n\n")

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

	// Add type-of-change specific information
	switch file.Action {
	case sqlite.ActionInsert:
		// Add <file>
		prompt.WriteString("<file>\n")
		prompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.ID))
		prompt.WriteString("  <file_content>\n")
		prompt.WriteString(file.RawData)
		prompt.WriteString("\n")
		prompt.WriteString("  </file_content>\n")
		prompt.WriteString("</file>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was created, ", file.ID))
		prompt.WriteString("and that the contents of the created file are in <file>, ")
	case sqlite.ActionModify:
		// Add <diff>
		prompt.WriteString("<diff>\n")
		prompt.WriteString(textDiff)
		prompt.WriteString("</diff>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was modified, ", file.ID))
		prompt.WriteString("and that a patch representing the changes to that file is in <diff>, ")
	case sqlite.ActionRename:
		// Add <diff> optionally
		if textDiff != "" {
			prompt.WriteString("<diff>\n")
			prompt.WriteString(textDiff)
			prompt.WriteString("</diff>\n")
			prompt.WriteString("\n")
		}
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was renamed to %s, ", file.OriginalID, file.ID))
		if textDiff != "" {
			prompt.WriteString("and that a patch representing the changes to the renamed file is in <diff>, ")
		}
	case sqlite.ActionDelete:
		// Add <file>
		prompt.WriteString("<file>\n")
		prompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.ID))
		prompt.WriteString("  <file_content>\n")
		prompt.WriteString(file.RawData)
		prompt.WriteString("\n")
		prompt.WriteString("  </file_content>\n")
		prompt.WriteString("</file>\n")
		prompt.WriteString("\n")
		// Add prompt
		prompt.WriteString(fmt.Sprintf("Given that the file %s was deleted, ", file.ID))
		prompt.WriteString("and that the contents of the deleted file are in <file>, ")
	default:
		// Do nothing and return
		slog.Warn("check.Change encountered an unknown action", "file", file.ID, "action", file.Action)
		return
	}

	// Add prompt instructions for pull requests and/or issue(s)
	if numPullRequests > 0 {
		prompt.WriteString("and that the contents of related pull request(s) are in <pull_request>, ")
	}
	if numIssues > 0 {
		prompt.WriteString("and that the contents of related issue(s) are in <issue>, ")
	}

	// Add instructions
	prompt.WriteString("look at the documentation provided in <documents> and determine which documents, if any, should be updated based on this change.\n")
	prompt.WriteString(fmt.Sprintf("Then, call the provided %s tool with a list of ids of the documents and/or sections that should be updated along with the reason they should be updated.\n", checkLLMNeedsUpdateName))
	prompt.WriteString(fmt.Sprintf("If there are no documents that need to be updated call the %s tool instead.", checkLLMNoUpdateNeededName))

	// Create tool(s)
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	tools := []*llm.Tool{
		{
			Name:        checkLLMNeedsUpdateName,
			Description: "Identify a set of documents and/or sections that need to be updated for this change",
			Schema:      reflector.Reflect(&checkLLMNeedsUpdateSchema{}),
			Callback: func(input string) (bool, string, error) {
				slog.Debug("check.Change - checkLLM determined updates are needed")
				// Parse the input
				var needsUpdate checkLLMNeedsUpdateSchema
				err := json.Unmarshal([]byte(input), &needsUpdate)
				if err != nil {
					slog.Debug("check.Change - checkLLM could not parse tool call input, invalid json", "tool", checkLLMNeedsUpdateName, "input", input, "error", err)
					return true, "", err
				}

				// Loop through and handle each document/section identified by the llm
				for _, update := range needsUpdate.Entries {
					mapEntry, ok := documentMap[update.ID]
					if !ok {
						slog.Debug("check.Change - checkLLM could not find referenced documentation ID", "tool", checkLLMNeedsUpdateName, "ID", update.ID, "error", err)
						continue
					}

					// Add a result for this entry
					results = append(results, ChangeResult{
						DocumentationSource: mapEntry.DocumentationSource,
						Document:            mapEntry.Document,
						Section:             mapEntry.Section,
						Reasons:             []string{update.Reason},
					})
				}

				// Return with done = true so we stop
				return true, "", nil
			},
		},
		{
			Name:        checkLLMNoUpdateNeededName,
			Description: "Identify that there are no documents that need to be updated for this change",
			Schema:      reflector.Reflect(&checkLLMNoUpdateNeededSchema{}),
			Callback: func(params string) (bool, string, error) {
				slog.Debug("check.Change - checkLLM determined no updates needed")
				// Return with done = true so we stop
				return true, "", nil
			},
		},
	}

	// Call LLM
	userPrompt := prompt.String()
	// fmt.Println(userPrompt)
	slog.Debug("check.Change calling the llm")
	// slog.Debug("check.Change calling the llm", "systemPrompt", systemPrompt, "userPrompt", userPrompt)
	_, err = llm.CallLLM(systemPrompt, userPrompt, tools, cfg)
	if err != nil {
		slog.Debug("check.Change encountered an error when calling the llm", "error", err)
		return
	}

	return
}

type documentMapEntry struct {
	DocumentationSource string
	Document            string
	Section             []string
}

func formatDocuments(desiredDocsMap map[string][]config.Document) (string, map[string]documentMapEntry) {
	var documents strings.Builder
	documentMap := make(map[string]documentMapEntry)
	indent := 0

	documents.WriteString("<documents>\n")

	indent += 2

	for docID, desiredDocs := range desiredDocsMap {
		for idx, desiredDoc := range desiredDocs {
			id := fmt.Sprintf("%s.%d", docID, idx+1)
			documentMap[id] = documentMapEntry{
				DocumentationSource: docID,
				Document:            desiredDoc.Path,
				Section:             []string{},
			}

			// <document>
			documents.WriteString(fmt.Sprintf("%s<document id=\"%s\">\n", strings.Repeat(" ", indent), id))
			indent += 2

			// <document_name>{{NAME}}<document_name>
			documents.WriteString(fmt.Sprintf("%s<document_name>%s</document_name>\n", strings.Repeat(" ", indent), desiredDoc.Path))

			// <document_purpose>{{PURPOSE}}</document_purpose>
			if desiredDoc.Purpose != "" {
				documents.WriteString(fmt.Sprintf("%s<document_purpose>%s</document_purpose>\n", strings.Repeat(" ", indent), desiredDoc.Purpose))
			}

			// <sections>
			if len(desiredDoc.Sections) > 0 {
				documents.WriteString(formatSections(desiredDoc.Sections, id, []string{}, indent, docID, desiredDoc.Path, &documentMap))
			}

			indent -= 2

			// </document>
			documents.WriteString(fmt.Sprintf("%s</document>\n", strings.Repeat(" ", indent)))
		}
	}

	indent -= 2

	documents.WriteString("<documents>\n")

	return documents.String(), documentMap
}

// Note: only call this if len(sections) > 0
func formatSections(sections []config.DocumentSection, prefix string, parents []string, indent int, documentSource string, document string, documentMap *map[string]documentMapEntry) string {
	var str strings.Builder

	// <sections>
	str.WriteString(fmt.Sprintf("%s<sections>\n", strings.Repeat(" ", indent)))

	indent += 2

	for idx, section := range sections {
		id := fmt.Sprintf("%s.%d", prefix, idx+1)
		sectionArr := []string{}
		sectionArr = append(sectionArr, parents...)
		sectionArr = append(sectionArr, section.ID)
		(*documentMap)[id] = documentMapEntry{
			DocumentationSource: documentSource,
			Document:            document,
			Section:             sectionArr,
		}

		// <section id="">
		str.WriteString(fmt.Sprintf("%s<section id=\"%s\">\n", strings.Repeat(" ", indent), id))

		indent += 2

		// <section_name>{{NAME}}<section_name>
		str.WriteString(fmt.Sprintf("%s<section_name>%s</section_name>\n", strings.Repeat(" ", indent), section.ID))

		// <section_purpose>{{PURPOSE}}</section_purpose>
		if section.Purpose != "" {
			str.WriteString(fmt.Sprintf("%s<section_purpose>%s</section_purpose>\n", strings.Repeat(" ", indent), section.Purpose))
		}

		// <sections> if present
		if len(section.Sections) > 0 {
			str.WriteString(formatSections(section.Sections, id, sectionArr, indent, documentSource, document, documentMap))
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

func checkUpdateIfs(id string, originalID string, action sqlite.Action, ruleDocsMap map[string][]config.Document) (results []ChangeResult) {
	for docSource, ruleDocs := range ruleDocsMap {
		for _, ruleDoc := range ruleDocs {
			// Check touched
			for _, glob := range ruleDoc.UpdateIf.Touched {
				if doublestar.MatchUnvalidated(glob, id) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            ruleDoc.Path,
						Section:             []string{},
						Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were touched", glob)},
					})
				} else if originalID != "" && doublestar.MatchUnvalidated(glob, originalID) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            ruleDoc.Path,
						Section:             []string{},
						Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were touched (%s was renamed to %s)", glob, originalID, id)},
					})
				}
			}
			// Check action specific updates
			switch action {
			case sqlite.ActionInsert:
				// Check added
				for _, glob := range ruleDoc.UpdateIf.Added {
					if doublestar.MatchUnvalidated(glob, id) {
						results = append(results, ChangeResult{
							DocumentationSource: docSource,
							Document:            ruleDoc.Path,
							Section:             []string{},
							Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were added", glob)},
						})
					}
				}
			case sqlite.ActionModify:
				// Check modified
				for _, glob := range ruleDoc.UpdateIf.Modified {
					if doublestar.MatchUnvalidated(glob, id) {
						results = append(results, ChangeResult{
							DocumentationSource: docSource,
							Document:            ruleDoc.Path,
							Section:             []string{},
							Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were modified", glob)},
						})
					}
				}
			case sqlite.ActionRename:
				// Check renamed
				for _, glob := range ruleDoc.UpdateIf.Renamed {
					if doublestar.MatchUnvalidated(glob, id) || doublestar.MatchUnvalidated(glob, originalID) {
						results = append(results, ChangeResult{
							DocumentationSource: docSource,
							Document:            ruleDoc.Path,
							Section:             []string{},
							Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were renamed (%s was renamed to %s)", glob, id, originalID)},
						})
					}
				}
			case sqlite.ActionDelete:
				// Check deleted
				for _, glob := range ruleDoc.UpdateIf.Deleted {
					if doublestar.MatchUnvalidated(glob, id) {
						results = append(results, ChangeResult{
							DocumentationSource: docSource,
							Document:            ruleDoc.Path,
							Section:             []string{},
							Reasons:             []string{fmt.Sprintf("Update this document if any files matching %s were deleted", glob)},
						})
					}
				}
			}

			// Check sections
			results = append(results, checkSectionUpdateIfs(id, originalID, action, docSource, ruleDoc.Path, ruleDoc.Sections, []string{})...)
		}
	}

	return
}

func checkSectionUpdateIfs(id string, originalID string, action sqlite.Action, docSource string, document string, sections []config.DocumentSection, parents []string) (results []ChangeResult) {
	for _, section := range sections {
		sectionArr := []string{}
		sectionArr = append(sectionArr, parents...)
		sectionArr = append(sectionArr, section.ID)

		// Check touched
		for _, glob := range section.UpdateIf.Touched {
			if doublestar.MatchUnvalidated(glob, id) {
				results = append(results, ChangeResult{
					DocumentationSource: docSource,
					Document:            document,
					Section:             sectionArr,
					Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were touched", glob)},
				})
			} else if originalID != "" && doublestar.MatchUnvalidated(glob, originalID) {
				results = append(results, ChangeResult{
					DocumentationSource: docSource,
					Document:            document,
					Section:             sectionArr,
					Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were touched (%s was renamed to %s)", glob, originalID, id)},
				})
			}
		}
		// Check action specific updates
		switch action {
		case sqlite.ActionInsert:
			// Check added
			for _, glob := range section.UpdateIf.Added {
				if doublestar.MatchUnvalidated(glob, id) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            document,
						Section:             sectionArr,
						Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were added", glob)},
					})
				}
			}
		case sqlite.ActionModify:
			// Check modified
			for _, glob := range section.UpdateIf.Modified {
				if doublestar.MatchUnvalidated(glob, id) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            document,
						Section:             sectionArr,
						Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were modified", glob)},
					})
				}
			}
		case sqlite.ActionRename:
			// Check renamed
			for _, glob := range section.UpdateIf.Renamed {
				if doublestar.MatchUnvalidated(glob, id) || doublestar.MatchUnvalidated(glob, originalID) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            document,
						Section:             sectionArr,
						Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were renamed (%s was renamed to %s)", glob, id, originalID)},
					})
				}
			}
		case sqlite.ActionDelete:
			// Check deleted
			for _, glob := range section.UpdateIf.Deleted {
				if doublestar.MatchUnvalidated(glob, id) {
					results = append(results, ChangeResult{
						DocumentationSource: docSource,
						Document:            document,
						Section:             sectionArr,
						Reasons:             []string{fmt.Sprintf("Update this section if any files matching %s were deleted", glob)},
					})
				}
			}
		}

		// Check sections
		results = append(results, checkSectionUpdateIfs(id, originalID, action, docSource, document, section.Sections, sectionArr)...)
	}

	return
}
