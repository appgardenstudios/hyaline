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
			CodeSources: []CodeSource{
				{
					ID: "app-fs",
					Extractor: Extractor{
						Type: "fs",
						Options: ExtractorOptions{
							Path: "./",
						},
						Include: []string{"package.json", "./**/*.js"},
						Exclude: []string{"./**/*.test.js"},
					},
				}, {
					ID: "app-git-http",
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "http",
								Options: ExtractorAuthOptions{
									Username: "bob",
									Password: "nope",
								},
							},
						},
						Include: []string{"package.json", "./**/*.js"},
						Exclude: []string{"./**/*.test.js"},
					},
				}, {
					ID: "app-git-ssh",
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "ssh",
								Options: ExtractorAuthOptions{
									User:     "bob",
									PEM:      "my-pem",
									Password: "nope",
								},
							},
						},
						Include: []string{"package.json", "./**/*.js"},
						Exclude: []string{"./**/*.test.js"},
					},
				},
			},
			DocumentationSources: []DocumentationSource{
				{
					ID:   "md-docs-fs",
					Type: "md",
					Extractor: Extractor{
						Type: "fs",
						Options: ExtractorOptions{
							Path: "./",
						},
						Include: []string{"./**/*.md"},
					},
				}, {
					ID:   "md-docs-git-http",
					Type: "md",
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "http",
								Options: ExtractorAuthOptions{
									Username: "bob",
									Password: "nope",
								},
							},
						},
						Include: []string{"./**/*.md"},
					},
				}, {
					ID:   "md-docs-git-ssh",
					Type: "md",
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "ssh",
								Options: ExtractorAuthOptions{
									User:     "bob",
									PEM:      "my-pem",
									Password: "nope",
								},
							},
						},
						Include: []string{"./**/*.md"},
					},
				}, {
					ID:   "html-docs-fs",
					Type: "html",
					Options: DocumentationSourceOptions{
						Selector: "main",
					},
					Extractor: Extractor{
						Type: "fs",
						Options: ExtractorOptions{
							Path: "./",
						},
						Include: []string{"./**/*.md"},
					},
				}, {
					ID:   "html-docs-git-http",
					Type: "html",
					Options: DocumentationSourceOptions{
						Selector: "main",
					},
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "http",
								Options: ExtractorAuthOptions{
									Username: "bob",
									Password: "nope",
								},
							},
						},
						Include: []string{"./**/*.md"},
					},
				}, {
					ID:   "html-docs-git-ssh",
					Type: "html",
					Options: DocumentationSourceOptions{
						Selector: "main",
					},
					Extractor: Extractor{
						Type: "git",
						Options: ExtractorOptions{
							Repo:   "git@github.com:appgardenstudios/hyaline-example.git",
							Branch: "main",
							Path:   "my/path",
							Clone:  true,
							Auth: ExtractorAuth{
								Type: "ssh",
								Options: ExtractorAuthOptions{
									User:     "bob",
									PEM:      "my-pem",
									Password: "nope",
								},
							},
						},
						Include: []string{"./**/*.md"},
					},
				},
			},
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
		ID: "1234",
		Extractor: Extractor{
			Type:    "fs",
			Include: []string{"**/*.js"},
			Exclude: []string{"**/*.test.js"},
		},
	}
	invalidCodeInclude := CodeSource{
		ID: "1234",
		Extractor: Extractor{
			Type:    "fs",
			Include: []string{"{a"},
		},
	}
	invalidCodeExclude := CodeSource{
		ID: "1234",
		Extractor: Extractor{
			Type:    "fs",
			Exclude: []string{"{a"},
		},
	}
	invalidCodeExtractor := CodeSource{
		ID: "1234",
		Extractor: Extractor{
			Type: "invalid",
		},
	}
	doc := DocumentationSource{
		ID:   "1234",
		Type: "md",
		Extractor: Extractor{
			Type:    "fs",
			Include: []string{"**/*.md"},
			Exclude: []string{"random.md"},
		},
	}
	invalidDoc := DocumentationSource{
		ID:   "1234",
		Type: "invalid",
		Extractor: Extractor{
			Type: "fs",
		},
	}
	invalidDocInclude := DocumentationSource{
		ID:   "1234",
		Type: "md",
		Extractor: Extractor{
			Type:    "fs",
			Include: []string{"{a"},
		},
	}
	invalidDocExclude := DocumentationSource{
		ID:   "1234",
		Type: "md",
		Extractor: Extractor{
			Type:    "fs",
			Include: []string{"{a"},
		},
	}
	invalidDocExtractor := DocumentationSource{
		ID:   "1234",
		Type: "md",
		Extractor: Extractor{
			Type: "invalid",
		},
	}
	invalidLLM := LLM{
		Provider: "invalid",
	}
	rule := DocumentSet{
		ID: "test",
	}

	var tests = []struct {
		llm         LLM
		code        []CodeSource
		docs        []DocumentationSource
		rules       []DocumentSet
		shouldError bool
	}{
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []DocumentSet{}, false},
		{LLM{}, []CodeSource{code}, []DocumentationSource{}, []DocumentSet{}, false},
		{LLM{}, []CodeSource{}, []DocumentationSource{doc}, []DocumentSet{}, false},
		{LLM{}, []CodeSource{code}, []DocumentationSource{doc}, []DocumentSet{}, false},
		{LLM{}, []CodeSource{code, code}, []DocumentationSource{doc}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{code}, []DocumentationSource{doc, doc}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{code}, []DocumentationSource{invalidDoc}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{invalidCodeInclude}, []DocumentationSource{}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{invalidCodeExclude}, []DocumentationSource{}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocInclude}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocExclude}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{invalidCodeExtractor}, []DocumentationSource{}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{invalidDocExtractor}, []DocumentSet{}, true},
		{invalidLLM, []CodeSource{}, []DocumentationSource{}, []DocumentSet{}, true},
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []DocumentSet{rule}, false},
		{LLM{}, []CodeSource{}, []DocumentationSource{}, []DocumentSet{rule, rule}, true},
	}

	for i, test := range tests {
		cfg := &Config{
			LLM: test.llm,
			Systems: []System{{
				ID:                   "test-system",
				CodeSources:          test.code,
				DocumentationSources: test.docs,
			}},
			CommonDocuments: test.rules,
		}

		err := validate(cfg)
		if (err == nil && test.shouldError) || (err != nil && !test.shouldError) {
			t.Logf("Error detected on test %d", i)
			t.Errorf("got %v, want %t", err, test.shouldError)
		}
	}
}
