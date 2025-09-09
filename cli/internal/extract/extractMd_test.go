package extract

import "testing"

func TestExtractMdDocumentPurpose(t *testing.T) {
	empty := ""
	emptyFrontMatter := `---
---`
	basicFrontMatter := `---
purpose: I'm Here!
---`
	complexFrontMatter := `---
not-purpose: true
purpose: document purpose
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
purpose: My Purpose!
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
