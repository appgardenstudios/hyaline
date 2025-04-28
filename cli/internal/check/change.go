package check

import (
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
)

// https://pkg.go.dev/github.com/sergi/go-diff#section-readme

type ChangeResult struct {
	DocumentationSource string
	Document            string
	Section             string
	Reasons             []string
}

func Change(file *sqlite.File, codeSource config.CodeSource, ruleDocsMap map[string][]config.RuleDocument) (results []ChangeResult, err error) {
	// Calculate the diff and ignore whitespace only changes
	switch file.Action {
	case sqlite.ActionInsert:
		// TODO
	case sqlite.ActionModify:
		// TODO
	case sqlite.ActionRename:
		// TODO
	case sqlite.ActionDelete:
		// TODO
	default:
		// Do nothing and return
		// TODO
	}

	for docSource, _ := range ruleDocsMap {
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
