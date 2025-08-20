package action

import (
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"strings"
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
		slog.Debug("action.ExportDocumentation could not get filtered documents", "error", err)
		return err
	}
	slog.Info("Retrieved filtered documents", "documents", len(documents))

	// Output documentation
	err = exportFs(documents, outputAbsPath)
	if err != nil {
		slog.Debug("action.ExportDocumentation could not export to path", "error", err)
		return err
	}

	slog.Info("Export Complete")

	return nil
}

func exportFs(documents []*docs.FilteredDoc, outputPath string) (err error) {
	slog.Info(fmt.Sprintf("Exporting %d documents to %s", len(documents), outputPath))

	// Create path
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		slog.Debug("action.exportFs could not MkdirAll for outputPath", "outputPath", outputPath, "error", err)
		return
	}

	// Output documents
	for _, document := range documents {
		// Get dir and filename
		dir := path.Join(outputPath, document.Document.SourceID, filepath.Dir(document.Document.ID))
		var filename string
		if strings.HasSuffix(document.Document.ID, "/") {
			filename = "index.md"
		} else {
			filename = filepath.Base(document.Document.ID)
		}
		if !strings.HasSuffix(filename, ".md") {
			filename = filename + ".md"
		}
		finalPath := path.Join(dir, filename)
		slog.Debug("Writing document",
			"source", document.Document.SourceID,
			"document", document.Document.ID,
			"finalPath", finalPath)

		// Ensure dir exists
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			slog.Debug("action.exportFs could not MkdirAll for dir", "dir", dir, "error", err)
			return
		}

		// Output file
		var file *os.File
		file, err = os.Create(finalPath)
		if err != nil {
			slog.Debug("action.exportFs could not create file", "finalPath", finalPath, "error", err)
			return
		}
		defer file.Close()
		_, err = file.WriteString(document.Document.ExtractedData)
		if err != nil {
			slog.Debug("action.exportFs could not write string to file", "finalPath", finalPath, "error", err)
			return
		}
	}

	// Format and output README.md
	// TODO

	return
}
