package extract

import (
	"errors"
	"fmt"
	"hyaline/internal/config"
	"log/slog"
	"net/url"
	"path"
	"sync"

	"github.com/gocolly/colly/v2"
)

func crawlHttp(cfg *config.ExtractCrawler, cb extractorCallback) error {
	// Collect encountered errors into an array and check it at the end
	var errs []error

	// Use baseURL to calculate includes/excludes
	baseUrl, err := url.Parse(cfg.Options.BaseURL)
	if err != nil {
		slog.Debug("extract.crawlHttp could not parse baseUrl", "baseUrl", cfg.Options.BaseURL, "error", err)
		return err
	}
	var includes []string
	for _, include := range cfg.Include {
		includes = append(includes, path.Join(baseUrl.Path, include))
	}
	var excludes []string
	for _, exclude := range cfg.Exclude {
		excludes = append(excludes, path.Join(baseUrl.Path, exclude))
	}
	slog.Debug("extract.crawlHttp includes/excludes", "includes", includes, "excludes", excludes, "basePath", baseUrl.Path)

	// Determine start URL
	var startUrl *url.URL
	if cfg.Options.Start != "" {
		startUrl, err = url.Parse(cfg.Options.Start)
		if err != nil {
			slog.Debug("extract.crawlHttp could not parse start", "start", cfg.Options.Start, "error", err)
			return err
		}
		startUrl = baseUrl.ResolveReference(startUrl)
	} else {
		startUrl = baseUrl
	}
	slog.Info("Crawling documentation using http", "startUrl", startUrl.String())
	slog.Debug("extract.crawlHttp startUrl", "startUrl", startUrl, "start", cfg.Options.Start, "baseUrl", baseUrl.String())

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

	// Create our mutex to prevent simultaneous writes to SQLite
	var mutex sync.Mutex

	// Add headers (if any)
	for key, val := range cfg.Options.Headers {
		c.Headers.Add(key, val)
	}

	// Find and visit all links in returned html
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		// Get a resolved URL for the href relative to the base URL of the requested page
		href := e.Attr("href")
		raw, err := url.Parse(href)
		if err != nil {
			slog.Debug("extract.crawlHttp unable to parse href", "href", href, "error", err)
			errs = append(errs, err)
			return
		}
		u := e.Request.URL.ResolveReference(raw)
		slog.Debug("extract.crawlHttp evaluating href", "href", href, "url", u.String())

		// Only visit pages on this same host
		if e.Request.URL.Host != u.Host {
			slog.Debug("extract.crawlHttp skipping external link", "href", href)
			return
		}

		// Only visit if this path matches an include (and does not match an exclude)
		if config.PathIsIncluded(u.Path, includes, excludes) {
			// Visit only the main part of the URL (protocol, host, path) without fragments or query params
			urlToVisit := &url.URL{
				Scheme: u.Scheme,
				Host:   u.Host,
				Path:   u.Path,
			}

			slog.Debug("extract.crawlHttp visiting URL", "href", href, "urlToVisit", urlToVisit.String(), "currentPage", e.Request.URL.String())
			err = e.Request.Visit(urlToVisit.String())

			if err != nil {
				var alreadyVisitedError *colly.AlreadyVisitedError
				if errors.As(err, &alreadyVisitedError) {
					slog.Debug("extract.crawlHttp skipping already visited URL", "href", href, "urlToVisit", urlToVisit.String())
					return
				}

				slog.Debug("extract.crawlHttp could not visit URL", "href", href, "urlToVisit", urlToVisit.String(), "error", err)
				errs = append(errs, err)
				return
			}
		} else {
			slog.Debug("extract.crawlHttp URL excluded", "href", href, "url", u.String())
		}
	})

	// Save documents we scrape
	c.OnResponse(func(r *colly.Response) {
		// Acquire lock to serialize writing to sqlite
		mutex.Lock()
		defer mutex.Unlock()

		// Call extractor callback
		err = cb(r.Request.URL.Path, r.Body)
		if err != nil {
			slog.Debug("extract.crawlHttp could not extract page", "error", err)
			errs = append(errs, err)
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
		slog.Debug("extract.crawlHttp encountered errors", "errors", errs, "error", err)
		return err
	}

	return nil
}
