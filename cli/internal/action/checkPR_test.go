package action

import (
	"hyaline/internal/check"
	"strings"
	"testing"
)

func createTestReason(reason string, checkType check.DiffCheckType, file string, contextHash string, outdated bool) check.Reason {
	return check.Reason{
		Reason:   reason,
		Outdated: outdated,
		Check: check.DiffCheck{
			Type:        checkType,
			File:        file,
			ContextHash: contextHash,
		},
	}
}

func TestMergeCheckRecommendations_NewOnly(t *testing.T) {
	newRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("New file added", check.DiffCheckTypeLLM, "main.go", "hash1", false)},
	}

	contextHashes := check.FileCheckContextHashes{
		"main.go": {check.DiffCheckTypeLLM: "hash1"},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{newRec}, []CheckRecommendation{}, contextHashes)

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if result[0].Outdated {
		t.Errorf("Expected recommendation to not be outdated")
	}
}

func TestMergeCheckRecommendations_ExistingOnly(t *testing.T) {
	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Existing reason", check.DiffCheckTypeLLM, "main.go", "hash1", false)},
		Checked:  true,
	}

	contextHashes := check.FileCheckContextHashes{
		"main.go": {check.DiffCheckTypeLLM: "hash1"},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{}, []CheckRecommendation{existingRec}, contextHashes)

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if result[0].Outdated {
		t.Errorf("Expected recommendation to not be outdated")
	}
	if !result[0].Checked {
		t.Errorf("Expected recommendation to remain checked")
	}
}

func TestMergeCheckRecommendations_ExistingOutdated(t *testing.T) {
	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Old reason", check.DiffCheckTypeLLM, "main.go", "old-hash", false)},
		Checked:  true,
	}

	contextHashes := check.FileCheckContextHashes{
		"main.go": {check.DiffCheckTypeLLM: "new-hash"},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{}, []CheckRecommendation{existingRec}, contextHashes)

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if !result[0].Outdated {
		t.Errorf("Expected recommendation to be outdated")
	}
	if !result[0].Reasons[0].Outdated {
		t.Errorf("Expected reason to be marked as outdated")
	}
}

func TestMergeCheckRecommendations_MergeMatching(t *testing.T) {
	newRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("New LLM reason", check.DiffCheckTypeLLM, "main.go", "hash2", false)},
	}

	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Existing touched reason", check.DiffCheckTypeUpdateIfTouched, "main.go", "hash1", false)},
		Checked:  true,
	}

	contextHashes := check.FileCheckContextHashes{
		"main.go": {
			check.DiffCheckTypeLLM:             "hash2",
			check.DiffCheckTypeUpdateIfTouched: "hash3",
		},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{newRec}, []CheckRecommendation{existingRec}, contextHashes)

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if len(result[0].Reasons) != 2 {
		t.Errorf("Expected 2 reasons, got %d", len(result[0].Reasons))
	}
	if !result[0].Checked {
		t.Errorf("Expected to preserve checked state")
	}

	// Non-outdated reasons sort first
	if result[0].Reasons[0].Outdated {
		t.Errorf("Expected first reason to be current")
	}
	if !result[0].Reasons[1].Outdated {
		t.Errorf("Expected second reason to be outdated")
	}
	if result[0].Outdated {
		t.Errorf("Expected recommendation to not be outdated")
	}
}

func TestMergeCheckRecommendations_DifferentDocuments(t *testing.T) {
	newRec := CheckRecommendation{
		Source:   "docs",
		Document: "API.md",
		Reasons:  []check.Reason{createTestReason("New API change", check.DiffCheckTypeLLM, "api.go", "hash1", false)},
	}

	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Existing readme update", check.DiffCheckTypeLLM, "main.go", "hash2", false)},
		Checked:  true,
	}

	contextHashes := check.FileCheckContextHashes{
		"api.go":  {check.DiffCheckTypeLLM: "hash1"},
		"main.go": {check.DiffCheckTypeLLM: "hash2"},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{newRec}, []CheckRecommendation{existingRec}, contextHashes)

	if len(result) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(result))
	}

	// Sorted alphabetically by document: API.md, README.md
	if result[0].Document != "API.md" {
		t.Errorf("Expected first document to be API.md, got %s", result[0].Document)
	}
	if result[0].Outdated {
		t.Errorf("Expected API.md to not be outdated")
	}
	if result[1].Document != "README.md" {
		t.Errorf("Expected second document to be README.md, got %s", result[1].Document)
	}
	if result[1].Outdated {
		t.Errorf("Expected README.md to not be outdated")
	}
}

func TestMergeCheckRecommendations_EmptyContextHashes(t *testing.T) {
	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Existing reason", check.DiffCheckTypeLLM, "main.go", "some-hash", false)},
		Checked:  true,
	}

	result := mergeCheckRecommendations([]CheckRecommendation{}, []CheckRecommendation{existingRec}, make(check.FileCheckContextHashes))

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if !result[0].Outdated {
		t.Errorf("Expected recommendation to be outdated")
	}
	if !result[0].Reasons[0].Outdated {
		t.Errorf("Expected reason to be outdated")
	}
}

func TestMergeCheckRecommendations_MatchingReasons(t *testing.T) {
	newRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("New reason", check.DiffCheckTypeLLM, "main.go", "new-hash", false)},
	}

	existingRec := CheckRecommendation{
		Source:   "docs",
		Document: "README.md",
		Reasons:  []check.Reason{createTestReason("Old reason", check.DiffCheckTypeLLM, "main.go", "old-hash", true)},
		Checked:  true,
		Outdated: true,
	}

	contextHashes := check.FileCheckContextHashes{
		"main.go": {check.DiffCheckTypeLLM: "new-hash"},
	}

	result := mergeCheckRecommendations([]CheckRecommendation{newRec}, []CheckRecommendation{existingRec}, contextHashes)

	if len(result) != 1 {
		t.Errorf("Expected 1 recommendation, got %d", len(result))
	}
	if len(result[0].Reasons) != 1 {
		t.Errorf("Expected 1 reason, got %d", len(result[0].Reasons))
	}
	if !result[0].Checked {
		t.Errorf("Expected to preserve checked state")
	}
	if result[0].Reasons[0].Outdated {
		t.Errorf("Expected reason not to be outdated")
	}
	if result[0].Reasons[0].Reason != "New reason" {
		t.Errorf("Expected new reason, got: %s", result[0].Reasons[0].Reason)
	}
}

func TestCheckPRParseComment(t *testing.T) {
	recs := []CheckRecommendation{
		{
			Source:   "docs",
			Document: "README.md",
			Section:  []string{"getting-started"},
			Reasons:  []check.Reason{createTestReason("API endpoint changed", check.DiffCheckTypeLLM, "api.go", "hash1", false)},
		},
		{
			Source:   "docs",
			Document: "API.md",
			Reasons: []check.Reason{
				createTestReason("Old requirement", check.DiffCheckTypeLLM, "old-file.go", "old-hash", true),
				createTestReason("This is outdated", check.DiffCheckTypeUpdateIfTouched, "removed.go", "outdated-hash", true),
			},
			Checked:  true,
			Outdated: true,
		},
	}

	output := CheckOutput{
		Recommendations: recs,
		Head:            "commit-sha-123",
		Base:            "main",
	}

	formattedComment := formatCheckPRComment(&output)

	if !strings.Contains(formattedComment, "### Recommendations") {
		t.Errorf("Missing recommendations section")
	}
	if !strings.Contains(formattedComment, "Changes have caused the following recommendations to be outdated") {
		t.Errorf("Missing outdated recommendations section")
	}

	formattedComment = strings.ReplaceAll(formattedComment, "- [ ]", "- [x]")

	parsedOutput, err := parseCheckPRComment(formattedComment)
	if err != nil {
		t.Fatal(err)
	}

	if len(parsedOutput.Recommendations) != 2 {
		t.Errorf("Expected 2 recommendations, got %d", len(parsedOutput.Recommendations))
	}
	if parsedOutput.Head != "commit-sha-123" {
		t.Errorf("Expected head 'commit-sha-123', got '%s'", parsedOutput.Head)
	}
	if parsedOutput.Base != "main" {
		t.Errorf("Expected base 'main', got '%s'", parsedOutput.Base)
	}

	// Non-outdated first, then by document: README.md, API.md
	currentRec := &parsedOutput.Recommendations[0]
	outdatedRec := &parsedOutput.Recommendations[1]

	if currentRec.Document != "README.md" {
		t.Errorf("Expected first to be README.md, got %s", currentRec.Document)
	}
	if !currentRec.Checked {
		t.Errorf("Expected current rec to be checked")
	}
	if currentRec.Outdated {
		t.Errorf("Expected current rec to not be outdated")
	}

	if outdatedRec.Document != "API.md" {
		t.Errorf("Expected second to be API.md, got %s", outdatedRec.Document)
	}
	if !outdatedRec.Outdated {
		t.Errorf("Expected outdated rec to be outdated")
	}
	if !outdatedRec.Checked {
		t.Errorf("Expected outdated rec to preserve checked state")
	}

	if len(currentRec.Reasons) != 1 {
		t.Errorf("Expected current rec to have 1 reason, got %d", len(currentRec.Reasons))
	}
	if len(outdatedRec.Reasons) != 2 {
		t.Errorf("Expected outdated rec to have 2 reasons, got %d", len(outdatedRec.Reasons))
	}
}
