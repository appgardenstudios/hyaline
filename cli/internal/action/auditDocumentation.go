package action

import (
	"encoding/json"
	"fmt"
	"hyaline/internal/audit"
	"hyaline/internal/config"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type AuditDocumentationArgs struct {
	Config        string
	Documentation string
	Sources       []string
	Output        string
}

func AuditDocumentation(args *AuditDocumentationArgs) error {
	slog.Info("Auditing documentation",
		"config", args.Config,
		"documentation", args.Documentation,
		"sources", args.Sources,
		"output", args.Output)

	// Load Config
	cfg, err := config.Load(args.Config)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not load the config", "error", err)
		return err
	}

	// Ensure audit configuration exists
	if cfg.Audit == nil {
		return fmt.Errorf("audit configuration not found in config file")
	}

	// Ensure output JSON file does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.AuditDocumentation detected that output already exists", "absPath", outputAbsPath)
		return fmt.Errorf("output file already exists")
	}

	// Initialize documentation database
	db, err := sqlite.InitInput(args.Documentation)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not initialize documentation database", "documentation", args.Documentation, "error", err)
		return err
	}

	slog.Debug("action.AuditDocumentation initialized documentation database", "documentation", args.Documentation)

	auditResults, err := audit.Documentation(cfg, db, args.Sources)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not run audit", "error", err)
		return err
	}

	// Write results to JSON file
	jsonData, err := json.MarshalIndent(auditResults, "", "  ")
	if err != nil {
		slog.Debug("action.AuditDocumentation could not marshal JSON", "error", err)
		return err
	}

	outputFile, err := os.Create(outputAbsPath)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not create output file", "error", err)
		return err
	}
	defer outputFile.Close()

	_, err = outputFile.Write(jsonData)
	if err != nil {
		slog.Debug("action.AuditDocumentation could not write output file", "error", err)
		return err
	}

	slog.Info("Audit documentation completed successfully")
	return nil
}
