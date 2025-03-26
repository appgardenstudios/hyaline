package config

import "testing"

func TestValidate(t *testing.T) {
	code := Code{
		ID: "1234",
	}
	doc := Doc{
		ID: "1234",
	}

	var tests = []struct {
		code        []Code
		docs        []Doc
		shouldError bool
	}{
		{[]Code{}, []Doc{}, false},
		{[]Code{code}, []Doc{}, false},
		{[]Code{}, []Doc{doc}, false},
		{[]Code{code}, []Doc{doc}, false},
		{[]Code{code, code}, []Doc{doc}, true},
		{[]Code{code}, []Doc{doc, doc}, true},
	}

	for _, test := range tests {
		cfg := Config{
			Systems: []System{{
				ID:   "test-system",
				Code: test.code,
				Docs: test.docs,
			}},
		}

		err := validate(&cfg)
		if (err == nil && test.shouldError) || (err != nil && !test.shouldError) {
			t.Errorf("got %v, want %t", err, test.shouldError)
		}
	}
}
