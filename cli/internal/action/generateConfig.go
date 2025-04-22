package action

import (
	"database/sql"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"path/filepath"

	"gopkg.in/yaml.v3"

	_ "modernc.org/sqlite"
)

type GenerateConfigArgs struct {
	Config         string
	Current        string
	System         string
	Output         string
	IncludePurpose bool
}

func GenerateConfig(args *GenerateConfigArgs) error {
	slog.Info("Generating Config")
	slog.Debug("action.GenerateConfig Args", slog.Group("args",
		"config", args.Config,
		"current", args.Current,
		"system", args.System,
		"output", args.Output,
		"include-purpose", args.IncludePurpose,
	))

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.GenerateConfig could not load the config", "error", err)
		return err
	}

	// Ensure output location does not exist
	// TODO

	// Open current db
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.Check could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.Check could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.Check opened current database", "current", args.Current, "path", currentAbsPath)
	defer currentDB.Close()

	// Get System
	system, err := config.GetSystem(args.System, cfg)
	if err != nil {
		slog.Debug("action.GenerateConfig could not locate the system", "system", args.System, "error", err)
		return err
	}

	// New config
	newCfg := config.Config{}

	// Loop through docs in our current system and generate a config for each
	for _, d := range system.Docs {
		// Create a rule for this documentation
		rule := config.Rule{
			ID:        d.ID,
			Documents: []config.RuleDocument{},
		}

		// Get a list of Documents from the db for this doc ID
		documents, err := sqlite.GetAllDocument(d.ID, system.ID, currentDB)
		if err != nil {
			slog.Debug("action.GenerateConfig could not get documents from current db", "doc", d.ID, "system", system.ID, "error", err)
			return err
		}

		// Loop through each document to generate rules for it
		for _, doc := range documents {
			// Create rule for the document
			ruleDoc := config.RuleDocument{
				Path:     doc.ID,
				Required: true,
			}

			// If IncludePurpose flag is set, generate purpose
			// TODO

			// Get and add sections for this document
			sections, err := sqlite.GetAllSectionsForDocument(doc.ID, d.ID, system.ID, currentDB)
			if err != nil {
				slog.Debug("action.GenerateConfig could not get sections for a document from current db", "document", doc.ID, "doc", d.ID, "system", system.ID, "error", err)
				return err
			}
			ruleDoc.Sections = append(ruleDoc.Sections, createRuleSections(sections, "")...)

			// Add ruleDoc to rule
			rule.Documents = append(rule.Documents, ruleDoc)
		}

		// Add rule to config
		newCfg.Rules = append(newCfg.Rules, rule)
	}

	// Output new config
	fmt.Println("system", system)
	fmt.Println("rules", newCfg.Rules)

	yml, err := yaml.Marshal(newCfg)
	if err != nil {
		slog.Debug("action.GenerateConfig could not marshal yaml", "error", err)
		return err
	}
	fmt.Println("config", string(yml))

	return nil
}

func createRuleSections(sections []*sqlite.Section, parentID string) []config.RuleDocumentSection {
	docSections := []config.RuleDocumentSection{}

	for _, section := range sections {
		// TODO guard against circular issues by ensuring that no ID is the same as its parent ID

		if section.ParentID == parentID {
			// If IncludePurpose flag is set, generate purpose
			// TODO
			purpose := ""

			docSections = append(docSections, config.RuleDocumentSection{
				ID:       section.Name,
				Purpose:  purpose,
				Required: true,
				Sections: createRuleSections(sections, section.ID),
			})
		}
	}

	return docSections
}
