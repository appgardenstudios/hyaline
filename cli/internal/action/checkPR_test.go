package action

import (
	"hyaline/internal/check"
	"reflect"
	"strings"
	"testing"
)

func TestCheckPRMergeRecs(t *testing.T) {
	newRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "newRec1"}},
		},
		Checked: false,
	}
	existingRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "existingRec1"}},
		},
		Checked: false,
	}
	newToExistingRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "newRec1"}},
		},
		Checked: false,
	}
	mergedRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		},
		Checked: false,
	}
	newRec1wSection := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Section:  []string{"section1", "section2"},
			Reasons:  []check.Reason{{Reason: "newRec1"}},
		},
		Checked: false,
	}
	existingRec1wSection := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Section:  []string{"section1", "section2"},
			Reasons:  []check.Reason{{Reason: "existingRec1"}},
		},
		Checked: false,
	}
	mergedRec1wSection := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Section:  []string{"section1", "section2"},
			Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		},
		Checked: false,
	}
	newRec2 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc0",
			Reasons:  []check.Reason{{Reason: "newRec2"}},
		},
		Checked: false,
	}
	newToExistingRec2 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc0",
			Reasons:  []check.Reason{{Reason: "newRec2"}},
		},
		Checked: false,
	}
	newChangedRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "newRec1"}},
		},
		Checked: true,
	}
	existingCheckedRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "existingRec1"}},
		},
		Checked: true,
	}
	mergedCheckedRec1 := CheckPRCommentRecommendation{
		CheckRecommendation: CheckRecommendation{
			Source:   "source",
			Document: "doc1",
			Reasons:  []check.Reason{{Reason: "existingRec1"}, {Reason: "newRec1"}},
		},
		Checked: true,
	}

	tests := []struct {
		newRecs      []CheckPRCommentRecommendation
		existingRecs []CheckPRCommentRecommendation
		mergedRecs   []CheckPRCommentRecommendation
	}{
		// 1 existing rec, no new recs
		{
			[]CheckPRCommentRecommendation{},
			[]CheckPRCommentRecommendation{existingRec1},
			[]CheckPRCommentRecommendation{existingRec1},
		},
		// 0 existing recs, 1 new rec
		{
			[]CheckPRCommentRecommendation{newRec1},
			[]CheckPRCommentRecommendation{},
			[]CheckPRCommentRecommendation{newToExistingRec1},
		},
		// Existing and new recs the same (+ merge reasons)
		{
			[]CheckPRCommentRecommendation{newRec1},
			[]CheckPRCommentRecommendation{existingRec1},
			[]CheckPRCommentRecommendation{mergedRec1},
		},
		// Existing and new recs the same w/ sections (+ merge reasons)
		{
			[]CheckPRCommentRecommendation{newRec1wSection},
			[]CheckPRCommentRecommendation{existingRec1wSection},
			[]CheckPRCommentRecommendation{mergedRec1wSection},
		},
		// Existing rec and new rec (new rec sorts before existing rec)
		{
			[]CheckPRCommentRecommendation{newRec2},
			[]CheckPRCommentRecommendation{existingRec1},
			[]CheckPRCommentRecommendation{newToExistingRec2, existingRec1},
		},
		// Existing rec w/ section, new rec same document (but not section)
		{
			[]CheckPRCommentRecommendation{newRec1},
			[]CheckPRCommentRecommendation{existingRec1wSection},
			[]CheckPRCommentRecommendation{newToExistingRec1, existingRec1wSection},
		},
		// Existing rec unchecked, new rec changed
		{
			[]CheckPRCommentRecommendation{newChangedRec1},
			[]CheckPRCommentRecommendation{existingRec1},
			[]CheckPRCommentRecommendation{mergedCheckedRec1},
		},
		// Existing rec checked, new rec unchanged
		{
			[]CheckPRCommentRecommendation{newRec1},
			[]CheckPRCommentRecommendation{existingCheckedRec1},
			[]CheckPRCommentRecommendation{mergedCheckedRec1},
		},
	}

	for i, test := range tests {
		recs := mergeCheckPRRecs(test.newRecs, test.existingRecs)

		if !reflect.DeepEqual(recs, test.mergedRecs) {
			t.Errorf("%d: got %v, expected %v", i, recs, test.mergedRecs)
		}
	}
}

func TestCheckPRParseComment(t *testing.T) {
	recs := []CheckPRCommentRecommendation{
		{
			CheckRecommendation: CheckRecommendation{
				Source:   "source",
				Document: "document",
				Section:  []string{"section1", "section2"},
				Reasons:  []check.Reason{{Reason: "reason1"}, {Reason: "reason2"}},
			},
			Checked: false,
		},
	}
	rawData, err := formatCheckPRRawData(&recs)
	if err != nil {
		t.Fatal(err)
	}
	comment := CheckPRComment{
		Sha:             "sha",
		Recommendations: recs,
		RawData:         rawData,
	}
	formattedComment := formatCheckPRComment(&comment)

	// mark everything as checked
	formattedComment = strings.ReplaceAll(formattedComment, "- [ ]", "- [x]")

	expectedRecs := []CheckPRCommentRecommendation{
		{
			CheckRecommendation: CheckRecommendation{
				Source:   "source",
				Document: "document",
				Section:  []string{"section1", "section2"},
				Reasons:  []check.Reason{{Reason: "reason1"}, {Reason: "reason2"}},
			},
			Checked: true,
		},
	}

	existingRecs, err := parseCheckPRComment(formattedComment)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedRecs, existingRecs) {
		t.Errorf("got %v, expected %v", existingRecs, expectedRecs)
	}
}
