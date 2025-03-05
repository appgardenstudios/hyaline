package docs

import "strings"

type section struct {
	Parent   *section
	Depth    int
	Title    string
	Content  string
	Children []*section
}

func getMarkdownSections(lines []string) *section {
	// Create our root
	root := &section{
		Parent:   nil,
		Depth:    0,
		Title:    "",
		Content:  "",
		Children: []*section{},
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
			newSection := &section{
				Parent:   current,
				Depth:    level,
				Title:    strings.TrimSpace(line[level:]),
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
