package docs

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/repo"
	"hyaline/internal/sqlite"
	"io/fs"
	"log/slog"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/gocolly/colly/v2"
)

func ExtractCurrent(system *config.System, db *sql.DB) (err error) {
	// Process each docs source
	for _, d := range system.DocumentationSources {
		slog.Debug("docs.ExtractCurrent extracting docs", "system", system.ID, "docs", d.ID)

		// Get document path
		var path string
		switch d.Extractor.Type {
		case config.ExtractorTypeFs:
			path = d.Extractor.Options.Path
		case config.ExtractorTypeGit:
			path = d.Extractor.Options.Path
			if path == "" {
				path = d.Extractor.Options.Repo
			}
		case config.ExtractorTypeHttp:
			u, err := url.Parse(d.Extractor.Options.BaseURL)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not parse base url", "system", system.ID, "docs", d.ID, "baseUrl", d.Extractor.Options.BaseURL)
			}
			path = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
		}

		// Insert Documentation
		err = sqlite.InsertSystemDocumentation(sqlite.SystemDocumentation{
			ID:       d.ID,
			SystemID: system.ID,
			Type:     d.Type.String(),
			Path:     path,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrent could not insert documentation", "error", err, "doc", d.ID)
			return err
		}

		// Extract based on the extractor
		switch d.Extractor.Type {
		case config.ExtractorTypeFs:
			err = ExtractCurrentFs(system.ID, &d, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract docs using fs extractor", "error", err, "doc", d.ID)
				return
			}
		case config.ExtractorTypeGit:
			err = ExtractCurrentGit(system.ID, &d, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract docs using git extractor", "error", err, "doc", d.ID)
				return
			}
		case config.ExtractorTypeHttp:
			err = ExtractCurrentHttp(system.ID, &d, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrent could not extract docs using http extractor", "error", err, "doc", d.ID)
				return
			}
		default:
			slog.Debug("docs.ExtractCurrent unknown extractor", "extractor", d.Extractor.Type.String(), "doc", d.ID)
			return errors.New("Unknown Extractor '" + d.Extractor.Type.String() + "' for doc " + d.ID)
		}
	}

	return
}

func ExtractCurrentFs(systemID string, d *config.DocumentationSource, db *sql.DB) (err error) {
	// Get our absolute path
	absPath, err := filepath.Abs(d.Extractor.Options.Path)
	if err != nil {
		slog.Debug("docs.ExtractCurrentFs could not determine absolute docs path", "error", err, "path", d.Extractor.Options.Path)
		return err
	}
	slog.Debug("docs.ExtractCurrentFs extracting docs from path", "absPath", absPath)

	// Get our root FS
	// We use a root FS so symlinks and relative paths don't escape our path
	// https://pkg.go.dev/os@go1.24.1#Root
	root, err := os.OpenRoot(absPath)
	if err != nil {
		slog.Debug("code.ExtractCurrentFs could not open fs root", "error", err, "path", d.Extractor.Options.Path)
		return err
	}
	fsRoot := root.FS()

	// Our set of files (as a map so we don't get dupes)
	docs := map[string]struct{}{}

	// Loop through our includes and get files
	for _, include := range d.Extractor.Include {
		slog.Debug("docs.ExtractCurrentFs extracting docs using include", "include", include, "doc", d.ID)

		// Get matched docs
		matches, err := doublestar.Glob(fsRoot, include)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs could not find docs files with include", "include", include, "error", err)
			return err
		}

		// Loop through docs and add those that match includes and don't match excludes
		for _, doc := range matches {
			if config.PathIsIncluded(doc, []string{include}, d.Extractor.Exclude) {
				docs[doc] = struct{}{}
			} else {
				slog.Debug("docs.ExtractCurrentFs doc excluded", "doc", doc)
			}
		}
	}

	// Insert docs
	for doc := range docs {
		// Get file rawData
		rawData, err := fs.ReadFile(fsRoot, doc)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs could not read doc file", "error", err, "doc", doc)
			return err
		}

		// Extract and clean data (trim whitespace and remove carriage returns)
		var extractedData string
		switch d.Type {
		case config.DocTypeHTML:
			extractedData, err = extractHTMLDocument(string(rawData), d.Options.Selector)
			if err != nil {
				slog.Debug("docs.ExtractCurrentFs could not extract html document", "error", err, "doc", doc)
				return err
			}
		default:
			extractedData = strings.TrimSpace(string(rawData))
		}
		extractedData = strings.ReplaceAll(extractedData, "\r", "")

		// Insert our document
		err = sqlite.InsertSystemDocument(sqlite.SystemDocument{
			ID:              doc,
			DocumentationID: d.ID,
			SystemID:        systemID,
			Type:            d.Type.String(),
			Action:          sqlite.ActionNone,
			OriginalID:      "",
			RawData:         string(rawData),
			ExtractedData:   extractedData,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs could not insert document", "error", err)
			return err
		}

		// Get and insert sections
		sections := getMarkdownSections(strings.Split(extractedData, "\n"))
		err = insertMarkdownSectionAndChildren(sections, 0, doc, d.ID, systemID, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrentFs could not insert section", "error", err)
			return err
		}
	}

	return
}

func ExtractCurrentGit(systemID string, d *config.DocumentationSource, db *sql.DB) (err error) {
	// Initialize go-git repo (on disk or in mem)
	var r *git.Repository
	r, err = repo.GetRepo(d.Extractor.Options)
	if err != nil {
		slog.Debug("docs.ExtractCurrentGit could not get repo", "error", err)
		return
	}

	// Extract files from branch
	branch := "main"
	if d.Extractor.Options.Branch != "" {
		branch = d.Extractor.Options.Branch
	}
	err = repo.GetFiles(branch, r, func(f *object.File) error {
		if config.PathIsIncluded(f.Name, d.Extractor.Include, d.Extractor.Exclude) {
			// Get contents of file and insert into db
			var bytes []byte
			bytes, err = repo.GetBlobBytes(f.Blob)
			if err != nil {
				slog.Debug("docs.ExtractCurrentGit could not get blob bytes", "error", err)
				return err
			}

			// Extract and clean data (trim whitespace and remove carriage returns)
			var extractedData string
			switch d.Type {
			case config.DocTypeHTML:
				extractedData, err = extractHTMLDocument(string(bytes), d.Options.Selector)
				if err != nil {
					slog.Debug("docs.ExtractCurrentGit could not extract html document", "error", err, "doc", f.Name)
					return err
				}
			default:
				extractedData = strings.TrimSpace(string(bytes))
			}
			extractedData = strings.ReplaceAll(extractedData, "\r", "")

			// Insert our document
			err = sqlite.InsertSystemDocument(sqlite.SystemDocument{
				ID:              f.Name,
				DocumentationID: d.ID,
				SystemID:        systemID,
				Type:            d.Type.String(),
				Action:          sqlite.ActionNone,
				OriginalID:      "",
				RawData:         string(bytes),
				ExtractedData:   extractedData,
			}, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrentGit could not insert document", "error", err)
				return err
			}

			// Get and insert sections
			sections := getMarkdownSections(strings.Split(extractedData, "\n"))
			err = insertMarkdownSectionAndChildren(sections, 0, f.Name, d.ID, systemID, db)
			if err != nil {
				slog.Debug("docs.ExtractCurrentGit could not insert section", "error", err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		slog.Debug("code.ExtractCurrentGit could not get files", "error", err)
		return
	}

	return
}

func ExtractCurrentHttp(systemID string, d *config.DocumentationSource, db *sql.DB) error {
	// Collect encountered errors into an array and check it at the end
	var errs []error

	// Use baseURL to calculate includes/excludes
	baseUrl, err := url.Parse(d.Extractor.Options.BaseURL)
	if err != nil {
		slog.Debug("docs.ExtractCurrentHttp could not parse baseUrl", "baseUrl", d.Extractor.Options.BaseURL, "error", err)
		return err
	}
	var includes []string
	for _, include := range d.Extractor.Include {
		includes = append(includes, path.Join(baseUrl.Path, include))
	}
	var excludes []string
	for _, exclude := range d.Extractor.Exclude {
		excludes = append(excludes, path.Join(baseUrl.Path, exclude))
	}
	slog.Debug("docs.ExtractCurrentHttp includes/excludes", "includes", includes, "excludes", excludes, "basePath", baseUrl.Path)

	// Determine start URL
	var startUrl *url.URL
	if d.Extractor.Options.Start != "" {
		startUrl, err = url.Parse(d.Extractor.Options.Start)
		if err != nil {
			slog.Debug("docs.ExtractCurrentHttp could not parse start", "start", d.Extractor.Options.Start, "error", err)
			return err
		}
		startUrl = baseUrl.ResolveReference(startUrl)
	} else {
		startUrl = baseUrl
	}
	slog.Debug("docs.ExtractCurrentHttp startUrl", "startUrl", startUrl, "start", d.Extractor.Options.Start, "baseUrl", baseUrl.String())

	// Initialize our collector
	c := colly.NewCollector(
		colly.Async(),
	)

	// Create our default limits
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Delay:       0,
		Parallelism: 1,
	})

	// Add headers (if any)
	for key, val := range d.Extractor.Options.Headers {
		c.Headers.Add(key, val)
	}

	// Find and visit all links in returned html
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Get a resolved URL for the href relative to the base URL of the requested page
		href := e.Attr("href")
		raw, err := url.Parse(href)
		if err != nil {
			slog.Debug("docs.ExtractCurrentHttp unable to parse href", "href", href, "error", err)
			errs = append(errs, err)
			return
		}
		u := e.Request.URL.ResolveReference(raw)
		slog.Debug("docs.ExtractCurrentHttp evaluating href", "href", href, "url", u.String())

		// Only visit pages on this same host
		if e.Request.URL.Host != u.Host {
			slog.Debug("docs.ExtractCurrentHttp skipping external link", "href", href)
			return
		}

		// Only visit if this path matches an include (and does not match an exclude)
		if config.PathIsIncluded(u.Path, includes, excludes) {
			slog.Debug("docs.ExtractCurrentHttp visiting URL", "href", href, "url", u.String(), "currentPage", e.Request.URL.String())
			e.Request.Visit(href)
		} else {
			slog.Debug("docs.ExtractCurrentHttp URL excluded", "href", href, "url", u.String())
		}
	})

	// Save documents we scrape
	c.OnResponse(func(r *colly.Response) {
		path := r.Request.URL.Path

		// Extract and clean data (trim whitespace and remove carriage returns)
		var extractedData string
		var err error
		switch d.Type {
		case config.DocTypeHTML:
			extractedData, err = extractHTMLDocument(string(r.Body), d.Options.Selector)
			if err != nil {
				slog.Debug("docs.ExtractCurrentGit could not extract html document", "error", err, "doc", path)
				errs = append(errs, err)
				return
			}
		default:
			extractedData = strings.TrimSpace(string(r.Body))
		}
		extractedData = strings.ReplaceAll(extractedData, "\r", "")

		// Insert our document
		err = sqlite.InsertSystemDocument(sqlite.SystemDocument{
			ID:              path,
			DocumentationID: d.ID,
			SystemID:        systemID,
			Type:            d.Type.String(),
			Action:          sqlite.ActionNone,
			OriginalID:      "",
			RawData:         string(r.Body),
			ExtractedData:   extractedData,
		}, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrentGit could not insert document", "error", err)
			errs = append(errs, err)
			return
		}

		// Get and insert sections
		sections := getMarkdownSections(strings.Split(extractedData, "\n"))
		err = insertMarkdownSectionAndChildren(sections, 0, path, d.ID, systemID, db)
		if err != nil {
			slog.Debug("docs.ExtractCurrentGit could not insert section", "error", err)
			errs = append(errs, err)
			return
		}
	})

	// Record any encountered errors
	c.OnError(func(r *colly.Response, e error) {
		errs = append(errs, e)
	})

	// Visit and wait for all routines to return
	c.Visit(startUrl.String())
	c.Wait()

	// Log and handle any errors we encountered
	if len(errs) > 0 {
		err := errors.New("http extractor encountered " + fmt.Sprint(len(errs)) + " errors")
		slog.Debug("docs.ExtractCurrentHttp encountered errors", "errors", errs, "error", err)
		return err
	}

	return nil
}

func insertMarkdownSectionAndChildren(s *section, order int, documentId string, documentationId string, systemId string, db *sql.DB) error {
	// Insert this section
	parentId := ""
	if s.Parent != nil {
		parentId = documentId + s.Parent.FullName
	}
	err := sqlite.InsertSystemSection(sqlite.SystemSection{
		ID:              documentId + s.FullName,
		DocumentID:      documentId,
		DocumentationID: documentationId,
		SystemID:        systemId,
		Name:            s.Name,
		ParentID:        parentId,
		PeerOrder:       order,
		ExtractedData:   strings.TrimSpace(s.Content),
	}, db)
	if err != nil {
		slog.Debug("docs.insertSectionAndChildren could not insert section", "error", err)
		return err
	}

	// Insert children
	for i, child := range s.Children {
		err = insertMarkdownSectionAndChildren(child, i, documentId, documentationId, systemId, db)
		if err != nil {
			return err
		}
	}

	return nil
}
