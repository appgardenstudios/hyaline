package docs

import (
	"fmt"
	"strings"
)

type section struct {
	Parent   *section
	Depth    int
	Name     string
	FullName string
	Content  string
	Children []*section
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

	for _, line := range lines {
		// If line starts with #, modify current to the correct level
		level := countPounds(line)
		if level == 0 {
			// Add to current section
			current.Content = current.Content + "\n" + line
		} else {
			// recurse up to put this section where it goes
			for current.Depth >= level {
				current = current.Parent
			}
			name := strings.TrimSpace(strings.ReplaceAll(line[level:], "#", ""))

			fullName := name
			if current.FullName != "" {
				fullName = fmt.Sprintf("%s#%s", current.FullName, name)
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
