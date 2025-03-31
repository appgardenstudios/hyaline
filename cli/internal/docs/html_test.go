package docs

import "testing"

func TestLoad(t *testing.T) {
	html := `
<!doctype html>
<html lang=en>
<head>
  <meta charset=utf-8>
  <title>blah</title>
</head>
<body>
  <nav>
    <ul>
      <li>Nav 1</li>
      <li>Nav 2</li>
      <li>Nav 3</li>
    </ul>
  </nav>
  <main>
    <p>I am the content</p>
  
    <h1>First Section</h1>
		<p>Some section one content</p>
  </main>
</body>
</html>
	`
	expectedBody := `- Nav 1
- Nav 2
- Nav 3

I am the content

# First Section

Some section one content`
	expectedMain := `I am the content

# First Section

Some section one content`

	var tests = []struct {
		selector    string
		markdown    string
		shouldError bool
	}{
		{"", expectedBody, false},
		{"main", expectedMain, false},
		{"invalid", "", true},
		{"missing", "", true},
	}

	for _, test := range tests {
		result, err := extractHTMLDocument(html, test.selector)

		if test.shouldError {
			if err == nil {
				t.Errorf("should have errored for selector %s", test.selector)
			}
		} else {
			if result != test.markdown {
				t.Errorf("for selector '%s', got %s, want %s", test.selector, result, test.markdown)
			}
		}
	}

}
