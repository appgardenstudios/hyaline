package config

import "testing"

func TestValidateCheck(t *testing.T) {
	validCode := CheckCode{
		Include: []string{"**/*"},
	}
	invalidCodeNoIncludes := CheckCode{}
	invalidCodeBlankInclude := CheckCode{
		Include: []string{""},
	}
	invalidCodeInvalidInclude := CheckCode{
		Include: []string{"{a"},
	}
	invalidCodeBlankExclude := CheckCode{
		Include: []string{"**/*"},
		Exclude: []string{""},
	}
	invalidCodeInvalidExclude := CheckCode{
		Include: []string{"**/*"},
		Exclude: []string{"{a"},
	}

	validDocumentation := CheckDocumentation{
		Include: []DocumentationFilter{
			{Source: "foo"},
		},
	}
	invalidDocumentationNoIncludes := CheckDocumentation{}
	invalidDocumentationInvalidInclude := CheckDocumentation{
		Include: []DocumentationFilter{
			{Source: ""},
		},
	}
	invalidDocumentationInvalidExclude := CheckDocumentation{
		Include: []DocumentationFilter{
			{Source: "foo"},
		},
		Exclude: []DocumentationFilter{
			{Source: ""},
		},
	}

	validOptions := CheckOptions{}
	invalidOptionsInvalidUpdateSource := CheckOptions{
		DetectDocumentationUpdates: CheckOptionsDetectDocumentationUpdates{
			Source: "**invalid**",
		},
	}
	invalidOptionsInvalidUpdateIfTouched := CheckOptions{
		UpdateIf: CheckOptionsUpdateIf{
			Touched: []CheckOptionsUpdateIfEntry{
				{},
			},
		},
	}
	invalidOptionsInvalidUpdateIfAdded := CheckOptions{
		UpdateIf: CheckOptionsUpdateIf{
			Added: []CheckOptionsUpdateIfEntry{
				{},
			},
		},
	}
	invalidOptionsInvalidUpdateIfModified := CheckOptions{
		UpdateIf: CheckOptionsUpdateIf{
			Modified: []CheckOptionsUpdateIfEntry{
				{},
			},
		},
	}
	invalidOptionsInvalidUpdateIfDeleted := CheckOptions{
		UpdateIf: CheckOptionsUpdateIf{
			Deleted: []CheckOptionsUpdateIfEntry{
				{},
			},
		},
	}
	invalidOptionsInvalidUpdateIfRenamed := CheckOptions{
		UpdateIf: CheckOptionsUpdateIf{
			Renamed: []CheckOptionsUpdateIfEntry{
				{},
			},
		},
	}

	var tests = []struct {
		check *Check
		err   string
	}{
		{nil, ``},
		{&Check{false, validCode, validDocumentation, validOptions}, ``},
		{&Check{Disabled: true}, ``},
		{&Check{false, invalidCodeNoIncludes, validDocumentation, validOptions}, `check.code.include must contain at least one entry, none found`},
		{&Check{false, invalidCodeBlankInclude, validDocumentation, validOptions}, `check.code.include[0] must be a valid pattern, found: `},
		{&Check{false, invalidCodeInvalidInclude, validDocumentation, validOptions}, `check.code.include[0] must be a valid pattern, found: {a`},
		{&Check{false, invalidCodeBlankExclude, validDocumentation, validOptions}, `check.code.exclude[0] must be a valid pattern, found: `},
		{&Check{false, invalidCodeInvalidExclude, validDocumentation, validOptions}, `check.code.exclude[0] must be a valid pattern, found: {a`},
		{&Check{false, validCode, invalidDocumentationNoIncludes, validOptions}, `check.documentation.include must contain at least one entry, none found`},
		{&Check{false, validCode, invalidDocumentationInvalidInclude, validOptions}, `check.documentation.include[0].source must be a valid pattern, found: `},
		{&Check{false, validCode, invalidDocumentationInvalidExclude, validOptions}, `check.documentation.exclude[0].source must be a valid pattern, found: `},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateSource}, `extract.options.detectDocumentationUpdates.source must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: **invalid**`},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateIfTouched}, `check.options.updateIf.touched[0].code.path must be a valid pattern, found: `},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateIfAdded}, `check.options.updateIf.added[0].code.path must be a valid pattern, found: `},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateIfModified}, `check.options.updateIf.modified[0].code.path must be a valid pattern, found: `},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateIfDeleted}, `check.options.updateIf.deleted[0].code.path must be a valid pattern, found: `},
		{&Check{false, validCode, validDocumentation, invalidOptionsInvalidUpdateIfRenamed}, `check.options.updateIf.renamed[0].code.path must be a valid pattern, found: `},
	}

	for i, test := range tests {
		cfg := &Config{
			Check: test.check,
		}

		err := ValidateCheck(cfg)

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

func TestValidateCheckUpdateIf(t *testing.T) {
	validCodeFilter := CheckCodeFilter{
		Path: "**/*",
	}
	invalidCodeFilter := CheckCodeFilter{
		Path: "",
	}

	validDocumentationFilter := DocumentationFilter{
		Source: "*",
	}
	invalidDocumentationFilter := DocumentationFilter{
		Source: "",
	}

	var tests = []struct {
		location string
		entries  []CheckOptionsUpdateIfEntry
		err      string
	}{
		{"location", []CheckOptionsUpdateIfEntry{}, ``},
		{"location", []CheckOptionsUpdateIfEntry{{validCodeFilter, validDocumentationFilter}}, ``},
		{"location", []CheckOptionsUpdateIfEntry{{invalidCodeFilter, validDocumentationFilter}}, `location[0].code.path must be a valid pattern, found: `},
		{"location", []CheckOptionsUpdateIfEntry{{validCodeFilter, invalidDocumentationFilter}}, `location[0].documentation.source must be a valid pattern, found: `},
	}

	for i, test := range tests {
		err := validateCheckUpdateIf(test.location, test.entries)

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

func TestValidateCheckCodeFilter(t *testing.T) {
	validFilter := CheckCodeFilter{
		Path: "foo",
	}
	invalidFilterBlankPath := CheckCodeFilter{
		Path: "",
	}
	invalidFilterInvalidPath := CheckCodeFilter{
		Path: "{a",
	}

	var tests = []struct {
		location string
		filter   CheckCodeFilter
		err      string
	}{
		{"location", validFilter, ``},
		{"location", invalidFilterBlankPath, `location.path must be a valid pattern, found: `},
		{"location", invalidFilterInvalidPath, `location.path must be a valid pattern, found: {a`},
	}

	for i, test := range tests {
		err := validateCheckCodeFilter(test.location, test.filter)

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

func TestValidateDocumentationFilter(t *testing.T) {
	validFilter := DocumentationFilter{
		Source: "foo",
	}
	invalidFilterUriInvalidPrefix := DocumentationFilter{
		URI: "invalid",
	}
	invalidFilterUriMissingSlash := DocumentationFilter{
		URI: "document://foo",
	}
	invalidFilterUriBlankSource := DocumentationFilter{
		URI: "document:///",
	}
	invalidFilterUriInvalidSource := DocumentationFilter{
		URI: "document://{a/",
	}
	invalidFilterUriBlankDocument := DocumentationFilter{
		URI: "document://foo/",
	}
	invalidFilterUriInvalidDocument := DocumentationFilter{
		URI: "document://foo/{a",
	}
	invalidFilterUriBlankSection := DocumentationFilter{
		URI: "document://foo/**/*#",
	}
	invalidFilterUriInvalidSection := DocumentationFilter{
		URI: "document://foo/bar#{a",
	}
	validFilterUri := DocumentationFilter{
		URI: "document://foo/**/*#**/*",
	}
	invalidFilterSourceBlank := DocumentationFilter{
		Source: "",
	}
	invalidFilterSourceInvalid := DocumentationFilter{
		Source: "{a",
	}
	invalidFilterDocumentInvalid := DocumentationFilter{
		Source:   "foo",
		Document: "{a",
	}
	invalidFilterSectionWithMissingDocument := DocumentationFilter{
		Source:  "foo",
		Section: "**/*",
	}
	invalidFilterSectionInvalid := DocumentationFilter{
		Source:   "foo",
		Document: "**/*",
		Section:  "{a",
	}
	validFilterTags := DocumentationFilter{
		Source: "foo",
		Tags: []DocumentationFilterTag{
			{"foo", "bar"},
		},
	}
	invalidFilterTagKey := DocumentationFilter{
		Source: "foo",
		Tags: []DocumentationFilterTag{
			{"**foo**", "bar"},
		},
	}
	invalidFilterTagValue := DocumentationFilter{
		Source: "foo",
		Tags: []DocumentationFilterTag{
			{"foo", ""},
		},
	}

	var tests = []struct {
		location string
		filter   DocumentationFilter
		err      string
	}{
		{"location", validFilter, ``},
		{"location", invalidFilterUriInvalidPrefix, `location.uri must start with document://, found: invalid`},
		{"location", invalidFilterUriMissingSlash, `location.uri must contain at least one /, found: document://foo`},
		{"location", invalidFilterUriBlankSource, `location.uri must contain a valid source pattern, found:  in document:///`},
		{"location", invalidFilterUriInvalidSource, `location.uri must contain a valid source pattern, found: {a in document://{a/`},
		{"location", invalidFilterUriBlankDocument, `location.uri must contain a valid document pattern, found:  in document://foo/`},
		{"location", invalidFilterUriInvalidDocument, `location.uri must contain a valid document pattern, found: {a in document://foo/{a`},
		{"location", invalidFilterUriBlankSection, `location.uri must contain a valid section pattern, found:  in document://foo/**/*#`},
		{"location", invalidFilterUriInvalidSection, `location.uri must contain a valid section pattern, found: {a in document://foo/bar#{a`},
		{"location", validFilterUri, ``},
		{"location", invalidFilterSourceBlank, `location.source must be a valid pattern, found: `},
		{"location", invalidFilterSourceInvalid, `location.source must be a valid pattern, found: {a`},
		{"location", invalidFilterDocumentInvalid, `location.document must be a valid pattern, found: {a`},
		{"location", invalidFilterSectionWithMissingDocument, `location.document must be set if location.section is set`},
		{"location", invalidFilterSectionInvalid, `location.section must be a valid pattern, found: {a`},
		{"location", validFilterTags, ``},
		{"location", invalidFilterTagKey, `location.tags[0].key must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: **foo**`},
		{"location", invalidFilterTagValue, `location.tags[0].value must match regex /^[A-z0-9][A-z0-9_-]{0,63}$/, found: `},
	}

	for i, test := range tests {
		err := validateDocumentationFilter(test.location, &test.filter)

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
