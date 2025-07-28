package config

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/bmatcuk/doublestar/v4"
)

const extractSourceIDRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const extractMetadataTagKeyRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`
const extractMetadataTagValueRegex = `^[A-z0-9][A-z0-9_-]{0,63}$`

func validateExtract(cfg *Config) (err error) {
	// If extract was not defined in the config check nothing as extract is not always required
	// Note that actions requiring the config need to check for nil themselves
	if cfg.Extract == nil {
		return
	}

	// Check source
	if !regexp.MustCompile(extractSourceIDRegex).MatchString(cfg.Extract.Source.ID) {
		return fmt.Errorf("extract.source.id must match regex /%s/, found: %s", extractSourceIDRegex, cfg.Extract.Source.ID)
	}

	// Check crawler
	if !cfg.Extract.Crawler.Type.IsValid() {
		return fmt.Errorf("extract.crawler.type must be one of %s, found: %s", cfg.Extract.Crawler.Type.PossibleValues(), cfg.Extract.Crawler.Type)
	}
	for i, include := range cfg.Extract.Crawler.Include {
		if include == "" || !doublestar.ValidatePattern(include) {
			return fmt.Errorf("extract.crawler.include[%d] must be a valid pattern, found: %s", i, include)
		}
	}
	for i, exclude := range cfg.Extract.Crawler.Exclude {
		if exclude == "" || !doublestar.ValidatePattern(exclude) {
			return fmt.Errorf("extract.crawler.exclude[%d] must be a valid pattern, found: %s", i, exclude)
		}
	}

	// Check extractors
	if len(cfg.Extract.Extractors) == 0 {
		return errors.New("extract.extractors must contain at least one extractor, none found")
	}
	for i, extractor := range cfg.Extract.Extractors {
		if !extractor.Type.IsValid() {
			return fmt.Errorf("extract.extractors[%d].type must be one of %s, found: %s", i, extractor.Type.PossibleValues(), extractor.Type)
		}
		for j, include := range extractor.Include {
			if include == "" || !doublestar.ValidatePattern(include) {
				return fmt.Errorf("extract.extractors[%d].include[%d] must be a valid pattern, found: %s", i, j, include)
			}
		}
		for j, exclude := range extractor.Exclude {
			if exclude == "" || !doublestar.ValidatePattern(exclude) {
				return fmt.Errorf("extract.extractors[%d].exclude[%d] must be a valid pattern, found: %s", i, j, exclude)
			}
		}
	}

	// Check metadata
	keyRegex := regexp.MustCompile(extractMetadataTagKeyRegex)
	valueRegex := regexp.MustCompile(extractMetadataTagValueRegex)
	for i, metadata := range cfg.Extract.Metadata {
		if metadata.Document == "" || !doublestar.ValidatePattern(metadata.Document) {
			return fmt.Errorf("extract.metadata[%d].document must be a valid pattern, found: %s", i, metadata.Document)
		}
		if metadata.Section != "" && !doublestar.ValidatePattern(metadata.Section) {
			return fmt.Errorf("extract.metadata[%d].section must be a valid pattern if not empty, found: %s", i, metadata.Section)
		}
		for j, tag := range metadata.Tags {
			if !keyRegex.MatchString(tag.Key) {
				return fmt.Errorf("extract.metadata[%d].tags[%d].key must match regex /%s/, found: %s", i, j, extractMetadataTagKeyRegex, tag.Key)
			}
			if !valueRegex.MatchString(tag.Value) {
				return fmt.Errorf("extract.metadata[%d].tags[%d].value must match regex /%s/, found: %s", i, j, extractMetadataTagValueRegex, tag.Value)
			}
		}
	}

	return
}
