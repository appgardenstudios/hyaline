package suggest

import (
	"fmt"
	"strings"
)

type SuggestionData struct {
	Diffs []SuggestionDataDiff
}

type SuggestionDataDiff struct {
	// TODO
}

// Eventually we should group this up and handle a document and section updates in the same call
func Change(systemID string, documentationSource string, document string, section []string) (string, error) {
	systemPrompt := "You are a senior technical writer who writes clear and accurate system documentation."
	var prompt strings.Builder

	if len(section) > 0 {

	} else {

	}

	// update document based off of diffs
	// vs
	// update document and section(s) based off of diffs

	// Prompt: Based on the changes contained in <diffs>, and the prs/issues, what changes should be made to document/section contained in <doc/sec>. Call tools with change, or toolX for no change. supply the entire document, and do NOT put diffs. stay concise and match the existing style of the document, ...

	fmt.Println(systemPrompt)
	fmt.Println(prompt.String())

	return "", nil
}
