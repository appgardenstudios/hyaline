package check

import (
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"strings"
)

type ChangeResult struct {
	DocumentationSource string
	Document            string
	Section             string
	Reasons             []string
}

// https://pkg.go.dev/github.com/sergi/go-diff#section-readme

func Change(file *sqlite.File, codeSource config.CodeSource, ruleDocsMap map[string][]config.RuleDocument) (results []ChangeResult, err error) {
	// Calculate the diff and ignore whitespace only changes
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
		// TODO add <file>
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was created, ", file.ID))
		userPrompt.WriteString("and that the contents of the created file are in <file>, ")
	case sqlite.ActionModify:
		// TODO add <diff>
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was modified, ", file.ID))
		userPrompt.WriteString("and that a patch representing the changes to that file is in <diff>, ")
	case sqlite.ActionRename:
		// TODO add <diff> optionally
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was renamed to %s, ", file.OriginalID, file.ID))
		// TODO handle a diff + rename by adding a reference to the patch
	case sqlite.ActionDelete:
		// TODO add <file>
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was deleted, ", file.ID))
		userPrompt.WriteString("and that the contents of the deleted file are in <file>, ")
	default:
		// Do nothing and return
		// TODO log this
		return
	}
	userPrompt.WriteString("look at the documentation sources provided in <sources> and determine which documents, if any, should be updated based on this change.")
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
