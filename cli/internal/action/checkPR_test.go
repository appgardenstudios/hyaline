package action

import (
	"hyaline/internal/check"
	"reflect"
	"strings"
	"testing"
)

func TestCheckPRMergeRecs(t *testing.T) {
	newRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "newRec1"}},
		Checked:  false,
	}
	existingRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "existingRec1"}},
		Checked:  false,
	}
	newToExistingRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "newRec1"}},
		Checked:  false,
	}
	mergedRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		Checked:  false,
	}
	newRec1wSection := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Section:  []string{"section1", "section2"},
		Reasons:  []check.Reason{{Reason: "newRec1"}},
		Checked:  false,
	}
	existingRec1wSection := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Section:  []string{"section1", "section2"},
		Reasons:  []check.Reason{{Reason: "existingRec1"}},
		Checked:  false,
	}
	mergedRec1wSection := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Section:  []string{"section1", "section2"},
		Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		Checked:  false,
	}
	newRec2 := CheckRecommendation{
		Source:   "source",
		Document: "doc0",
		Reasons:  []check.Reason{{Reason: "newRec2"}},
		Checked:  false,
	}
	newToExistingRec2 := CheckRecommendation{
		Source:   "source",
		Document: "doc0",
		Reasons:  []check.Reason{{Reason: "newRec2"}},
		Checked:  false,
	}
	newChangedRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "newRec1"}},
		Checked:  true,
	}
	existingCheckedRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "existingRec1"}},
		Checked:  true,
	}
	mergedCheckedRec1 := CheckRecommendation{
		Source:   "source",
		Document: "doc1",
		Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		Checked:  true,
	}

	tests := []struct {
		newRecs      []CheckRecommendation
		existingRecs []CheckRecommendation
		mergedRecs   []CheckRecommendation
	}{
		// 1 existing rec, no new recs
		{
			[]CheckRecommendation{},
			[]CheckRecommendation{existingRec1},
			[]CheckRecommendation{existingRec1},
		},
		// 0 existing recs, 1 new rec
		{
			[]CheckRecommendation{newRec1},
			[]CheckRecommendation{},
			[]CheckRecommendation{newToExistingRec1},
		},
		// Existing and new recs the same (+ merge reasons)
		{
			[]CheckRecommendation{newRec1},
			[]CheckRecommendation{existingRec1},
			[]CheckRecommendation{mergedRec1},
		},
		// Existing and new recs the same w/ sections (+ merge reasons)
		{
			[]CheckRecommendation{newRec1wSection},
			[]CheckRecommendation{existingRec1wSection},
			[]CheckRecommendation{mergedRec1wSection},
		},
		// Existing rec and new rec (new rec sorts before existing rec)
		{
			[]CheckRecommendation{newRec2},
			[]CheckRecommendation{existingRec1},
			[]CheckRecommendation{newToExistingRec2, existingRec1},
		},
		// Existing rec w/ section, new rec same document (but not section)
		{
			[]CheckRecommendation{newRec1},
			[]CheckRecommendation{existingRec1wSection},
			[]CheckRecommendation{newToExistingRec1, existingRec1wSection},
		},
		// Existing rec unchecked, new rec changed
		{
			[]CheckRecommendation{newChangedRec1},
			[]CheckRecommendation{existingRec1},
			[]CheckRecommendation{mergedCheckedRec1},
		},
		// Existing rec checked, new rec unchanged
		{
			[]CheckRecommendation{newRec1},
			[]CheckRecommendation{existingCheckedRec1},
			[]CheckRecommendation{mergedCheckedRec1},
		},
	}

	for i, test := range tests {
		recs := mergeCheckRecommendations(test.newRecs, test.existingRecs)

		if !reflect.DeepEqual(recs, test.mergedRecs) {
			t.Errorf("%d: got %v, expected %v", i, recs, test.mergedRecs)
		}
	}
}

func TestCheckPRParseComment(t *testing.T) {
	recs := []CheckRecommendation{
		{
			Source:   "source",
			Document: "document",
			Section:  []string{"section1", "section2"},
			Reasons:  []check.Reason{{Reason: "reason1"}, {Reason: "reason2"}},
			Checked:  false,
		},
	}
	output := CheckOutput{
		Recommendations: recs,
		Head:            "sha",
		Base:            "base",
	}
	formattedComment := formatCheckPRComment(&output)

	// mark everything as checked
	formattedComment = strings.ReplaceAll(formattedComment, "- [ ]", "- [x]")

	expectedOutput := CheckOutput{
		Recommendations: []CheckRecommendation{
			{
				Source:   "source",
				Document: "document",
				Section:  []string{"section1", "section2"},
				Reasons:  []check.Reason{{Reason: "reason1"}, {Reason: "reason2"}},
				Checked:  true,
			},
		},
		Head: "sha",
		Base: "base",
	}

	existingOutput, err := parseCheckPRComment(formattedComment)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedOutput, *existingOutput) {
		t.Errorf("got %v, expected %v", *existingOutput, expectedOutput)
	}
}
