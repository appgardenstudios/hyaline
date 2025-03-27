package docs

import (
	"fmt"
	"strings"
)

type markdownSection struct {
	Parent   *markdownSection
	Depth    int
	Name     string
	FullName string
	Content  string
	Children []*markdownSection
}

func getMarkdownSections(lines []string) *markdownSection {
	// Create our root
	root := &markdownSection{
		Parent:   nil,
		Depth:    0,
		Name:     "",
		FullName: "",
		Content:  "",
		Children: []*markdownSection{},
	}
	current := root

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
			newSection := &markdownSection{
				Parent:   current,
				Depth:    level,
				Name:     name,
				FullName: fmt.Sprintf("%s#%s", current.FullName, name),
				Content:  "",
				Children: []*markdownSection{},
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
