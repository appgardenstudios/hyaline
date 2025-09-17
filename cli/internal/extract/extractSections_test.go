package extract

import (
	"reflect"
	"testing"
)

func TestGetMarkdownSections_BasicFunctionality(t *testing.T) {
	lines := []string{
		"# Section A",
		"Some content",
		"## Subsection A1",
		"More content",
		"# Section B",
		"Different content",
	}

	root := getMarkdownSections(lines)
	actual := getAllFullNames(root)
	expected := []string{
		"Section A",
		"Section A/Subsection A1",
		"Section B",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FullNames don't match.\nActual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetMarkdownSections_DuplicateNames(t *testing.T) {
	lines := []string{
		"# Section A (1)",
		"# Section B",
		"# Section A",
		"# Section A",
		"# Section B",
		"# Section B",
		"# Section A (1)",
		"# Section B (2)",
		"## Section A",
		"## Section B",
	}

	root := getMarkdownSections(lines)
	actual := getAllFullNames(root)
	expected := []string{
		"Section A (1)",
		"Section B",
		"Section A",
		"Section A (1) (1)",
		"Section B (1)",
		"Section B (2)",
		"Section A (1) (2)",
		"Section B (2) (1)",
		"Section B (2) (1)/Section A",
		"Section B (2) (1)/Section B",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FullNames don't match.\nActual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetMarkdownSections_ComplexDuplicates(t *testing.T) {
	lines := []string{
		"# Test",
		"# Test (1)",
		"# Test",
		"# Test (1) (1)",
		"# Test",
		"# Test (1)",
	}

	root := getMarkdownSections(lines)
	actual := getAllFullNames(root)
	expected := []string{
		"Test",
		"Test (1)",
		"Test (1) (1)",
		"Test (1) (1) (1)",
		"Test (1) (2)",
		"Test (1) (3)",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FullNames don't match.\nActual: %v\nExpected: %v", actual, expected)
	}
}

func TestGetMarkdownSections_MultiLevelDuplicates(t *testing.T) {
	lines := []string{
		"# Parent",
		"## Child",
		"## Child",
		"# Parent",
		"## Child",
		"## Child (1)",
		"## Child",
	}

	root := getMarkdownSections(lines)
	actual := getAllFullNames(root)
	expected := []string{
		"Parent",
		"Parent/Child",
		"Parent/Child (1)",
		"Parent (1)",
		"Parent (1)/Child",
		"Parent (1)/Child (1)",
		"Parent (1)/Child (1) (1)",
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("FullNames don't match.\nActual: %v\nExpected: %v", actual, expected)
	}
}

func TestExtractMarkdownSectionPurposes(t *testing.T) {
	noContent := &section{}
	emptyContent := &section{
		Content: "",
	}
	noComment := &section{
		Content: "section information\nanother line",
	}
	comment := &section{
		Content: "<!-- key: value -->\nsection information\nanother line",
	}
	purpose := &section{
		Content: "<!-- \npurpose: value\n -->\nsection information\nanother line",
	}

	s := &section{
		Content:  "<!-- purpose: main purpose -->section content",
		Children: []*section{noContent, emptyContent, noComment, comment, purpose},
	}

	expected := "main purpose"
	childrenExpected := []string{"", "", "", "", "value"}

	extractMarkdownSectionPurposes(s, "purpose")

	if s.Purpose != expected {
		t.Errorf("%s - expected %s, got %s", "s", expected, s.Purpose)
	}

	for i, purpose := range childrenExpected {
		if s.Children[i].Purpose != purpose {
			t.Errorf("%d - expected %s, got %s", i, purpose, s.Children[i].Purpose)
		}
	}
}

// getAllFullNames recursively collects all FullName values from a section tree
func getAllFullNames(s *section) []string {
	var fullNames []string

	// Add current section's FullName (skip root with empty FullName)
	if s.FullName != "" {
		fullNames = append(fullNames, s.FullName)
	}

	// Add children's FullNames
	for _, child := range s.Children {
		fullNames = append(fullNames, getAllFullNames(child)...)
	}

	return fullNames
}
