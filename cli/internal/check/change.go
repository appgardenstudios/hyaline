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
	userPrompt.WriteString("The documentation for this system is given in the <documents> tag, which contains a list of documents and the sections contained within the document") // TODO finish this description
	// TODO put this after the documents?

	// https://docs.anthropic.com/en/docs/build-with-claude/prompt-engineering/long-context-tips#example-quote-extraction

	// Build document structure for prompt
	// TODO
	_ = `
<documents>
	<document id="">
	  <document_name>{{NAME}}<document_name>
		<document_purpose>{{PURPOSE}}</document_purpose>
		<document_sections>
		  <section id="">
			  <section_name>{{NAME}}<section_name>
			  <section_purpose>{{PURPOSE}}</section_purpose>
			  <section_content>
			  {{CONTENTS}}
			  </section_content>
				<sections>...</sections>
			</section>
		</document_sections>
	</document>
<documents>
`

	switch file.Action {
	case sqlite.ActionInsert:
		// TODO add <file>
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was created, ", file.ID))
		userPrompt.WriteString("and that the contents of the created file are in <file>, ")
	case sqlite.ActionModify:
		// TODO add <diff>
		userPrompt.WriteString(fmt.Sprintf("Given that the file %s was modified, ", file.ID))
		userPrompt.WriteString("and that a patch representing the changes to that file are in <diff>, ")
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
