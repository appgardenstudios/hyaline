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

	_ "modernc.org/sqlite"
)

type CheckChangeArgs struct {
	Config  string
	Current string
	Change  string
	System  string
	Output  string
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
	System              string   `json:"system"`
	DocumentationSource string   `json:"documentationSource"`
	Document            string   `json:"document"`
	Section             string   `json:"section,omitempty"`
	Recommendation      string   `json:"recommendation"`
	Reasons             []string `json:"reasons"`
	Changed             bool     `json:"changed"`
}

func CheckChange(args *CheckChangeArgs) error {
	slog.Info("Checking changed code and docs")
	slog.Debug("action.CheckChange Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"change", args.Change,
		"system", args.System,
		"output", args.Output,
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

	// Initialize our output
	output := CheckChangeOutput{
		Recommendations: []CheckChangeOutputEntry{},
	}

	// Get the full set of ruleDocuments that apply to this system mapped by documentationSource
	ruleDocsMap := make(map[string][]config.RuleDocument)
	for _, doc := range system.DocumentationSources {
		ruleDocs := []config.RuleDocument{}
		for _, ruleID := range doc.Rules {
			rules := config.GetRule(cfg.Rules, ruleID)
			ruleDocs = append(ruleDocs, rules.Documents...)
		}
		ruleDocsMap[doc.ID] = ruleDocs
	}

	// Get a map of the documents that have been updated as a part of this change (by documentation source)
	updatedDocumentMap := make(map[string][]*sqlite.Document)
	for _, doc := range system.DocumentationSources {
		documents, err := sqlite.GetAllDocument(doc.ID, system.ID, changeDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get changed documents", "doc", doc.ID, "system", args.System, "error", err)
			return err
		}
		updatedDocumentMap[doc.ID] = documents
	}

	// Initalize our results map used to collect results across all code sources
	resultsMap := make(map[CheckChangeResultKey]*check.ChangeResult)

	// Get the set of documents/sections that need to be updated for each code change in each code source
	for _, c := range system.CodeSources {
		results := []check.ChangeResult{}

		// Get the set of files changed for this code source
		files, err := sqlite.GetAllFiles(c.ID, system.ID, changeDB)
		if err != nil {
			slog.Debug("action.CheckChange could not get files for codeSource", "codeSource", c.ID, "system", args.System, "error", err)
			return err
		}

		// Check each file against our full set of documentation
		for _, file := range files {
			arr, err := check.Change(file, c, ruleDocsMap)
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
				Section:       result.Section,
			}
			_, ok := resultsMap[key]
			if ok {
				resultsMap[key].Reasons = append(resultsMap[key].Reasons, result.Reasons...)
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

		output.Recommendations = append(output.Recommendations, CheckChangeOutputEntry{
			System:              system.ID,
			DocumentationSource: result.DocumentationSource,
			Document:            result.Document,
			Section:             result.Section,
			Changed:             changed,
			Reasons:             result.Reasons,
		})
	}

	// Sort the output list
	// TODO

	// Output the results
	jsonData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		slog.Debug("action.CheckChange could not marshal json", "error", err)
		return err
	}

	// TODO output to output file
	fmt.Println(string(jsonData))

	return nil
}
