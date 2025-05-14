package action

import (
	"database/sql"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/llm"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
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
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.GenerateConfig could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.GenerateConfig detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open current db
	currentAbsPath, err := filepath.Abs(args.Current)
	if err != nil {
		slog.Debug("action.GenerateConfig could not get an absolute path for current", "current", args.Current, "error", err)
		return err
	}
	currentDB, err := sql.Open("sqlite", currentAbsPath)
	if err != nil {
		slog.Debug("action.GenerateConfig could not open current SQLite DB", "dataSourceName", currentAbsPath, "error", err)
		return err
	}
	slog.Debug("action.GenerateConfig opened current database", "current", args.Current, "path", currentAbsPath)
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
	for _, d := range system.DocumentationSources {
		// Create a desired document for this documentation
		desiredSocSet := config.DocumentSet{
			ID:        d.ID,
			Documents: []config.Document{},
		}

		// Get a list of Documents from the db for this doc ID
		documents, err := sqlite.GetAllDocument(d.ID, system.ID, currentDB)
		if err != nil {
			slog.Debug("action.GenerateConfig could not get documents from current db", "doc", d.ID, "system", system.ID, "error", err)
			return err
		}

		// Loop through each document to generate desired documents for it
		for _, doc := range documents {
			slog.Info("Processing document", "document", doc.ID)
			// Get the desiredDoc for this document (if any)
			desiredDoc, desiredDocFound := d.GetDocument(cfg, doc.ID)

			// If there is no desired document found, create it
			if !desiredDocFound {
				// If IncludePurpose flag is set, get purpose
				purpose := ""
				if args.IncludePurpose {
					purpose, err = llm.GetDocumentPurpose(doc.ID, doc.ExtractedData, &cfg.LLM)
					if err != nil {
						slog.Debug("action.GenerateConfig could not get purpose for document", "document", doc.ID, "doc", d.ID, "system", system.ID, "error", err)
						return err
					}
				}

				// Create desired document for the document
				desiredDoc = config.Document{
					Name:     doc.ID,
					Purpose:  purpose,
					Required: true,
				}
			}

			// Get and add sections for this document
			sections, err := sqlite.GetAllSectionsForDocument(doc.ID, d.ID, system.ID, currentDB)
			if err != nil {
				slog.Debug("action.GenerateConfig could not get sections for a document from current db", "document", doc.ID, "doc", d.ID, "system", system.ID, "error", err)
				return err
			}
			newSections, err := createRuleSections(sections, doc.ID, desiredDoc.Sections, args.IncludePurpose, doc.ID, desiredDoc.Purpose, &cfg.LLM)
			if err != nil {
				slog.Debug("action.GenerateConfig could not generate sections for a document from current db", "document", doc.ID, "doc", d.ID, "system", system.ID, "error", err)
				return err
			}
			desiredDoc.Sections = newSections

			// Add desired document to the set
			desiredSocSet.Documents = append(desiredSocSet.Documents, desiredDoc)
		}

		// Add desired document set to config
		newCfg.CommonDocuments = append(newCfg.CommonDocuments, desiredSocSet)
	}

	// Output new config
	yml, err := yaml.Marshal(newCfg)
	if err != nil {
		slog.Debug("action.GenerateConfig could not marshal yaml", "error", err)
		return err
	}
	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.GenerateConfig could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()

	// Write the byte slice to the file
	_, err = outputFile.Write(yml)
	if err != nil {
		slog.Debug("action.GenerateConfig could not write output file", "error", err)
		return err
	}

	return nil
}

// Note: sections MUST be in PEER_ORDER so that the doc sections are added in the correct order
func createRuleSections(sections []*sqlite.Section, parentID string, existingSections []config.DocumentSection, includePurpose bool, documentName string, documentPurpose string, cfg *config.LLM) (docSections []config.DocumentSection, err error) {
	for _, section := range sections {
		// Guard against circular issues by ensuring that no ID is the same as its parent ID
		if section.ID == section.ParentID {
			err = fmt.Errorf("circular section found: %s", section.ID)
			return
		}

		// Add this section if it is a child of the parent we are currently building out
		if section.ParentID == parentID {
			// See if section already exists
			sectionFound, docSection := getRuleSection(section.Name, existingSections)

			// If section not found, create it
			if !sectionFound {
				// If IncludePurpose flag is set, get purpose
				purpose := ""
				if includePurpose {
					purpose, err = llm.GetSectionPurpose(documentName, documentPurpose, section.Name, section.ExtractedData, cfg)
					if err != nil {
						slog.Debug("action.GenerateConfig could not get purpose for section", "section", section.ID, "error", err)
						return
					}
				}

				// Create new doc section
				docSection = config.DocumentSection{
					Name:     section.Name,
					Purpose:  purpose,
					Required: true,
				}
			}

			// Get and add child doc sections
			var childDocSections []config.DocumentSection
			childDocSections, err = createRuleSections(sections, section.ID, docSection.Sections, includePurpose, documentName, documentPurpose, cfg)
			if err != nil {
				return
			}
			docSection.Sections = childDocSections

			// Add the section to the list
			docSections = append(docSections, docSection)
		}
	}

	return docSections, nil
}

func getRuleSection(sectionID string, sections []config.DocumentSection) (sectionFound bool, section config.DocumentSection) {
	for _, section = range sections {
		if section.Name == sectionID {
			return true, section
		}
	}

	return
}
