package audit

import (
	"context"
	"fmt"
	"hyaline/internal/audit/checks"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/sqlite"
	"log/slog"
	"sort"
	"strings"
)

const (
	CheckContentExists        = "CONTENT_EXISTS"
	CheckContentMinLength     = "CONTENT_MIN_LENGTH"
	CheckContentMatchesRegex  = "CONTENT_MATCHES_REGEX"
	CheckContentMatchesPrompt = "CONTENT_MATCHES_PROMPT"
	CheckContentMatchesPurpose = "CONTENT_MATCHES_PURPOSE"
	CheckPurposeExists        = "PURPOSE_EXISTS"
	CheckTagsContains         = "TAGS_CONTAINS"
)

// AuditOutput represents the top-level audit results
type AuditOutput struct {
	Results []AuditRuleResult `json:"results"`
}

// AuditRuleResult represents the result of a single audit rule
type AuditRuleResult struct {
	Rule        string             `json:"rule"`
	Description string             `json:"description"`
	Pass        bool               `json:"pass"`
	Checks      []AuditCheckResult `json:"checks"`
}

// AuditCheckResult represents the result of a single audit check
type AuditCheckResult struct {
	Source   string   `json:"source"`
	Document string   `json:"document"`
	Section  []string `json:"section,omitempty"`
	URI      string   `json:"uri"`
	Rule     string   `json:"rule"`
	Check    string   `json:"check"`
	Pass     bool     `json:"pass"`
	Message  string   `json:"message"`
}

// Documentation executes the audit process against the provided database
func Documentation(cfg *config.Config, db *sqlite.Queries, sources []string) (*AuditOutput, error) {
	slog.Debug("audit.Documentation starting")

	// Load all data from database
	documents, err := db.GetAllDocuments(context.Background())
	if err != nil {
		slog.Debug("audit.Documentation could not get all documents", "error", err)
		return nil, err
	}

	documentTags, err := db.GetAllDocumentTags(context.Background())
	if err != nil {
		slog.Debug("audit.Documentation could not get all document tags", "error", err)
		return nil, err
	}

	sections, err := db.GetAllSections(context.Background())
	if err != nil {
		slog.Debug("audit.Documentation could not get all sections", "error", err)
		return nil, err
	}

	sectionTags, err := db.GetAllSectionTags(context.Background())
	if err != nil {
		slog.Debug("audit.Documentation could not get all section tags", "error", err)
		return nil, err
	}

	// Filter data by sources if specified
	if len(sources) > 0 {
		documents = filterBySource(documents, sources, getDocumentSourceID)
		documentTags = filterBySource(documentTags, sources, getDocumentTagSourceID)
		sections = filterBySource(sections, sources, getSectionSourceID)
		sectionTags = filterBySource(sectionTags, sources, getSectionTagSourceID)
	}

	// Build maps for efficient lookup
	documentTagMap := docs.GetDocumentTagMap(documentTags)
	sectionTagMap := docs.GetSectionTagMap(sectionTags)

	output := &AuditOutput{
		Results: []AuditRuleResult{},
	}

	// Process each rule
	for _, rule := range cfg.Audit.Rules {
		slog.Debug("audit.Documentation processing rule", "ruleID", rule.ID)

		ruleResult := AuditRuleResult{
			Rule:        rule.ID,
			Description: rule.Description,
			Checks:      []AuditCheckResult{},
		}

		// Process the rule
		err := processRule(&rule, documents, documentTagMap, sections, sectionTagMap, &ruleResult, cfg)
		if err != nil {
			slog.Debug("audit.Documentation error processing rule", "ruleID", rule.ID, "error", err)
			return nil, err
		}

		// Calculate rule pass status based on checks
		ruleResult.Pass = true
		for _, check := range ruleResult.Checks {
			if !check.Pass {
				ruleResult.Pass = false
				break
			}
		}

		// Sort checks within each rule
		sort.Slice(ruleResult.Checks, func(i, j int) bool {
			a, b := ruleResult.Checks[i], ruleResult.Checks[j]

			// Sort by URI first
			if a.URI != b.URI {
				return a.URI < b.URI
			}

			// Sort by check type
			return a.Check < b.Check
		})

		output.Results = append(output.Results, ruleResult)
	}

	// Sort rules by rule ID
	sort.Slice(output.Results, func(i, j int) bool {
		return output.Results[i].Rule < output.Results[j].Rule
	})

	slog.Debug("audit.Documentation completed")
	return output, nil
}

func processRule(rule *config.AuditRule, documents []sqlite.DOCUMENT, documentTagMap map[string][]docs.FilteredTag, sections []sqlite.SECTION, sectionTagMap map[string][]docs.FilteredTag, ruleResult *AuditRuleResult, cfg *config.Config) error {
	// Track if we found any matches for CONTENT_EXISTS check
	var firstMatchSource, firstMatchDocument string
	var firstMatchSection []string
	var firstMatchURI string
	foundMatch := false

	// Process documents
	for _, document := range documents {
		documentKey := document.SourceID + "/" + document.ID
		documentTags := documentTagMap[documentKey]

		// Check if document matches any documentation filter and not any ignore filter
		if matchesAnyFilter(document.ID, document.SourceID, "", documentTags, rule.Documentation) &&
			!matchesAnyFilter(document.ID, document.SourceID, "", documentTags, rule.Ignore) {

			// Track first match for CONTENT_EXISTS
			if !foundMatch {
				foundMatch = true
				firstMatchSource = document.SourceID
				firstMatchDocument = document.ID
				firstMatchSection = nil
				firstMatchURI = (&docs.DocumentURI{
					SourceID:     document.SourceID,
					DocumentPath: document.ID,
				}).String()
			}

			// Perform document-level checks (except CONTENT_EXISTS)
			baseResult := AuditCheckResult{
				Source:   document.SourceID,
				Document: document.ID,
				URI: (&docs.DocumentURI{
					SourceID:     document.SourceID,
					DocumentPath: document.ID,
				}).String(),
				Rule: rule.ID,
			}
			performContentChecks(rule, baseResult, document.SourceID, document.ID, "", document.ExtractedData, document.Purpose, ruleResult, cfg)
			performPurposeChecks(rule, baseResult, document.Purpose, ruleResult)
			performTagsChecks(rule, baseResult, documentTags, ruleResult)
		}
	}

	// Process sections
	for _, section := range sections {
		sectionKey := section.SourceID + "/" + section.DocumentID + "#" + section.ID
		sectionTags := sectionTagMap[sectionKey]

		// Check if section matches any documentation filter and not any ignore filter
		if matchesAnyFilter(section.DocumentID, section.SourceID, section.ID, sectionTags, rule.Documentation) &&
			!matchesAnyFilter(section.DocumentID, section.SourceID, section.ID, sectionTags, rule.Ignore) {

			// Track first match for CONTENT_EXISTS
			if !foundMatch {
				foundMatch = true
				firstMatchSource = section.SourceID
				firstMatchDocument = section.DocumentID
				firstMatchSection = strings.Split(section.ID, "/") // Assuming section ID is path-like
				firstMatchURI = (&docs.DocumentURI{
					SourceID:     section.SourceID,
					DocumentPath: section.DocumentID,
					Section:      section.ID,
				}).String()
			}

			// Perform section-level checks (except CONTENT_EXISTS)
			baseResult := AuditCheckResult{
				Source:   section.SourceID,
				Document: section.DocumentID,
				Section:  strings.Split(section.ID, "/"),
				URI: (&docs.DocumentURI{
					SourceID:     section.SourceID,
					DocumentPath: section.DocumentID,
					Section:      section.ID,
				}).String(),
				Rule: rule.ID,
			}
			performContentChecks(rule, baseResult, section.SourceID, section.DocumentID, section.ID, section.ExtractedData, section.Purpose, ruleResult, cfg)
			performPurposeChecks(rule, baseResult, section.Purpose, ruleResult)
			performTagsChecks(rule, baseResult, sectionTags, ruleResult)
		}
	}

	// Handle CONTENT_EXISTS check
	if rule.Checks.Content.Exists {
		checkResult := AuditCheckResult{
			Rule:  rule.ID,
			Check: CheckContentExists,
			Pass:  foundMatch,
		}

		if foundMatch {
			checkResult.Source = firstMatchSource
			checkResult.Document = firstMatchDocument
			checkResult.Section = firstMatchSection
			checkResult.URI = firstMatchURI
		} else {
			checkResult.Message = "This content does not exist."
		}

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}

	return nil
}

func performContentChecks(rule *config.AuditRule, baseResult AuditCheckResult, sourceID, documentID, sectionID, content, purpose string, ruleResult *AuditRuleResult, cfg *config.Config) {
	// CONTENT_MIN_LENGTH check
	if rule.Checks.Content.MinLength > 0 {
		pass, message := checks.ContentMinLength(content, rule.Checks.Content.MinLength)

		checkResult := baseResult
		checkResult.Check = CheckContentMinLength
		checkResult.Pass = pass
		checkResult.Message = message

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}

	// CONTENT_MATCHES_REGEX check
	if rule.Checks.Content.MatchesRegex != "" {
		pass, message := checks.ContentMatchesRegex(content, rule.Checks.Content.MatchesRegex)

		checkResult := baseResult
		checkResult.Check = CheckContentMatchesRegex
		checkResult.Pass = pass
		checkResult.Message = message

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}

	// CONTENT_MATCHES_PROMPT check
	if rule.Checks.Content.MatchesPrompt != "" {
		pass, message, err := checks.ContentMatchesPrompt(sourceID, documentID, sectionID, rule.Checks.Content.MatchesPrompt, content, &cfg.LLM)
		if err != nil {
			slog.Debug("audit.performContentChecks error in CONTENT_MATCHES_PROMPT", "error", err)
			pass = false
			message = fmt.Sprintf("Error checking prompt: %v", err)
		}

		checkResult := baseResult
		checkResult.Check = CheckContentMatchesPrompt
		checkResult.Pass = pass
		checkResult.Message = message

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}

	// CONTENT_MATCHES_PURPOSE check
	if rule.Checks.Content.MatchesPurpose {
		if purpose != "" {
			pass, message, err := checks.ContentMatchesPurpose(sourceID, documentID, sectionID, purpose, content, &cfg.LLM)
			if err != nil {
				slog.Debug("audit.performContentChecks error in CONTENT_MATCHES_PURPOSE", "error", err)
				pass = false
				message = fmt.Sprintf("Error checking purpose: %v", err)
			}

			checkResult := baseResult
			checkResult.Check = CheckContentMatchesPurpose
			checkResult.Pass = pass
			checkResult.Message = message

			ruleResult.Checks = append(ruleResult.Checks, checkResult)
		}
	}
}

func performPurposeChecks(rule *config.AuditRule, baseResult AuditCheckResult, purpose string, ruleResult *AuditRuleResult) {
	// PURPOSE_EXISTS check
	if rule.Checks.Purpose.Exists {
		pass, message := checks.PurposeExists(purpose)

		checkResult := baseResult
		checkResult.Check = CheckPurposeExists
		checkResult.Pass = pass
		if !pass {
			checkResult.Message = message
		}

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}

}

func performTagsChecks(rule *config.AuditRule, baseResult AuditCheckResult, tags []docs.FilteredTag, ruleResult *AuditRuleResult) {
	// TAGS_CONTAINS check
	if len(rule.Checks.Tags.Contains) > 0 {
		pass, message := checks.TagsContains(tags, rule.Checks.Tags.Contains)

		checkResult := baseResult
		checkResult.Check = CheckTagsContains
		checkResult.Pass = pass
		checkResult.Message = message

		ruleResult.Checks = append(ruleResult.Checks, checkResult)
	}
}

func matchesAnyFilter(documentID, sourceID, sectionID string, tags []docs.FilteredTag, filters []config.DocumentationFilter) bool {
	for _, filter := range filters {
		if sectionID != "" {
			// Check section match
			if docs.SectionMatches(sectionID, documentID, sourceID, tags, &filter, true) {
				return true
			}
		} else {
			// Check document match
			if docs.DocumentMatches(documentID, sourceID, tags, &filter) {
				return true
			}
		}
	}
	return false
}

// Generic helper function to filter any slice by source
func filterBySource[T any](items []T, sources []string, getSourceID func(T) string) []T {
	var filtered []T
	sourceMap := make(map[string]bool)
	for _, source := range sources {
		sourceMap[source] = true
	}

	for _, item := range items {
		if sourceMap[getSourceID(item)] {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// Helper functions to extract SourceID from each type
func getDocumentSourceID(d sqlite.DOCUMENT) string        { return d.SourceID }
func getDocumentTagSourceID(dt sqlite.DOCUMENTTAG) string { return dt.SourceID }
func getSectionSourceID(s sqlite.SECTION) string          { return s.SourceID }
func getSectionTagSourceID(st sqlite.SECTIONTAG) string   { return st.SourceID }
