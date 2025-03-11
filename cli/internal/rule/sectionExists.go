package rule

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/recommend"
	"hyaline/internal/sqlite"
	"log/slog"
	"strings"
)

var SectionExistsRule = "section-exists"

type SectionExistsOptions struct {
	Document   string `json:"document"`
	Section    string `json:"section"`
	Severity   string `json:"severity"`
	AllowTodos bool   `json:"allowTodos"`
}

func GetSectionExistsOptions(raw map[string]interface{}) (options SectionExistsOptions, err error) {
	// Get options w/ defaults
	options = SectionExistsOptions{
		Document:   getString("document", raw, ""),
		Section:    getString("section", raw, ""),
		Severity:   getString("severity", raw, "WARNING"),
		AllowTodos: getBool("allowTodos", raw, true),
	}

	// Validate
	if options.Section == "" {
		err = errors.New("options.section is required")
	}
	if options.Document == "" {
		err = errors.New("options.document is required")
	}

	return
}

func sectionExistsWithoutTodos(allowTodos bool) string {
	if allowTodos {
		return ""
	} else {
		return " without TODOs"
	}
}

func RunSectionExists(id string, description string, options SectionExistsOptions, system string, current *sql.DB, recommendAction bool, llmOpts config.LLM) (result *Result, err error) {
	result = &Result{
		System:      system,
		ID:          id,
		Description: description,
		Rule:        SectionExistsRule,
		Options:     options,
	}

	// Retrieve section (if exists)
	section, err := sqlite.GetDocumentSection(options.Document, options.Section, system, current)
	if err != nil {
		return
	}

	// Ensure section exists
	if section == nil {
		result.Pass = false
		result.Severity = options.Severity
		result.Message = fmt.Sprintf("The section '%s' must exist in '%s'%s.", options.Section, options.Document, sectionExistsWithoutTodos(options.AllowTodos))
		if recommendAction {
			action, err := recommend.SectionExists(true, options.Section, options.Document, system, current, llmOpts)
			if err != nil {
				slog.Debug("RunSectionExists could not generate recommendation", "error", err)
			} else {
				result.Action = action
			}
		}
		return
	}

	// Ensure section is not empty
	if strings.TrimSpace(section.RawData) == "" {
		result.Pass = false
		result.Severity = options.Severity
		result.Message = fmt.Sprintf("The section '%s' in '%s' must contain text%s.", options.Section, options.Document, sectionExistsWithoutTodos(options.AllowTodos))
		if recommendAction {
			action, err := recommend.SectionExists(false, options.Section, options.Document, system, current, llmOpts)
			if err != nil {
				slog.Debug("RunSectionExists could not generate recommendation", "error", err)
			} else {
				result.Action = action
			}
		}
		return
	}

	// If allowTodos is false, ensure there are no TODOs
	if !options.AllowTodos && strings.Contains(section.RawData, "TODO") {
		result.Pass = false
		result.Severity = options.Severity
		result.Message = fmt.Sprintf("The section '%s' in '%s' must not contain TODOs.", options.Section, options.Document)
		if recommendAction {
			result.Action = fmt.Sprintf("You should resolve the TODOs remaining in the section '%s' in '%s'.", options.Section, options.Document)
		}
		return
	}

	// If we are here, pass
	result.Pass = true

	return
}
