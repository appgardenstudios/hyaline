package check

import (
	"encoding/json"
	"hyaline/internal/code"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"testing"
)

func TestDiff_IgnoresInvalidIDsFromLLM(t *testing.T) {
	// 1. Mock llm.CallLLM
	originalCallLLM := callLLM
	defer func() { callLLM = originalCallLLM }()

	callLLM = func(systemPrompt string, prompt string, tools []*llm.Tool, cfg *config.LLM) (string, error) {
		// Simulate LLM calling the 'needs_update' tool with valid and invalid IDs
		for _, tool := range tools {
			if tool.Name == checkNeedsUpdateName {
				updates := checkNeedsUpdateSchema{
					Entries: []checkNeedsUpdateSchemaEntry{
						{ID: "document://docs/valid.md", Reason: "This is a valid doc"},
						{ID: "document://invalid/doc.md", Reason: "This is an invalid doc"},
						{ID: "document://docs/valid.md#valid-section", Reason: "This is a valid section"},
						{ID: "document://docs/valid.md#valid-section/nested-section", Reason: "This is a nested section"},
						{ID: "document://docs/valid.md#invalid-section", Reason: "This is an invalid section"},
					},
				}
				inputBytes, _ := json.Marshal(updates)
				tool.Callback(string(inputBytes))
				break
			}
		}
		return "", nil
	}

	// 2. Define valid documents and sections
	documents := []*docs.FilteredDoc{
		{
			Document: &sqlite.DOCUMENT{ID: "valid.md", SourceID: "docs"},
			Sections: []docs.FilteredSection{
				{
					Section: &sqlite.SECTION{ID: "valid-section", DocumentID: "valid.md", SourceID: "docs"},
					Sections: []docs.FilteredSection{
						{
							Section: &sqlite.SECTION{ID: "valid-section/nested-section", DocumentID: "valid.md", SourceID: "docs"},
						},
					},
				},
			},
		},
	}

	// 3. Define other inputs for Diff
	files := []code.FilteredFile{
		{Filename: "some/file.go", Action: code.ActionModify, Contents: []byte("hello")},
	}
	checkCfg := &config.Check{Options: config.CheckOptions{UpdateIf: config.CheckOptionsUpdateIf{}}}
	llmCfg := &config.LLM{}

	// 4. Call Diff
	results, _, err := Diff(files, documents, nil, nil, checkCfg, llmCfg)
	if err != nil {
		t.Fatalf("Diff returned an error: %v", err)
	}

	// 5. Assert results
	if len(results) != 3 {
		t.Errorf("Expected 3 results (doc, section, nested section), but got %d", len(results))
	}

	foundValidDoc := false
	foundValidSection := false
	foundNestedSection := false
	for _, result := range results {
		// The result struct splits the ID, so we check the components
		if result.Source == "docs" && result.Document == "valid.md" && len(result.Section) == 0 {
			foundValidDoc = true
		}
		if result.Source == "docs" && result.Document == "valid.md" && len(result.Section) == 1 && result.Section[0] == "valid-section" {
			foundValidSection = true
		}
		if result.Source == "docs" && result.Document == "valid.md" && len(result.Section) == 2 && result.Section[0] == "valid-section" && result.Section[1] == "nested-section" {
			foundNestedSection = true
		}
	}

	if !foundValidDoc {
		t.Errorf("Did not find expected result for docs/valid.md")
	}
	if !foundValidSection {
		t.Errorf("Did not find expected result for docs/valid.md#valid-section")
	}
	if !foundNestedSection {
		t.Errorf("Did not find expected result for docs/valid.md#valid-section/nested-section")
	}
}
