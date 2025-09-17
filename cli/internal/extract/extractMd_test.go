package extract

import (
	"strings"
	"testing"
)

func TestExtractMdDocumentPurpose(t *testing.T) {
	empty := ""
	emptyFrontMatter := `---
---`
	basicFrontMatter := `---
purpose: I'm Here!
---`
	complexFrontMatter := `---
not-purpose: true
purpose: "document purpose"
still-not-purpose: Blargh
---
`
	purposeOutsideFrontMatter := `---
not-purpose: true
still-not-purpose: Blargh
---
purpose: document purpose
# Section 1
`
	emptyComment := `<!-- -->`
	basicComment := `<!-- purpose: My Purpose! -->`
	multiLineCommentFirstLine := `<!-- purpose: My Purpose!
-->`
	multiLineComment := `<!--
purpose: "My Purpose!"
-->`
	multiLineCommentEndingLine := `<!--
purpose: My Purpose!-->`
	malformedComment := `<!-->`

	var tests = []struct {
		name     string
		document string
		key      string
		purpose  string
	}{
		{"Empty", empty, "purpose", ""},
		{"Empty Front Matter", emptyFrontMatter, "purpose", ""},
		{"Basic Front Matter", basicFrontMatter, "purpose", "I'm Here!"},
		{"Basic Front Matter Key Not Present", basicFrontMatter, "custom", ""},
		{"Complex Front Matter", complexFrontMatter, "purpose", "document purpose"},
		{"Purpose Outside Front Matter", purposeOutsideFrontMatter, "purpose", ""},
		{"Empty Comment", emptyComment, "purpose", ""},
		{"Basic Comment", basicComment, "purpose", "My Purpose!"},
		{"Multi Line Comment First Line", multiLineCommentFirstLine, "purpose", "My Purpose!"},
		{"Multi Line Comment", multiLineComment, "purpose", "My Purpose!"},
		{"Multi Line Comment Ending Line", multiLineCommentEndingLine, "purpose", "My Purpose!"},
		{"Malformed Comment", malformedComment, "purpose", ""},
	}

	for _, test := range tests {
		purpose := extractMdDocumentPurpose(test.document, test.key)

		if purpose != test.purpose {
			t.Errorf("%s - expected %s, got %s", test.name, purpose, test.purpose)
		}
	}

}

func TestExtractFrontMatter(t *testing.T) {
	basic := `---
purpose: here!
---
`
	multiline := `---
purpose: here!

another: line

---
`

	var tests = []struct {
		name     string
		lines    string
		contents string
	}{
		{"Empty String", "", ""},
		{"Empty Frontmatter", "---\n---", ""},
		{"No Frontmatter", "The contents line 1\nLine2\nLine3", ""},
		{"Basic", basic, "purpose: here!"},
		{"Multiline", multiline, "purpose: here!\n\nanother: line"},
	}

	for _, test := range tests {
		contents := extractFrontMatter(strings.Split(test.lines, "\n"))

		if contents != test.contents {
			t.Errorf("%s - expected %s, got %s", test.name, test.contents, contents)
		}
	}
}

func TestExtractHTMLComment(t *testing.T) {
	multiline := `<!--
The content
Line 2
-->`
	multiline2 := `<!--
The content
Line 2-->`

	var tests = []struct {
		name     string
		lines    string
		contents string
	}{
		{"Empty String", "", ""},
		{"No Comment", "The contents line 1\nLine2\nLine3", ""},
		{"Single Line", "<!-- The contents! -->", "The contents!"},
		{"Empty Single Line", "<!-- -->", ""},
		{"Comment Without end", "<!-->\nLine2", ">\nLine2"},
		{"Multiline", multiline, "The content\nLine 2"},
		{"Multiline2", multiline2, "The content\nLine 2"},
	}

	for _, test := range tests {
		contents := extractHTMLComment(strings.Split(test.lines, "\n"))

		if contents != test.contents {
			t.Errorf("%s - expected %s, got %s", test.name, test.contents, contents)
		}
	}
}

func TestExtractPurpose(t *testing.T) {
	var tests = []struct {
		name     string
		contents string
		key      string
		result   string
	}{
		{"Basic", "purpose: value", "purpose", "value"},
		{"Empty", "", "purpose", ""},
		{"Invalid YAML", "purpose:::: value", "purpose", ""},
		{"Missing", "other: value", "purpose", ""},
	}

	for _, test := range tests {
		purpose := extractPurpose(test.contents, test.key)

		if purpose != test.result {
			t.Errorf("%s - expected %s, got %s", test.name, test.result, purpose)
		}
	}
}
