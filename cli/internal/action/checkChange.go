package action

import (
	"database/sql"
	"encoding/json"
	"errors"
	"hyaline/internal/check"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"hyaline/internal/suggest"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	_ "modernc.org/sqlite"
)

type CheckChangeArgs struct {
	Config  string
	Current string
	Change  string
	System  string
	Output  string
	Suggest bool
}

type CheckChangeResultKey struct {
	Documentation string
	Document      string
	Section       string
}

type CheckChangeOutput struct {
	Recommendations []CheckChangeOutputEntry `json:"recommendations"`
}

type CheckChangeOutputEntry struct {
	System              string                        `json:"system"`
	DocumentationSource string                        `json:"documentationSource"`
	Document            string                        `json:"document"`
	Section             []string                      `json:"section,omitempty"`
	Recommendation      string                        `json:"recommendation"`
	Reasons             []string                      `json:"reasons"`
	Changed             bool                          `json:"changed"`
	Suggestion          string                        `json:"suggestion,omitempty"`
	_References         []check.ChangeResultReference `json:"-"` // Always omit this as it is for internal purposes only
}

type CheckChangeOutputEntrySort []CheckChangeOutputEntry

func (c CheckChangeOutputEntrySort) Len() int {
	return len(c)
}
func (c CheckChangeOutputEntrySort) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}
func (c CheckChangeOutputEntrySort) Less(i, j int) bool {
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
	return strings.Join(c[i].Section, "#") < strings.Join(c[j].Section, "#")
}

func CheckChange(args *CheckChangeArgs) error {
	slog.Info("Checking changed code and docs")
	slog.Debug("action.CheckChange Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"change", args.Change,
		"system", args.System,
		"output", args.Output,
		"suggest", args.Suggest,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.CheckChange could not load the config", "error", err)
		return err
	}

	// Ensure output file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.CheckChange could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.CheckChange detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open Current DB
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.CheckChange could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.CheckChange could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.CheckChange opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Open Change DB
	var changeDB *sql.DB
	if args.Change != "" {
		changeAbsPath, err := filepath.Abs(args.Change)
		if err != nil {
			slog.Debug("action.CheckChange could not get an absolute path for change", "change", args.Change, "error", err)
			return err
		}
		changeDB, err = sql.Open("sqlite", changeAbsPath)
		if err != nil {
			slog.Debug("action.CheckChange could not open change SQLite DB", "dataSourceName", changeAbsPath, "error", err)
			return err
		}
		slog.Debug("action.CheckChange opened change database", "change", args.Change, "path", changeAbsPath)
		defer changeDB.Close()
	}

	// Get system
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.CheckChange could not locate the system", "system", args.System, "error", err)
		return err
	}

	// Get Changes
	systemChanges, err := sqlite.GetAllSystemChange(system.ID, changeDB)
	if err != nil {
		slog.Debug("action.CheckChange could not get related changes", "error", err)
		return err
	}

	// Get Tasks
	systemTasks, err := sqlite.GetAllSystemTask(system.ID, changeDB)
	if err != nil {
		slog.Debug("action.CheckChange could not get related tasks", "error", err)
		return err
	}

	// Initialize our output recommendations
	output := CheckChangeOutput{
		Recommendations: []CheckChangeOutputEntry{},
	}

	// Get the full set of desiredDocuments that apply to this system mapped by documentationSource
	desiredDocsMap := make(map[string][]config.Document)
	for _, docSource := range system.DocumentationSources {
		desiredDocsMap[docSource.ID] = docSource.GetDocuments(cfg)
	}

	// Get a map of the documents that have been updated as a part of this change (by documentation source)
	updatedDocumentMap := make(map[string][]*sqlite.SystemDocument)
	for _, doc := range system.DocumentationSources {
		documents, err := sqlite.GetAllSystemDocument(doc.ID, system.ID, changeDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get changed documents", "doc", doc.ID, "system", args.System, "error", err)
			return err
		}
		updatedDocumentMap[doc.ID] = documents
	}

	// Initialize our results map used to collect results across all code sources
	resultsMap := make(map[CheckChangeResultKey]*check.ChangeResult)

	// Get the set of documents/sections that need to be updated for each code change in each code source
	for _, c := range system.CodeSources {
		results := []check.ChangeResult{}

		// Get the set of files changed for this code source
		files, err := sqlite.GetAllSystemFiles(c.ID, system.ID, changeDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get files for codeSource", "codeSource", c.ID, "system", args.System, "error", err)
			return err
		}

		// Check each file against our full set of documentation
		for _, file := range files {
			arr, err := check.Change(file, c, desiredDocsMap, systemChanges, systemTasks, currentDB, changeDB, &cfg.LLM)
			results = append(results, arr...)
			if err != nil {
				slog.Debug("action.CheckChange could not check change", "file", file.ID, "system", args.System, "error", err)
				return err
			}
		}

		// Merge results into a master list
		for _, result := range results {
			key := CheckChangeResultKey{
				Documentation: result.DocumentationSource,
				Document:      result.Document,
				Section:       strings.Join(result.Section, "#"),
			}
			_, ok := resultsMap[key]
			if ok {
				resultsMap[key].Reasons = append(resultsMap[key].Reasons, result.Reasons...)
				resultsMap[key].References = append(resultsMap[key].References, result.References...)
			} else {
				resultsMap[key] = &result
			}
		}
	}

	// Create our entries
	for _, result := range resultsMap {
		// See if this document was already updated
		changed := false
		documents, ok := updatedDocumentMap[result.DocumentationSource]
		if ok {
			for _, document := range documents {
				if result.Document == document.ID {
					changed = true
					break
				}
			}
		}

		// Add the recommendation
		output.Recommendations = append(output.Recommendations, CheckChangeOutputEntry{
			System:              system.ID,
			DocumentationSource: result.DocumentationSource,
			Document:            result.Document,
			Section:             result.Section,
			Recommendation:      "Consider reviewing and updating this documentation",
			Reasons:             result.Reasons,
			Changed:             changed,
			_References:         result.References,
		})
	}

	// Sort the output list
	sort.Sort(CheckChangeOutputEntrySort(output.Recommendations))

	// Suggest change(s) (if flag is set)
	if args.Suggest {
		for idx, entry := range output.Recommendations {
			// Get purpose from desiredDoc
			purpose, _ := config.GetPurpose(entry.System, entry.DocumentationSource, entry.Document, entry.Section, cfg)
			suggestion, err := suggest.Change(entry.System, entry.DocumentationSource, entry.Document, entry.Section, purpose, entry.Reasons, entry._References, systemChanges, systemTasks, &cfg.LLM, currentDB)
			if err != nil {
				slog.Debug("action.CheckChange could not get suggestion",
					"system", entry.System,
					"doc", entry.DocumentationSource,
					"document", entry.Document,
					"section", entry.Section,
					"error", err)
				return err
			}
			if suggestion != "" {
				output.Recommendations[idx].Suggestion = suggestion
			} else {
				output.Recommendations[idx].Suggestion = "(none)"
			}
		}
	}

	// Output the results
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
