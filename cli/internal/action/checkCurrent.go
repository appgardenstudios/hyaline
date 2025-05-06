package action

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"hyaline/internal/check"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type CheckCurrentArgs struct {
	Config       string
	Current      string
	System       string
	Output       string
	CheckPurpose bool
}

type CheckCurrentOutput struct {
	Results []CheckCurrentOutputEntry `json:"results"`
}

type CheckCurrentOutputEntry struct {
	System              string   `json:"system"`
	DocumentationSource string   `json:"documentationSource"`
	Document            string   `json:"document"`
	Section             []string `json:"section,omitempty"`
	Rule                string   `json:"rule"`
	Check               string   `json:"check"`
	Result              string   `json:"result"`
	Message             string   `json:"message"`
}

type CheckCurrentOutputEntrySort []CheckCurrentOutputEntry

func (c CheckCurrentOutputEntrySort) Len() int {
	return len(c)
}
func (c CheckCurrentOutputEntrySort) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CheckCurrentOutputEntrySort) Less(i, j int) bool {
	if c[i].System < c[j].System {
		return true
	}
	if c[i].System > c[j].System {
		return false
	}
	if c[i].DocumentationSource < c[j].DocumentationSource {
		return true
	}
	if c[i].DocumentationSource > c[j].DocumentationSource {
		return false
	}
	if c[i].Document < c[j].Document {
		return true
	}
	if c[i].Document > c[j].Document {
		return false
	}
	if strings.Join(c[i].Section, "#") < strings.Join(c[j].Section, "#") {
		return true
	}
	if strings.Join(c[i].Section, "#") > strings.Join(c[j].Section, "#") {
		return false
	}
	if c[i].Rule < c[j].Rule {
		return true
	}
	if c[i].Rule > c[j].Rule {
		return false
	}
	return c[i].Check < c[j].Check
}

func CheckCurrent(args *CheckCurrentArgs) error {
	slog.Info("Checking current docs")
	slog.Debug("action.CheckCurrent Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
		"output", args.Output,
		"check-purpose", args.CheckPurpose,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckCurrent could not load the config", "error", err)
		return err
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckCurrent could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckCurrent detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open Current DB
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.CheckCurrent could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.CheckCurrent could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.CheckCurrent opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Get system
	system, found := cfg.GetSystem(args.System)
	if !found {
		err = fmt.Errorf("system not found: %s", args.System)
		slog.Debug("action.CheckChange could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Initialize our output
	output := CheckCurrentOutput{
		Results: []CheckCurrentOutputEntry{},
	}

	// Process each documentation source in the system
	for _, docSource := range system.DocumentationSources {
		// Initialize our processed documents/section map
		processedDocumentMap := make(map[string]string)
		processedSectionMap := make(map[string]string)

		// Get all documents for system and put them into a map
		docMap := make(map[string]*sqlite.Document)
		docs, err := sqlite.GetAllDocument(docSource.ID, system.ID, currentDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get documents for documentationSource", "documentationSource", docSource.ID, "system", args.System, "error", err)
			return err
		}
		for _, doc := range docs {
			docMap[doc.ID] = doc
		}

		// Loop through docRules
		for _, ruleID := range docSource.Rules {
			// Get rule set
			ruleSet, _ := cfg.GetRuleSet(ruleID)

			// Loop through documents
			for _, ruleDoc := range ruleSet.Documents {
				// Check REQUIRED
				if ruleDoc.Required {
					result := "PASS"
					message := ""
					_, found := docMap[ruleDoc.Path]
					if !found {
						result = "ERROR"
						message = "This document is marked as required"
					}
					if ruleDoc.Ignore {
						result = "SKIPPED"
					}
					output.Results = append(output.Results, CheckCurrentOutputEntry{
						System:              system.ID,
						DocumentationSource: docSource.ID,
						Document:            ruleDoc.Path,
						Rule:                ruleID,
						Check:               "REQUIRED",
						Result:              result,
						Message:             message,
					})
				}

				// Check MATCHES_PURPOSE
				if args.CheckPurpose {
					result := "PASS"
					message := ""
					if ruleDoc.Ignore {
						result = "SKIPPED"
					} else {
						if ruleDoc.Purpose != "" {
							matches, reason, err := check.Purpose()
							if err != nil {
								slog.Debug("action.CheckChange could not check purpose for document", "document", ruleDoc.Path, "documentationSource", docSource.ID, "system", args.System, "error", err)
								return err
							}
							if !matches {
								result = "ERROR"
								message = fmt.Sprintf("This document does not match it's purpose. %s", reason)
							} else {
								message = fmt.Sprintf("This document does match it's purpose. %s", reason)
							}
						} else {
							result = "WARN"
							message = "This document does not have a purpose"
						}
					}

					output.Results = append(output.Results, CheckCurrentOutputEntry{
						System:              system.ID,
						DocumentationSource: docSource.ID,
						Document:            ruleDoc.Path,
						Rule:                ruleID,
						Check:               "MATCHES_PURPOSE",
						Result:              result,
						Message:             message,
					})
				}

				// Check sections (if not skipped)
				if !ruleDoc.Ignore {
					// Get section map
					sectionMap := make(map[string]*sqlite.Section)
					sections, err := sqlite.GetAllSectionsForDocument(ruleDoc.Path, docSource.ID, system.ID, currentDB)
					if err != nil {
						slog.Debug("action.CheckChange could not get sections for document", "document", ruleDoc.Path, "documentationSource", docSource.ID, "system", args.System, "error", err)
						return err
					}
					for _, sec := range sections {
						sectionMap[sec.ID] = sec
					}

					// Check section
					addtlResults, err := checkCurrentSections(ruleID, system.ID, docSource.ID, ruleDoc.Path, []string{}, ruleDoc.Sections, &sectionMap, &processedSectionMap, args.CheckPurpose)
					if err != nil {
						slog.Debug("action.CheckChange could not check current sections for document", "document", ruleDoc.Path, "documentationSource", docSource.ID, "system", args.System, "error", err)
						return err
					}
					output.Results = append(output.Results, addtlResults...)
				}

				// Mark doc as processed
				processedDocumentMap[ruleDoc.Path] = ruleID
			}
		}

		// Loop through docs and make sure each had a corresponding rule
		for _, doc := range docs {
			result := "PASS"
			message := ""
			correspondingRuleID, found := processedDocumentMap[doc.ID]
			if !found {
				result = "ERROR"
				message = fmt.Sprintf("Document %s does not have a corresponding rule", doc.ID)
			}
			output.Results = append(output.Results, CheckCurrentOutputEntry{
				System:              system.ID,
				DocumentationSource: docSource.ID,
				Document:            doc.ID,
				Rule:                correspondingRuleID,
				Check:               "RULE_EXISTS",
				Result:              result,
				Message:             message,
			})
		}

		// Loop through sections and make sure each has a corresponding rule
		allSections, err := sqlite.GetAllSection(docSource.ID, system.ID, currentDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get sections for documentationSource", "documentationSource", docSource.ID, "system", args.System, "error", err)
			return err
		}
		for _, sec := range allSections {
			arr := strings.Split(sec.ID, "#")
			// Skip root sections
			if len(arr) < 2 {
				continue
			}

			result := "PASS"
			message := ""
			correspondingRuleID, found := processedSectionMap[sec.ID]
			if !found {
				result = "ERROR"
				message = fmt.Sprintf("Section %s does not have a corresponding rule", sec.ID)
			}
			output.Results = append(output.Results, CheckCurrentOutputEntry{
				System:              system.ID,
				DocumentationSource: docSource.ID,
				Document:            sec.ID,
				Section:             arr[1:], // Split document off of the ID and take what is left, e.g. doc#sec1#sec1.1
				Rule:                correspondingRuleID,
				Check:               "RULE_EXISTS",
				Result:              result,
				Message:             message,
			})
		}
	}

	// Sort output
	sort.Sort(CheckCurrentOutputEntrySort(output.Results))

	// Output
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		slog.Debug("action.CheckChange could not marshal json", "error", err)
		return err
	}
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.CheckChange could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()
	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.GenerateConfig could not write output file", "error", err)
		return err
	}

	return nil
}

func checkCurrentSections(ruleID string, system string, documentationSource string, document string, section []string, ruleDocSections []config.RuleDocumentSection, sectionMap *map[string]*sqlite.Section, processedSectionMap *map[string]string, checkPurpose bool) (results []CheckCurrentOutputEntry, err error) {
	for _, ruleDocSection := range ruleDocSections {
		currentSection := []string{}
		currentSection = append(currentSection, section...)
		currentSection = append(currentSection, ruleDocSection.ID)
		sectionID := fmt.Sprintf("%s#%s", document, strings.Join(currentSection, "#"))

		// Check REQUIRED
		if ruleDocSection.Required {
			result := "PASS"
			message := ""
			_, found := (*sectionMap)[sectionID]
			if !found {
				result = "ERROR"
				message = "This section is marked as required"
			}
			if ruleDocSection.Ignore {
				result = "SKIPPED"
			}
			results = append(results, CheckCurrentOutputEntry{
				System:              system,
				DocumentationSource: documentationSource,
				Document:            document,
				Section:             currentSection,
				Rule:                ruleID,
				Check:               "REQUIRED",
				Result:              result,
				Message:             message,
			})
		}

		// Check MATCHES_PURPOSE
		if checkPurpose {
			result := "PASS"
			message := ""
			if ruleDocSection.Ignore {
				result = "SKIPPED"
			} else {
				if ruleDocSection.Purpose != "" {
					var matches bool
					var reason string
					matches, reason, err = check.Purpose()
					if err != nil {
						return
					}
					if !matches {
						result = "ERROR"
						message = fmt.Sprintf("This section does not match it's purpose. %s", reason)
					} else {
						message = fmt.Sprintf("This section does match it's purpose. %s", reason)
					}
				} else {
					result = "WARN"
					message = "This section does not have a purpose"
				}
			}

			results = append(results, CheckCurrentOutputEntry{
				System:              system,
				DocumentationSource: documentationSource,
				Document:            document,
				Section:             currentSection,
				Rule:                ruleID,
				Check:               "MATCHES_PURPOSE",
				Result:              result,
				Message:             message,
			})
		}

		// Check sections (if not skipped)
		if !ruleDocSection.Ignore {
			var addtlResults []CheckCurrentOutputEntry
			addtlResults, err = checkCurrentSections(ruleID, system, documentationSource, document, currentSection, ruleDocSection.Sections, sectionMap, processedSectionMap, checkPurpose)
			if err != nil {
				return
			}
			results = append(results, addtlResults...)
		}

		// Mark section as processed
		(*processedSectionMap)[sectionID] = ruleID
	}

	return
}
