package check

import (
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
)

// https://pkg.go.dev/github.com/sergi/go-diff#section-readme

type ChangeResult struct {
	Documentation string
	Document      string
	Section       string
	Reasons       []string
}

func Change(file *sqlite.File, ruleDocsMap map[string][]config.RuleDocument) (results []ChangeResult, err error) {
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

	results = append(results, ChangeResult{
		Documentation: "testDocumentation",
		Document:      "testDocument",
		Section:       "testSection",
		Reasons:       []string{"testReason"},
	})

	// TODO respect updateIfs

	// Prompt: given the set of system documentation in <documentation> and the change in <change>, what documentation should be updated? respond with a tool call to update_documentation(list)
	return
}
