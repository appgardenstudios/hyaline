package action

import (
	"reflect"
	"strings"
	"testing"
)

func TestMergeRecs(t *testing.T) {
	newRec1 := CheckChangeOutputEntry{
		System:              "system",
		DocumentationSource: "source",
		Document:            "doc1",
		Reasons:             []string{"newRec1"},
		Changed:             false,
	}
	existingRec1 := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Reasons:  []string{"existingRec1"},
	}
	newToExistingRec1 := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Reasons:  []string{"newRec1"},
	}
	mergedRec1 := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Reasons:  []string{"existingRec1", "newRec1"},
	}
	newRec1wSection := CheckChangeOutputEntry{
		System:              "system",
		DocumentationSource: "source",
		Document:            "doc1",
		Section:             []string{"section1", "section2"},
		Reasons:             []string{"newRec1"},
		Changed:             false,
	}
	existingRec1wSection := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Section:  []string{"section1", "section2"},
		Reasons:  []string{"existingRec1"},
	}
	mergedRec1wSection := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Section:  []string{"section1", "section2"},
		Reasons:  []string{"existingRec1", "newRec1"},
	}
	newRec2 := CheckChangeOutputEntry{
		System:              "system",
		DocumentationSource: "source",
		Document:            "doc0",
		Reasons:             []string{"newRec2"},
		Changed:             false,
	}
	newToExistingRec2 := UpdatePRCommentRecommendation{
		Checked:  false,
		System:   "system",
		Source:   "source",
		Document: "doc0",
		Reasons:  []string{"newRec2"},
	}
	newChangedRec1 := CheckChangeOutputEntry{
		System:              "system",
		DocumentationSource: "source",
		Document:            "doc1",
		Reasons:             []string{"newRec1"},
		Changed:             true,
	}
	existingCheckedRec1 := UpdatePRCommentRecommendation{
		Checked:  true,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Reasons:  []string{"existingRec1"},
	}
	mergedCheckedRec1 := UpdatePRCommentRecommendation{
		Checked:  true,
		System:   "system",
		Source:   "source",
		Document: "doc1",
		Reasons:  []string{"existingRec1", "newRec1"},
	}

	tests := []struct {
		newRecs      []CheckChangeOutputEntry
		existingRecs []UpdatePRCommentRecommendation
		mergedRecs   []UpdatePRCommentRecommendation
	}{
		// 1 existing rec, no new recs
		{
			[]CheckChangeOutputEntry{},
			[]UpdatePRCommentRecommendation{existingRec1},
			[]UpdatePRCommentRecommendation{existingRec1},
		},
		// 0 existing recs, 1 new rec
		{
			[]CheckChangeOutputEntry{newRec1},
			[]UpdatePRCommentRecommendation{},
			[]UpdatePRCommentRecommendation{newToExistingRec1},
		},
		// Existing and new recs the same (+ merge reasons)
		{
			[]CheckChangeOutputEntry{newRec1},
			[]UpdatePRCommentRecommendation{existingRec1},
			[]UpdatePRCommentRecommendation{mergedRec1},
		},
		// Existing and new recs the same w/ sections (+ merge reasons)
		{
			[]CheckChangeOutputEntry{newRec1wSection},
			[]UpdatePRCommentRecommendation{existingRec1wSection},
			[]UpdatePRCommentRecommendation{mergedRec1wSection},
		},
		// Existing rec and new rec (new rec sorts before existing rec)
		{
			[]CheckChangeOutputEntry{newRec2},
			[]UpdatePRCommentRecommendation{existingRec1},
			[]UpdatePRCommentRecommendation{newToExistingRec2, existingRec1},
		},
		// Existing rec w/ section, new rec same document (but not section)
		{
			[]CheckChangeOutputEntry{newRec1},
			[]UpdatePRCommentRecommendation{existingRec1wSection},
			[]UpdatePRCommentRecommendation{newToExistingRec1, existingRec1wSection},
		},
		// Existing rec unchecked, new rec changed
		{
			[]CheckChangeOutputEntry{newChangedRec1},
			[]UpdatePRCommentRecommendation{existingRec1},
			[]UpdatePRCommentRecommendation{mergedCheckedRec1},
		},
		// Existing rec checked, new rec unchanged
		{
			[]CheckChangeOutputEntry{newRec1},
			[]UpdatePRCommentRecommendation{existingCheckedRec1},
			[]UpdatePRCommentRecommendation{mergedCheckedRec1},
		},
	}

	for i, test := range tests {
		recs := mergeRecs(test.newRecs, test.existingRecs)

		if !reflect.DeepEqual(recs, test.mergedRecs) {
			t.Errorf("%d: got %v, expected %v", i, recs, test.mergedRecs)
		}
	}
}

func TestParsePRComment(t *testing.T) {
	recs := []UpdatePRCommentRecommendation{
		{
			Checked:  false,
			System:   "system",
			Source:   "source",
			Document: "document",
			Section:  []string{"section1", "section2"},
			Reasons:  []string{"reason1", "reason2"},
		},
	}
	rawData, err := formatRawData(&recs)
	if err != nil {
		t.Fatal(err)
	}
	comment := UpdatePRComment{
		Sha:             "sha",
		Recommendations: recs,
		RawData:         rawData,
	}
	formattedComment := formatPRComment(&comment)

	// mark everything as checked
	formattedComment = strings.ReplaceAll(formattedComment, "- [ ]", "- [x]")

	expectedRecs := []UpdatePRCommentRecommendation{
		{
			Checked:  true,
			System:   "system",
			Source:   "source",
			Document: "document",
			Section:  []string{"section1", "section2"},
			Reasons:  []string{"reason1", "reason2"},
		},
	}

	existingRecs, err := parsePRComment(formattedComment)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedRecs, existingRecs) {
		t.Errorf("got %v, expected %v", existingRecs, expectedRecs)
	}
}
