package config

import (
	"os"
	"path"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	llmKey := "TEST_KEY"

	expectedConfig := Config{
		LLM: LLM{
			Provider: "anthropic",
			Model:    "claude-3-5-sonnet-20241022",
			Key:      llmKey,
		},
		Systems: []System{{
			ID: "my-app",
			Code: []Code{{
				ID:        "app",
				Extractor: "fs",
				Path:      "./",
				Include:   []string{"package.json", "./**/*.js"},
				Exclude:   []string{"./**/*.test.js"},
			}},
			Docs: []Doc{{
				ID:        "docs",
				Type:      "md",
				Extractor: "fs",
				Path:      "./",
				Include:   []string{"./**/*.md"},
			}},
		}},
	}

	os.Setenv("ANTHROPIC_KEY", llmKey)

	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Could not get cwd: %v", err)
	}
	absPath := path.Join(dir, "test_config.yml")
	cfg, err := Load(absPath)
	if err != nil {
		t.Fatalf("Could not get config: %v", err)
	}

	if !reflect.DeepEqual(*cfg, expectedConfig) {
		t.Fatalf("Expected config to match. Got %v, Wanted %v", *cfg, expectedConfig)
	}
}

func TestValidate(t *testing.T) {
	code := Code{
		ID: "1234",
	}
	doc := Doc{
		ID:   "1234",
		Type: "md",
	}
	invalidDoc := Doc{
		ID:   "1234",
		Type: "invalid",
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
		{[]Code{code}, []Doc{invalidDoc}, true},
	}

	for _, test := range tests {
		cfg := &Config{
			Systems: []System{{
				ID:   "test-system",
				Code: test.code,
				Docs: test.docs,
			}},
		}

		err := validate(cfg)
		if (err == nil && test.shouldError) || (err != nil && !test.shouldError) {
			t.Errorf("got %v, want %t", err, test.shouldError)
		}
	}
}
