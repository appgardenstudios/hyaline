package check

import (
	"database/sql"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/diff"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

type ChangeResult struct {
	DocumentationSource string
	Document            string
	Section             string
	Reasons             []string
}

func Change(file *sqlite.File, codeSource config.CodeSource, ruleDocsMap map[string][]config.RuleDocument, currentDB *sql.DB) (results []ChangeResult, err error) {
	originalID := file.ID
	originalContents := ""
	if file.Action == sqlite.ActionModify {
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
	edits := diff.Strings(originalContents, file.RawData)
	textDiff, err := diff.ToUnified("a/"+originalID, "b/"+file.ID, originalContents, edits, 3)
	if err != nil {
		slog.Debug("check.Change could not generate diff", "file", file.ID, "error", err)
		return
	}

	// Ignore white space only changes?
	// TODO

	// Generate the user prompt
	var userPrompt strings.Builder

	// TODO give context
	userPrompt.WriteString("The documentation for this system is given in the <documents> tag, which contains a list of documents and the sections contained within the document.") // TODO finish this description
	userPrompt.WriteString("\n\n")

	// https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/long-context-tips#example-quote-extraction

	// Build document structure for prompt
	userPrompt.WriteString(formatDocuments(ruleDocsMap))
	userPrompt.WriteString("\n\n")

	switch file.Action {
	case sqlite.ActionInsert:
		// Add <file>
		userPrompt.WriteString("<file>\n")
		userPrompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.ID))
		userPrompt.WriteString("  <file_content>\n")
		userPrompt.WriteString(file.RawData)
		userPrompt.WriteString("\n")
		userPrompt.WriteString("  </file_content>\n")
		userPrompt.WriteString("</file>\n")
		userPrompt.WriteString("\n")
		// Add prompt
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was created, ", file.ID))
		userPrompt.WriteString("and that the contents of the created file are in <file>, ")
	case sqlite.ActionModify:
		// Add <diff>
		userPrompt.WriteString("<diff>\n")
		userPrompt.WriteString(textDiff)
		userPrompt.WriteString("\n")
		userPrompt.WriteString("</diff>\n")
		userPrompt.WriteString("\n")
		// Add prompt
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was modified, ", file.ID))
		userPrompt.WriteString("and that a patch representing the changes to that file is in <diff>, ")
	case sqlite.ActionRename:
		// Add <diff> optionally
		if textDiff != "" {
			userPrompt.WriteString("<diff>\n")
			userPrompt.WriteString(textDiff)
			userPrompt.WriteString("\n")
			userPrompt.WriteString("</diff>\n")
			userPrompt.WriteString("\n")
		}
		// Add prompt
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was renamed to %s, ", file.OriginalID, file.ID))
		if textDiff != "" {
			userPrompt.WriteString("and that a patch representing the changes to the renamed file is in <diff>, ")
		}
	case sqlite.ActionDelete:
		// Add <file>
		userPrompt.WriteString("<file>\n")
		userPrompt.WriteString(fmt.Sprintf("  <file_name>%s</file_name>\n", file.ID))
		userPrompt.WriteString("  <file_content>\n")
		userPrompt.WriteString(file.RawData)
		userPrompt.WriteString("\n")
		userPrompt.WriteString("  </file_content>\n")
		userPrompt.WriteString("</file>\n")
		userPrompt.WriteString("\n")
		// Add prompt
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was deleted, ", file.ID))
		userPrompt.WriteString("and that the contents of the deleted file are in <file>, ")
	default:
		// Do nothing and return
		// TODO log this
		return
	}
	userPrompt.WriteString("look at the documentation provided in <documents> and determine which documents, if any, should be updated based on this change.") // TODO call tool to update?
	userPrompt.WriteString("\n\n")

	fmt.Println(userPrompt.String())

	for docSource := range ruleDocsMap {
		results = append(results, ChangeResult{
			DocumentationSource: docSource,
			Document:            "README.md",
			Section:             "",
			Reasons:             []string{fmt.Sprintf("testReason for file %s in %s", file.ID, codeSource.ID)},
		})
		break
	}

	// TODO respect updateIfs

	// Prompt: given the set of system documentation in <documentation> and the change in <change>, what documentation should be updated? respond with a tool call to update_documentation(list)
	return
}

func formatDocuments(ruleDocsMap map[string][]config.RuleDocument) string {
	var documents strings.Builder
	indent := 0

	documents.WriteString("<documents>\n")

	indent += 2

	for docID, ruleDocs := range ruleDocsMap {
		for idx, ruleDoc := range ruleDocs {
			id := fmt.Sprintf("%s.%d", docID, idx+1)

			// <document>
			documents.WriteString(fmt.Sprintf("%s<document uid=\"%s\">\n", strings.Repeat(" ", indent), id))
			indent += 2

			// <document_name>{{NAME}}<document_name>
			documents.WriteString(fmt.Sprintf("%s<document_name>%s</document_name>\n", strings.Repeat(" ", indent), ruleDoc.Path))

			// <document_purpose>{{PURPOSE}}</document_purpose>
			if ruleDoc.Purpose != "" {
				documents.WriteString(fmt.Sprintf("%s<document_purpose>%s</document_purpose>\n", strings.Repeat(" ", indent), ruleDoc.Purpose))
			}

			// <sections>
			if len(ruleDoc.Sections) > 0 {
				documents.WriteString(formatSections(ruleDoc.Sections, id, indent))
			}

			indent -= 2

			// </document>
			documents.WriteString(fmt.Sprintf("%s</document>\n", strings.Repeat(" ", indent)))
		}
	}

	indent -= 2

	documents.WriteString("<documents>\n")

	return documents.String()
}

// Note: only call this if len(sections) > 0
// TODO add assert?
func formatSections(sections []config.RuleDocumentSection, prefix string, indent int) string {
	var str strings.Builder

	// <sections>
	str.WriteString(fmt.Sprintf("%s<sections>\n", strings.Repeat(" ", indent)))

	indent += 2

	for idx, section := range sections {
		id := fmt.Sprintf("%s.%d", prefix, idx+1)
		// <section id="">
		str.WriteString(fmt.Sprintf("%s<section uid=\"%s\">\n", strings.Repeat(" ", indent), id))

		indent += 2

		// <section_name>{{NAME}}<section_name>
		str.WriteString(fmt.Sprintf("%s<section_name>%s</section_name>\n", strings.Repeat(" ", indent), section.ID))

		// <section_purpose>{{PURPOSE}}</section_purpose>
		if section.Purpose != "" {
			str.WriteString(fmt.Sprintf("%s<section_purpose>%s</section_purpose>\n", strings.Repeat(" ", indent), section.Purpose))
		}

		// <sections> if present
		if len(section.Sections) > 0 {
			str.WriteString(formatSections(section.Sections, id, indent))
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
