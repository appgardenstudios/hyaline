package extract

import (
	"context"
	"fmt"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

type section struct {
	Parent   *section
	Depth    int
	Name     string
	FullName string
	Content  string
	Purpose  string
	Children []*section
}

func extractSections(documentID string, sourceID string, markdown string, extractPurpose bool, purposeKey string, db *sqlite.Queries) error {

	// Get our tree of sections
	sections := getMarkdownSections(strings.Split(markdown, "\n"))

	// Extract purpose (if enabled)
	if extractPurpose {
		extractMarkdownSectionPurposes(sections, purposeKey)
	}

	// Insert our sections
	return insertSections(sections, 0, documentID, sourceID, db)
}

func getMarkdownSections(lines []string) *section {
	// Create our root
	root := &section{
		Parent:   nil,
		Depth:    0,
		Name:     "",
		FullName: "",
		Content:  "",
		Children: []*section{},
	}
	current := root

	// Track original and generated full names globally
	// We use map[string]struct{} instead of []string for O(1) lookups vs O(n) array searches.
	// Values are empty structs since we only care about keys for existence checks.
	originalFullNames := make(map[string]struct{})
	generatedFullNames := make(map[string]struct{})

	inCodeBlock := false

	for _, line := range lines {
		// If the line starts with ```, enter or exit the code block
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
		}

		// If line starts with #, modify current to the correct level
		level := countPounds(line)
		if level == 0 || inCodeBlock {
			// Add to current section
			current.Content = current.Content + "\n" + line
		} else {
			// recurse up to put this section where it goes
			for current.Depth >= level {
				current = current.Parent
			}
			// Trim spaces and strip out all # and /
			name := strings.TrimSpace(strings.ReplaceAll(line[level:], "/", "_"))

			fullName := name
			if current.FullName != "" {
				fullName = fmt.Sprintf("%s/%s", current.FullName, name)
			}

			// Deduplicate the full name
			uniqueFullName := deduplicateSectionName(fullName, originalFullNames, generatedFullNames)

			newSection := &section{
				Parent:   current,
				Depth:    level,
				Name:     name,
				FullName: uniqueFullName,
				Content:  "",
				Children: []*section{},
			}
			current.Children = append(current.Children, newSection)
			current = newSection
		}

		// Insert this line up the chain to the root
		parent := current.Parent
		for parent != nil {
			parent.Content = parent.Content + "\n" + line
			parent = parent.Parent
		}
	}

	return root
}

func countPounds(line string) int {
	count := 0

	for _, c := range line {
		if c == '#' {
			count++
		} else {
			break
		}
	}

	return count
}

// deduplicateSectionName implements a two-step algorithm to handle duplicate section names:
// Step 1: Keep adding " (1)" until no conflict with original section names
// Step 2: Then increment the final number until no conflict with either original OR generated section names
func deduplicateSectionName(name string, originalNames map[string]struct{}, generatedNames map[string]struct{}) string {
	// If the name is not in original names AND not in generated names, it's unique
	if _, existsOriginal := originalNames[name]; !existsOriginal {
		if _, existsGenerated := generatedNames[name]; !existsGenerated {
			originalNames[name] = struct{}{}
			return name
		}
	}

	// Step 1: Keep adding " (1)" until no conflict with original section names
	candidate := name
	for {
		if _, exists := originalNames[candidate+" (1)"]; !exists {
			break
		}
		candidate = candidate + " (1)"
	}

	// Step 2: Increment the final number until no conflict with either original OR generated section names
	counter := 1
	for {
		finalName := fmt.Sprintf("%s (%d)", candidate, counter)

		// Check if this name conflicts with original or generated names
		if _, existsOriginal := originalNames[finalName]; existsOriginal {
			counter++
			continue
		}
		if _, existsGenerated := generatedNames[finalName]; existsGenerated {
			counter++
			continue
		}

		// Found a unique name
		generatedNames[finalName] = struct{}{}
		return finalName
	}
}

func extractMarkdownSectionPurposes(section *section, purposeKey string) {
	content := strings.TrimSpace(section.Content)

	// If the first line of the section content starts with <!--, attempt to extract purpose from the comment
	if strings.HasPrefix(content, "<!--") {
		lines := strings.Split(content, "\n")
		metadata := extractHTMLComment(lines)
		section.Purpose = extractPurpose(metadata, purposeKey)
	}

	// Extract child section purposes
	for _, child := range section.Children {
		extractMarkdownSectionPurposes(child, purposeKey)
	}
}

func insertSections(s *section, order int, documentID string, sourceID string, db *sqlite.Queries) error {
	// If Parent is nil, it's the root document and we don't insert it
	if s.Parent != nil {
		err := db.InsertSection(context.Background(), sqlite.InsertSectionParams{
			ID:            s.FullName,
			DocumentID:    documentID,
			SourceID:      sourceID,
			ParentID:      s.Parent.FullName,
			PeerOrder:     order,
			Name:          s.Name,
			Purpose:       s.Purpose,
			ExtractedData: strings.TrimSpace(s.Content),
		})
		if err != nil {
			slog.Debug("extract.insertSections could not insert section", "error", err)
			return err
		}
	}

	// Insert children
	for i, child := range s.Children {
		err := insertSections(child, i, documentID, sourceID, db)
		if err != nil {
			return err
		}
	}

	return nil
}
