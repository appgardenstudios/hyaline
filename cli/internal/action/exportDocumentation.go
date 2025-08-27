package action

import (
	"context"
	"errors"
	"fmt"
	"hyaline/internal/config"
	"hyaline/internal/docs"
	"hyaline/internal/io"
	"hyaline/internal/sqlite"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"sort"
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

	// Parse and validate includes/excludes
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
		return errors.New("output path already exists")
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

	// Get sources
	sources, err := docDB.GetAllSources(context.Background())
	if err != nil {
		slog.Debug("action.ExportDocumentation could not get sources", "error", err)
		return err
	}
	sourcesMap := make(map[string]*sqlite.SOURCE)
	for _, source := range sources {
		sourcesMap[source.ID] = &source
	}

	// Output documentation
	switch format {
	case ExportFormatFs:
		err = exportFs(documents, sourcesMap, args.Includes, args.Excludes, args.Documentation, outputAbsPath)
	case ExportFormatLlmsFullTxt:
		err = exportLlmsFullTxt(documents, outputAbsPath)
	case ExportFormatJson:
		err = exportJson(documents, outputAbsPath)
	case ExportFormatSqlite:
		err = exportSqlite(documents, sourcesMap, outputAbsPath)
	default:
		err = fmt.Errorf("unknown format %s", format.String())
	}
	if err != nil {
		slog.Debug("action.ExportDocumentation could not export", "error", err)
		return err
	}

	slog.Info("Export Complete")

	return nil
}

func exportFs(documents []*docs.FilteredDoc, sources map[string]*sqlite.SOURCE, includes []string, excludes []string, inputPath string, outputPath string) (err error) {
	slog.Info(fmt.Sprintf("Exporting %d documents to %s", len(documents), outputPath))

	// Create path
	err = os.MkdirAll(outputPath, 0755)
	if err != nil {
		slog.Debug("action.exportFs could not MkdirAll for outputPath", "outputPath", outputPath, "error", err)
		return
	}

	// Record the number of documents for each source for our README
	sourcesCount := make(map[string]int)

	// Output documents
	for _, document := range documents {
		count := sourcesCount[document.Document.SourceID]
		sourcesCount[document.Document.SourceID] = count + 1

		// Get dir and filename
		dir := path.Join(outputPath, document.Document.SourceID, filepath.Dir(document.Document.ID))
		var filename string
		if strings.HasSuffix(document.Document.ID, "/") {
			filename = "index.md"
		} else {
			filename = filepath.Base(document.Document.ID)
		}
		// Only add ".md" to the filename if it doesn't exist so we don't end up with file.md.md
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
	readmeIncludes := includes
	readmeExcludes := excludes
	if len(readmeIncludes) == 0 {
		readmeIncludes = append(readmeIncludes, "(all)")
	}
	if len(readmeExcludes) == 0 {
		readmeExcludes = append(readmeExcludes, "(none)")
	}
	readmeSources := []string{}
	for sourceID, count := range sourcesCount {
		description := ""
		source := sources[sourceID]
		if source != nil {
			description = source.Description
		}
		readmeSources = append(readmeSources, fmt.Sprintf("%s - %s (%d)", sourceID, description, count))
	}
	if len(readmeSources) == 0 {
		readmeSources = append(readmeSources, "(none)")
	}
	contents := fmt.Sprintf(`# Exported Documentation
Documentation exported from `+"`"+`%s`+"`"+`

**Includes**:
  - `+"`"+`%s`+"`"+`

**Excludes**:
  - `+"`"+`%s`+"`"+`

**Documents Exported**: %d

## Sources
- %s
`, inputPath, strings.Join(readmeIncludes, "`\n  - `"), strings.Join(readmeExcludes, "`\n  - `"), len(documents), strings.Join(readmeSources, "\n- "))
	readmePath := path.Join(outputPath, "README.md")
	var file *os.File
	file, err = os.Create(readmePath)
	if err != nil {
		slog.Debug("action.exportFs could not create file", "readmePath", readmePath, "error", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(contents)
	if err != nil {
		slog.Debug("action.exportFs could not write string to file", "readmePath", readmePath, "error", err)
		return
	}

	return
}

func exportLlmsFullTxt(documents []*docs.FilteredDoc, outputPath string) (err error) {
	slog.Info(fmt.Sprintf("Exporting %d documents to %s", len(documents), outputPath))

	// Sort documents
	sort.SliceStable(documents, func(i int, j int) bool {
		if documents[i].Document.SourceID < documents[j].Document.SourceID {
			return true
		}
		if documents[i].Document.SourceID > documents[j].Document.SourceID {
			return false
		}
		return documents[i].Document.ID < documents[j].Document.ID
	})

	// Format output
	var str strings.Builder
	for i, document := range documents {
		var title string
		if len(document.Sections) > 0 {
			title = document.Sections[0].Section.Name
		} else {
			title = document.Document.ID
		}
		if i > 0 {
			str.WriteString("\n\n\n")
		}
		str.WriteString(fmt.Sprintf("# %s\n", title))
		str.WriteString(fmt.Sprintf("Source: document://%s/%s\n", document.Document.SourceID, document.Document.ID))
		str.WriteString("\n")
		str.WriteString(strings.TrimSpace(document.Document.ExtractedData))
	}

	// Write output
	var file *os.File
	file, err = os.Create(outputPath)
	if err != nil {
		slog.Debug("action.exportLlmsFullTxt could not create file", "outputPath", outputPath, "error", err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(str.String())
	if err != nil {
		slog.Debug("action.exportLlmsFullTxt could not write string to file", "outputPath", outputPath, "error", err)
		return
	}

	return
}

func exportJson(documents []*docs.FilteredDoc, outputPath string) (err error) {
	slog.Info(fmt.Sprintf("Exporting %d documents to %s", len(documents), outputPath))

	type outputTag struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	type outputDocument struct {
		Source   string      `json:"source"`
		Document string      `json:"document"`
		URI      string      `json:"uri"`
		Purpose  string      `json:"purpose,omitempty"`
		Content  string      `json:"content"`
		Tags     []outputTag `json:"tags"`
	}

	// Format output
	output := make([]outputDocument, 0)
	for _, document := range documents {
		tags := make([]outputTag, 0)
		for _, tag := range document.Tags {
			tags = append(tags, outputTag{
				Key:   tag.Key,
				Value: tag.Value,
			})
		}
		uri := docs.DocumentURI{
			SourceID:     document.Document.SourceID,
			DocumentPath: document.Document.ID,
		}
		output = append(output, outputDocument{
			Source:   uri.SourceID,
			Document: uri.DocumentPath,
			URI:      uri.String(),
			Purpose:  document.Document.Purpose,
			Content:  document.Document.ExtractedData,
			Tags:     tags,
		})
	}

	// Sort output
	sort.SliceStable(output, func(i int, j int) bool {
		if output[i].Source < output[j].Source {
			return true
		}
		if output[i].Source > output[j].Source {
			return false
		}
		return output[i].Document < output[j].Document
	})

	// Write JSON
	outputFile, err := os.Create(outputPath)
	if err != nil {
		slog.Debug("action.exportJson could not open output file", "error", err)
		return err
	}
	defer outputFile.Close()
	io.WriteJSON(outputFile, output)
	return
}

func exportSqlite(documents []*docs.FilteredDoc, sources map[string]*sqlite.SOURCE, outputPath string) (err error) {
	// Initialize our output database
	docDb, close, err := sqlite.InitOutput(outputPath)
	if err != nil {
		slog.Debug("action.ExtractDocumentation could not initialize output", "error", err)
		return err
	}
	defer close()

	// Record the set of sources seen
	sourcesSeen := make(map[string]struct{})

	for _, document := range documents {
		// Mark source as seem
		sourcesSeen[document.Document.SourceID] = struct{}{}

		// Insert document
		err = docDb.InsertDocument(context.Background(), sqlite.InsertDocumentParams{
			ID:            document.Document.ID,
			SourceID:      document.Document.SourceID,
			Type:          document.Document.Type,
			Purpose:       document.Document.Purpose,
			RawData:       document.Document.RawData,
			ExtractedData: document.Document.ExtractedData,
		})
		if err != nil {
			// TODO debug
			return
		}

		// Insert document tags
		for _, tag := range document.Tags {
			err = docDb.UpsertDocumentTag(context.Background(), sqlite.UpsertDocumentTagParams{
				SourceID:   document.Document.SourceID,
				DocumentID: document.Document.ID,
				TagKey:     tag.Key,
				TagValue:   tag.Value,
			})
			if err != nil {
				// TODO debug
				return
			}
		}

		// Insert sections and tags
		if len(document.Sections) > 0 {
			err = exportSqliteSections(document.Sections, docDb)
			if err != nil {
				// TODO debug
				return
			}
		}
	}

	// Insert sources
	for id := range sourcesSeen {
		source := sources[id]
		err = docDb.InsertSource(context.Background(), sqlite.InsertSourceParams{
			ID:          source.ID,
			Description: source.Description,
			Crawler:     source.Crawler,
			Root:        source.Root,
		})
		if err != nil {
			// TODO debug
			return
		}
	}

	return
}

func exportSqliteSections(sections []docs.FilteredSection, docDb *sqlite.Queries) (err error) {
	for _, section := range sections {
		// Insert section
		err = docDb.InsertSection(context.Background(), sqlite.InsertSectionParams{
			ID:            section.Section.ID,
			DocumentID:    section.Section.DocumentID,
			SourceID:      section.Section.SourceID,
			ParentID:      section.Section.ParentID,
			PeerOrder:     section.Section.PeerOrder,
			Name:          section.Section.Name,
			Purpose:       section.Section.Purpose,
			ExtractedData: section.Section.ExtractedData,
		})
		if err != nil {
			return
		}

		// Insert section tags
		for _, tag := range section.Tags {
			err = docDb.UpsertSectionTag(context.Background(), sqlite.UpsertSectionTagParams{
				SourceID:   section.Section.SourceID,
				DocumentID: section.Section.DocumentID,
				SectionID:  section.Section.ID,
				TagKey:     tag.Key,
				TagValue:   tag.Value,
			})
			if err != nil {
				return
			}
		}

		// If children, recurse
		if len(section.Sections) > 0 {
			err = exportSqliteSections(section.Sections, docDb)
			if err != nil {
				return
			}
		}
	}

	return
}
