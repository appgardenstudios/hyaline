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
	Config            string
	Current           string
	System            string
	Output            string
	CheckPurpose      bool
	CheckCompleteness bool
}

type CheckCurrentOutput struct {
	Results []CheckCurrentOutputEntry `json:"results"`
}

type CheckCurrentOutputEntry struct {
	System              string   `json:"system"`
	DocumentationSource string   `json:"documentationSource"`
	Document            string   `json:"document"`
	Section             []string `json:"section,omitempty"`
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
		processedDocumentMap := make(map[string]struct{})
		processedSectionMap := make(map[string]struct{})

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

		// Loop through desiredDocuments
		for _, desiredDoc := range docSource.GetDocuments(cfg) {
			doc, found := docMap[desiredDoc.Name]
			// Check REQUIRED
			if desiredDoc.Required {
				result := "PASS"
				message := ""
				if !found {
					result = "ERROR"
					message = "This document is marked as required"
				}
				if desiredDoc.Ignore {
					result = "SKIPPED"
				}
				output.Results = append(output.Results, CheckCurrentOutputEntry{
					System:              system.ID,
					DocumentationSource: docSource.ID,
					Document:            desiredDoc.Name,
					Check:               "REQUIRED",
					Result:              result,
					Message:             message,
				})
			}

			// Check MATCHES_PURPOSE
			if args.CheckPurpose {
				result := "PASS"
				message := ""
				if desiredDoc.Ignore {
					result = "SKIPPED"
				} else if !found {
					result = "ERROR"
					message = "This document does not exist"
				} else {
					if desiredDoc.Purpose != "" {
						matches, reason, err := check.Purpose(system.ID, docSource.ID, desiredDoc.Name, []string{}, desiredDoc.Purpose, doc.ExtractedData, &cfg.LLM, currentDB)
						if err != nil {
							slog.Debug("action.CheckChange could not check purpose for document", "document", desiredDoc.Name, "documentationSource", docSource.ID, "system", args.System, "error", err)
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
					Document:            desiredDoc.Name,
					Check:               "MATCHES_PURPOSE",
					Result:              result,
					Message:             message,
				})
			}

			// Check COMPLETE
			if args.CheckCompleteness {
				result := "PASS"
				message := ""
				if desiredDoc.Ignore {
					result = "SKIPPED"
				} else if !found {
					result = "ERROR"
					message = "This document does not exist"
				} else {
					if desiredDoc.Purpose != "" {
						complete, reason, err := check.Completeness(system.ID, docSource.ID, desiredDoc.Name, []string{}, desiredDoc.Purpose, doc.ExtractedData, &cfg.LLM, currentDB)
						if err != nil {
							slog.Debug("action.CheckChange could not check completeness for document", "document", desiredDoc.Name, "documentationSource", docSource.ID, "system", args.System, "error", err)
							return err
						}
						if !complete {
							result = "ERROR"
							message = fmt.Sprintf("This document is not complete. %s", reason)
						} else {
							message = fmt.Sprintf("This document is complete. %s", reason)
						}
					} else {
						result = "WARN"
						message = "This document does not have a purpose"
					}
				}

				output.Results = append(output.Results, CheckCurrentOutputEntry{
					System:              system.ID,
					DocumentationSource: docSource.ID,
					Document:            desiredDoc.Name,
					Check:               "COMPLETE",
					Result:              result,
					Message:             message,
				})
			}

			// Check sections (if not skipped)
			if !desiredDoc.Ignore {
				// Get section map
				sectionMap := make(map[string]*sqlite.Section)
				sections, err := sqlite.GetAllSectionsForDocument(desiredDoc.Name, docSource.ID, system.ID, currentDB)
				if err != nil {
					slog.Debug("action.CheckChange could not get sections for document", "document", desiredDoc.Name, "documentationSource", docSource.ID, "system", args.System, "error", err)
					return err
				}
				for _, sec := range sections {
					sectionMap[sec.ID] = sec
				}

				// Check section
				addtlResults, err := checkCurrentSections(system.ID, docSource.ID, desiredDoc.Name, []string{}, desiredDoc.Sections, &sectionMap, &processedSectionMap, args.CheckPurpose, args.CheckCompleteness, &cfg.LLM, currentDB)
				if err != nil {
					slog.Debug("action.CheckChange could not check current sections for document", "document", desiredDoc.Name, "documentationSource", docSource.ID, "system", args.System, "error", err)
					return err
				}
				output.Results = append(output.Results, addtlResults...)
			}

			// Mark doc as processed
			processedDocumentMap[desiredDoc.Name] = struct{}{}
		}

		// Loop through docs and make sure each had a corresponding desired document
		for _, doc := range docs {
			result := "PASS"
			message := ""
			_, found := processedDocumentMap[doc.ID]
			if !found {
				result = "ERROR"
				message = fmt.Sprintf("Document %s does not have a corresponding desired document", doc.ID)
			}
			output.Results = append(output.Results, CheckCurrentOutputEntry{
				System:              system.ID,
				DocumentationSource: docSource.ID,
				Document:            doc.ID,
				Check:               "DESIRED_DOCUMENT_EXISTS",
				Result:              result,
				Message:             message,
			})
		}

		// Loop through sections and make sure each has a corresponding desired document
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
			_, found := processedSectionMap[sec.ID]
			if !found {
				result = "ERROR"
				message = fmt.Sprintf("Section %s does not have a corresponding desired document section", sec.ID)
			}
			output.Results = append(output.Results, CheckCurrentOutputEntry{
				System:              system.ID,
				DocumentationSource: docSource.ID,
				Document:            sec.ID,
				Section:             arr[1:], // Split document off of the ID and take what is left, e.g. doc#sec1#sec1.1
				Check:               "DESIRED_DOCUMENT_EXISTS",
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

func checkCurrentSections(system string, documentationSource string, document string, sectionArr []string, desiredDocSections []config.DocumentSection, sectionMap *map[string]*sqlite.Section, processedSectionMap *map[string]struct{}, checkPurpose bool, checkComplete bool, cfg *config.LLM, currentDB *sql.DB) (results []CheckCurrentOutputEntry, err error) {
	for _, desiredDocSection := range desiredDocSections {
		currentSection := []string{}
		currentSection = append(currentSection, sectionArr...)
		currentSection = append(currentSection, desiredDocSection.Name)
		sectionID := fmt.Sprintf("%s#%s", document, strings.Join(currentSection, "#"))

		section, found := (*sectionMap)[sectionID]

		// Check REQUIRED
		if desiredDocSection.Required {
			result := "PASS"
			message := ""
			if !found {
				result = "ERROR"
				message = "This section is marked as required"
			}
			if desiredDocSection.Ignore {
				result = "SKIPPED"
			}
			results = append(results, CheckCurrentOutputEntry{
				System:              system,
				DocumentationSource: documentationSource,
				Document:            document,
				Section:             currentSection,
				Check:               "REQUIRED",
				Result:              result,
				Message:             message,
			})
		}

		// Check MATCHES_PURPOSE
		if checkPurpose {
			result := "PASS"
			message := ""
			if desiredDocSection.Ignore {
				result = "SKIPPED"
			} else if !found {
				result = "ERROR"
				message = "This section does not exist"
			} else {
				if desiredDocSection.Purpose != "" {
					var matches bool
					var reason string
					matches, reason, err = check.Purpose(system, documentationSource, document, currentSection, desiredDocSection.Purpose, section.ExtractedData, cfg, currentDB)
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
				Check:               "MATCHES_PURPOSE",
				Result:              result,
				Message:             message,
			})
		}

		// Check COMPLETE
		if checkComplete {
			result := "PASS"
			message := ""
			if desiredDocSection.Ignore {
				result = "SKIPPED"
			} else if !found {
				result = "ERROR"
				message = "This section does not exist"
			} else {
				if desiredDocSection.Purpose != "" {
					var matches bool
					var reason string
					matches, reason, err = check.Completeness(system, documentationSource, document, currentSection, desiredDocSection.Purpose, section.ExtractedData, cfg, currentDB)
					if err != nil {
						return
					}
					if !matches {
						result = "ERROR"
						message = fmt.Sprintf("This section is not complete. %s", reason)
					} else {
						message = fmt.Sprintf("This section is complete. %s", reason)
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
				Check:               "COMPLETE",
				Result:              result,
				Message:             message,
			})
		}

		// Check sections (if not skipped)
		if !desiredDocSection.Ignore {
			var addtlResults []CheckCurrentOutputEntry
			addtlResults, err = checkCurrentSections(system, documentationSource, document, currentSection, desiredDocSection.Sections, sectionMap, processedSectionMap, checkPurpose, checkComplete, cfg, currentDB)
			if err != nil {
				return
			}
			results = append(results, addtlResults...)
		}

		// Mark section as processed
		(*processedSectionMap)[sectionID] = struct{}{}
	}

	return
}
