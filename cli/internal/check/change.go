package check

import (
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
)

// https://pkg.go.dev/github.com/sergi/go-diff#section-readme

func Change(file *sqlite.File, ruleDocs []config.RuleDocument) {
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

	// Prompt: given the set of system documentation in <documentation> and the change in <change>, what documentation should be updated? respond with a tool call to update_documentation(list)
}
