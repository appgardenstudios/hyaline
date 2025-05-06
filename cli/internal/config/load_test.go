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
			CodeSources: []CodeSource{{
				ID:        "app",
				Extractor: "fs",
				FsOptions: FsOptions{
					Path: "./",
				},
				GitOptions: GitOptions{
					Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
					Branch: "main",
					Path:   "my/path",
					Clone:  true,
					HTTPAuth: GitHTTPAuthOptions{
						Username: "bob",
						Password: "nope",
					},
					SSHAuth: GitSSHAuthOptions{
						User:     "bob",
						PEM:      "my-pem",
						Password: "nope",
					},
				},
				Include: []string{"package.json", "./**/*.js"},
				Exclude: []string{"./**/*.test.js"},
			}},
			DocumentationSources: []DocumentationSource{{
				ID:        "md-docs",
				Type:      "md",
				Extractor: "fs",
				FsOptions: FsOptions{
					Path: "./",
				},
				GitOptions: GitOptions{
					Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
					Branch: "main",
					Path:   "my/path",
					Clone:  true,
					HTTPAuth: GitHTTPAuthOptions{
						Username: "bob",
						Password: "nope",
					},
					SSHAuth: GitSSHAuthOptions{
						User:     "bob",
						PEM:      "my-pem",
						Password: "nope",
					},
				},
				Include: []string{"./**/*.md"},
			}, {
				ID:   "html-docs",
				Type: "html",
				HTML: DocHTMLOptions{
					Selector: "main",
				},
				Extractor: "fs",
				FsOptions: FsOptions{
					Path: "./",
				},
				GitOptions: GitOptions{
					Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
					Branch: "main",
					Path:   "my/path",
					Clone:  true,
					HTTPAuth: GitHTTPAuthOptions{
						Username: "bob",
						Password: "nope",
					},
					SSHAuth: GitSSHAuthOptions{
						User:     "bob",
						PEM:      "my-pem",
						Password: "nope",
					},
				},
				Include: []string{"./**/*.md"},
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

func TestGetEscapedEnv(t *testing.T) {
	var tests = []struct {
		env    string
		result string
	}{
		{"", ""},
		{"plain", "plain"},
		{`Line1
Line2`, `"Line1\nLine2"`},
		{`Line1"
Line2`, `"Line1\"\nLine2"`},
		{`Line1\nLine2`, `"Line1\nLine2"`},
		{`Line1"\nLine2`, `"Line1\"\nLine2"`},
		{"Line1\r\nLine2", `"Line1\nLine2"`},
	}

	for _, test := range tests {
		os.Setenv("TestGetEscapedEnv", test.env)
		result := getEscapedEnv("TestGetEscapedEnv")
		if result != test.result {
			t.Errorf("got %s, wanted %s", result, test.result)
		}
	}
}

func TestValidate(t *testing.T) {
	code := CodeSource{
		ID:        "1234",
		Extractor: "fs",
		Include:   []string{"**/*.js"},
		Exclude:   []string{"**/*.test.js"},
	}
	invalidCodeInclude := CodeSource{
		ID:        "1234",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidCodeExclude := CodeSource{
		ID:        "1234",
		Extractor: "fs",
		Exclude:   []string{"{a"},
	}
	invalidCodeExtractor := CodeSource{
		ID:        "1234",
		Extractor: "invalid",
	}
	doc := DocumentationSource{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"**/*.md"},
		Exclude:   []string{"random.md"},
	}
	invalidDoc := DocumentationSource{
		ID:        "1234",
		Type:      "invalid",
		Extractor: "fs",
	}
	invalidDocInclude := DocumentationSource{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidDocExclude := DocumentationSource{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidDocExtractor := DocumentationSource{
		ID:        "1234",
		Type:      "md",
		Extractor: "invalid",
	}
	invalidLLM := LLM{
		Provider: "invalid",
	}
	rule := RuleSet{
		ID: "test",
	}

	var tests = []struct {
		llm         LLM
		code        []CodeSource
		docs        []DocumentationSource
		rules       []RuleSet
		shouldError bool
	}{
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []RuleSet{}, false},
		{LLM{}, []CodeSource{code}, []DocumentationSource{}, []RuleSet{}, false},
		{LLM{}, []CodeSource{}, []DocumentationSource{doc}, []RuleSet{}, false},
		{LLM{}, []CodeSource{code}, []DocumentationSource{doc}, []RuleSet{}, false},
		{LLM{}, []CodeSource{code, code}, []DocumentationSource{doc}, []RuleSet{}, true},
		{LLM{}, []CodeSource{code}, []DocumentationSource{doc, doc}, []RuleSet{}, true},
		{LLM{}, []CodeSource{code}, []DocumentationSource{invalidDoc}, []RuleSet{}, true},
		{LLM{}, []CodeSource{invalidCodeInclude}, []DocumentationSource{}, []RuleSet{}, true},
		{LLM{}, []CodeSource{invalidCodeExclude}, []DocumentationSource{}, []RuleSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocInclude}, []RuleSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocExclude}, []RuleSet{}, true},
		{LLM{}, []CodeSource{invalidCodeExtractor}, []DocumentationSource{}, []RuleSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocExtractor}, []RuleSet{}, true},
		{invalidLLM, []CodeSource{}, []DocumentationSource{}, []RuleSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []RuleSet{rule}, false},
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []RuleSet{rule, rule}, true},
	}

	for i, test := range tests {
		cfg := &Config{
			LLM: test.llm,
			Systems: []System{{
				ID:                   "test-system",
				CodeSources:          test.code,
				DocumentationSources: test.docs,
			}},
			Rules: test.rules,
		}

		err := validate(cfg)
		if (err == nil && test.shouldError) || (err != nil && !test.shouldError) {
			t.Logf("Error detected on test %d", i)
			t.Errorf("got %v, want %t", err, test.shouldError)
		}
	}
}
