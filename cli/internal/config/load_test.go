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
			Docs: []Doc{{
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
	code := Code{
		ID:        "1234",
		Extractor: "fs",
		Include:   []string{"**/*.js"},
		Exclude:   []string{"**/*.test.js"},
	}
	invalidCodeInclude := Code{
		ID:        "1234",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidCodeExclude := Code{
		ID:        "1234",
		Extractor: "fs",
		Exclude:   []string{"{a"},
	}
	invalidCodeExtractor := Code{
		ID:        "1234",
		Extractor: "invalid",
	}
	doc := Doc{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"**/*.md"},
		Exclude:   []string{"random.md"},
	}
	invalidDoc := Doc{
		ID:        "1234",
		Type:      "invalid",
		Extractor: "fs",
	}
	invalidDocInclude := Doc{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidDocExclude := Doc{
		ID:        "1234",
		Type:      "md",
		Extractor: "fs",
		Include:   []string{"{a"},
	}
	invalidDocExtractor := Doc{
		ID:        "1234",
		Type:      "md",
		Extractor: "invalid",
	}
	invalidLLM := LLM{
		Provider: "invalid",
	}
	rule := Rule{
		ID: "test",
	}

	var tests = []struct {
		llm         LLM
		code        []Code
		docs        []Doc
		rules       []Rule
		shouldError bool
	}{
		{LLM{}, []Code{}, []Doc{}, []Rule{}, false},
		{LLM{}, []Code{code}, []Doc{}, []Rule{}, false},
		{LLM{}, []Code{}, []Doc{doc}, []Rule{}, false},
		{LLM{}, []Code{code}, []Doc{doc}, []Rule{}, false},
		{LLM{}, []Code{code, code}, []Doc{doc}, []Rule{}, true},
		{LLM{}, []Code{code}, []Doc{doc, doc}, []Rule{}, true},
		{LLM{}, []Code{code}, []Doc{invalidDoc}, []Rule{}, true},
		{LLM{}, []Code{invalidCodeInclude}, []Doc{}, []Rule{}, true},
		{LLM{}, []Code{invalidCodeExclude}, []Doc{}, []Rule{}, true},
		{LLM{}, []Code{}, []Doc{invalidDocInclude}, []Rule{}, true},
		{LLM{}, []Code{}, []Doc{invalidDocExclude}, []Rule{}, true},
		{LLM{}, []Code{invalidCodeExtractor}, []Doc{}, []Rule{}, true},
		{LLM{}, []Code{}, []Doc{invalidDocExtractor}, []Rule{}, true},
		{invalidLLM, []Code{}, []Doc{}, []Rule{}, true},
		{LLM{}, []Code{}, []Doc{}, []Rule{rule}, false},
		{LLM{}, []Code{}, []Doc{}, []Rule{rule, rule}, true},
	}

	for i, test := range tests {
		cfg := &Config{
			LLM: test.llm,
			Systems: []System{{
				ID:   "test-system",
				Code: test.code,
				Docs: test.docs,
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
