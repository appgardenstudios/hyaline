package config

import "testing"

func TestValidateExtract(t *testing.T) {
	validSource := ExtractSource{
		ID:          "my-app",
		Description: "The Description",
		Root:        "My Root",
	}
	invalidSourceID := ExtractSource{
		ID:          "my-app!",
		Description: "The Description",
		Root:        "My Root",
	}
	validCrawler := ExtractCrawler{
		Type: "fs",
	}
	invalidCrawlerType := ExtractCrawler{
		Type: "bogus",
	}
	invalidCrawlerInclude := ExtractCrawler{
		Type:    "fs",
		Include: []string{"{a"},
	}
	invalidCrawlerIncludeEmpty := ExtractCrawler{
		Type:    "fs",
		Include: []string{""},
	}
	invalidCrawlerExclude := ExtractCrawler{
		Type:    "fs",
		Exclude: []string{"{a"},
	}
	invalidCrawlerExcludeEmpty := ExtractCrawler{
		Type:    "fs",
		Exclude: []string{""},
	}
	validExtractors := []ExtractExtractor{
		{
			Type:    "md",
			Include: []string{"**/*"},
			Exclude: []string{"**/*"},
		},
	}
	invalidExtractorsEmpty := []ExtractExtractor{}
	invalidExtractorsType := []ExtractExtractor{
		{
			Type: "bogus",
		},
	}
	invalidExtractorsInclude := []ExtractExtractor{
		{
			Type:    "md",
			Include: []string{"{a"},
		},
	}
	invalidExtractorsIncludeEmpty := []ExtractExtractor{
		{
			Type:    "md",
			Include: []string{""},
		},
	}
	invalidExtractorsExclude := []ExtractExtractor{
		{
			Type:    "md",
			Exclude: []string{"{a"},
		},
	}
	invalidExtractorsExcludeEmpty := []ExtractExtractor{
		{
			Type:    "md",
			Exclude: []string{""},
		},
	}
	validMetadata := []ExtractMetadata{
		{
			Document: "**/*",
			Section:  "**/*",
			Tags: []ExtractMetadataTag{
				{
					Key:   "foo",
					Value: "bar",
				},
			},
		},
	}
	invalidMetadataDocument := []ExtractMetadata{
		{
			Document: "{a",
		},
	}
	invalidMetadataDocumentEmpty := []ExtractMetadata{
		{
			Document: "",
		},
	}
	invalidMetadataSection := []ExtractMetadata{
		{
			Document: "**/*",
			Section:  "{a",
		},
	}
	validMetadataSectionEmpty := []ExtractMetadata{
		{
			Document: "**/*",
			Section:  "",
		},
	}
	invalidMetadataTagKey := []ExtractMetadata{
		{
			Document: "**/*",
			Section:  "**/*",
			Tags: []ExtractMetadataTag{
				{
					Key:   "foo!",
					Value: "bar",
				},
			},
		},
	}
	invalidMetadataTagValue := []ExtractMetadata{
		{
			Document: "**/*",
			Section:  "**/*",
			Tags: []ExtractMetadataTag{
				{
					Key:   "foo",
					Value: "bar!",
				},
			},
		},
	}

	var tests = []struct {
		extract *Extract
		err     string
	}{
		{nil, ``},
		{&Extract{validSource, validCrawler, validExtractors, validMetadata}, ``},
		{&Extract{invalidSourceID, validCrawler, validExtractors, validMetadata}, `extract.source.id must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: my-app!`},
		{&Extract{validSource, invalidCrawlerType, validExtractors, validMetadata}, `extract.crawler.type must be one of fs, git, http, found: bogus`},
		{&Extract{validSource, invalidCrawlerInclude, validExtractors, validMetadata}, `extract.crawler.include[0] must be a valid pattern, found: {a`},
		{&Extract{validSource, invalidCrawlerIncludeEmpty, validExtractors, validMetadata}, `extract.crawler.include[0] must be a valid pattern, found: `},
		{&Extract{validSource, invalidCrawlerExclude, validExtractors, validMetadata}, `extract.crawler.exclude[0] must be a valid pattern, found: {a`},
		{&Extract{validSource, invalidCrawlerExcludeEmpty, validExtractors, validMetadata}, `extract.crawler.exclude[0] must be a valid pattern, found: `},
		{&Extract{validSource, validCrawler, invalidExtractorsEmpty, validMetadata}, `extract.extractors must contain at least one extractor, none found`},
		{&Extract{validSource, validCrawler, invalidExtractorsType, validMetadata}, `extract.extractors[0].type must be one of md, html, found: bogus`},
		{&Extract{validSource, validCrawler, invalidExtractorsInclude, validMetadata}, `extract.extractors[0].include[0] must be a valid pattern, found: {a`},
		{&Extract{validSource, validCrawler, invalidExtractorsIncludeEmpty, validMetadata}, `extract.extractors[0].include[0] must be a valid pattern, found: `},
		{&Extract{validSource, validCrawler, invalidExtractorsExclude, validMetadata}, `extract.extractors[0].exclude[0] must be a valid pattern, found: {a`},
		{&Extract{validSource, validCrawler, invalidExtractorsExcludeEmpty, validMetadata}, `extract.extractors[0].exclude[0] must be a valid pattern, found: `},
		{&Extract{validSource, validCrawler, validExtractors, invalidMetadataDocument}, `extract.metadata[0].document must be a valid pattern, found: {a`},
		{&Extract{validSource, validCrawler, validExtractors, invalidMetadataDocumentEmpty}, `extract.metadata[0].document must be a valid pattern, found: `},
		{&Extract{validSource, validCrawler, validExtractors, invalidMetadataSection}, `extract.metadata[0].section must be a valid pattern if not empty, found: {a`},
		{&Extract{validSource, validCrawler, validExtractors, validMetadataSectionEmpty}, ``},
		{&Extract{validSource, validCrawler, validExtractors, invalidMetadataTagKey}, `extract.metadata[0].tags[0].key must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: foo!`},
		{&Extract{validSource, validCrawler, validExtractors, invalidMetadataTagValue}, `extract.metadata[0].tags[0].value must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: bar!`},
	}

	for i, test := range tests {
		cfg := &Config{
			Extract: test.extract,
		}

		err := validateExtract(cfg)

		if test.err == "" && err != nil {
			t.Errorf("test %d - expected no error, got error: %s", i, err.Error())
		}
		if test.err != "" && err == nil {
			t.Errorf("test %d - expected error: %s, got no error", i, test.err)
		}
		if test.err != "" && err.Error() != test.err {
			t.Errorf("test %d - expected error: %s, got error: %s", i, test.err, err.Error())
		}
	}
}
