package docs

import (
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

// Based off of https://github.com/gomarkdown/markdown/blob/master/md/md_renderer.go
type textRenderer struct{}

func (r *textRenderer) RenderNode(w io.Writer, node ast.Node, entering bool) ast.WalkStatus {
	switch node := node.(type) {
	case *ast.Text:
		if entering {
			io.WriteString(w, string(node.Literal))
			io.WriteString(w, "\n")
		}
	case *ast.Heading:
		if entering {
			io.WriteString(w, string(node.Literal))
			io.WriteString(w, "\n")
		}
	// Note: This will eventually need to support additional types
	default:
		// Do nothing
	}

	return ast.GoToNext
}
func (r *textRenderer) RenderHeader(w io.Writer, ast ast.Node) {
	// do nothing
}

func (r *textRenderer) RenderFooter(w io.Writer, ast ast.Node) {
	// do nothing
}

func extractMarkdownText(content []byte) string {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(content)

	renderer := &textRenderer{}

	return strings.TrimSpace(string(markdown.Render(doc, renderer)))
}

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
