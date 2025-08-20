package action

import (
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path/filepath"
)

type ExportDocumentationArgs struct {
	Documentation string
	Format        string
	Includes      []string
	Excludes      []string
	Output        string
}

type ExportFormatType string

func (t ExportFormatType) String() string {
	return string(t)
}

func (t ExportFormatType) IsValid() bool {
	switch t {
	case ExportFormatFs, ExportFormatLlmsFullTxt, ExportFormatJson, ExportFormatSqlite:
		return true
	default:
		return false
	}
}

func (t ExportFormatType) PossibleValues() string {
	return fmt.Sprintf("%s, %s, %s, %s", ExportFormatFs, ExportFormatLlmsFullTxt, ExportFormatJson, ExportFormatSqlite)
}

const (
	ExportFormatFs          ExportFormatType = "fs"
	ExportFormatLlmsFullTxt ExportFormatType = "llmsfulltxt"
	ExportFormatJson        ExportFormatType = "json"
	ExportFormatSqlite      ExportFormatType = "sqlite"
)

func ExportDocumentation(args *ExportDocumentationArgs) error {
	slog.Info("Exporting Documentation",
		"documentation", args.Documentation,
		"format", args.Format,
		"includes", args.Includes,
		"excludes", args.Excludes,
		"output", args.Output,
	)

	// Validate format
	format := ExportFormatType(args.Format)
	if !format.IsValid() {
		slog.Debug("action.ExportDocumentation received an invalid format")
		return fmt.Errorf("invalid format, got: %s, wanted one of: %s", format.String(), format.PossibleValues())
	}

	// Parse/Validate includes/excludes
	includes := []*docs.DocumentURI{}
	excludes := []*docs.DocumentURI{}
	for _, include := range args.Includes {
		uri, err := docs.NewDocumentURI(include)
		if err != nil {
			slog.Debug("action.ExportDocumentation could not parse include", "include", include)
			return err
		}
		includes = append(includes, uri)
	}
	for _, exclude := range args.Excludes {
		uri, err := docs.NewDocumentURI(exclude)
		if err != nil {
			slog.Debug("action.ExportDocumentation could not parse exclude", "exclude", exclude)
			return err
		}
		excludes = append(excludes, uri)
	}

	// Ensure output path does not exist
	outputAbsPath, err := filepath.Abs(args.Output)
	if err != nil {
		slog.Debug("action.ExportDocumentation could not get an absolute path for output", "output", args.Output, "error", err)
		return err
	}
	_, err = os.Stat(outputAbsPath)
	if err == nil {
		slog.Debug("action.ExportDocumentation detected that output already exists", "absPath", outputAbsPath)
		return errors.New("output file already exists")
	}

	// Open Documentation database
	docDB, err := sqlite.InitInput(args.Documentation)
	if err != nil {
		slog.Debug("action.ExportDocumentation could not initialize documentation db", "documentation", args.Documentation, "error", err)
		return err
	}

	// Get documentation
	filters := &config.CheckDocumentation{}
	if len(includes) == 0 {
		filters.Include = append(filters.Include, config.DocumentationFilter{
			Source:   "*",
			Document: "**/*",
		})
	} else {
		for _, include := range includes {
			filters.Include = append(filters.Include, config.DocumentationFilter{
				Source:   include.SourceID,
				Document: include.DocumentPath,
				Tags:     include.Tags.ToDocumentationFilterTags(),
			})
		}
	}
	for _, exclude := range excludes {
		filters.Exclude = append(filters.Exclude, config.DocumentationFilter{
			Source:   exclude.SourceID,
			Document: exclude.DocumentPath,
			Tags:     exclude.Tags.ToDocumentationFilterTags(),
		})
	}
	documents, err := docs.GetFilteredDocs(filters, docDB)
	if err != nil {
		slog.Debug("action.CheckDiff could not get filtered documents", "error", err)
		return err
	}
	slog.Info("Retrieved filtered documents", "documents", len(documents))

	// Output documentation
	// TODO

	return nil
}
