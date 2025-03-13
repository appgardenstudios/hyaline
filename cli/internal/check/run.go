package check

import (
	"database/sql"
	"errors"
	"hyaline/internal/config"
	"hyaline/internal/rule"
)

func Run(check config.Check, system string, current *sql.DB, change *sql.DB, recommendAction bool, llmOpts config.LLM) (result *rule.Result, err error) {
	switch check.Rule {
	case rule.SectionExistsRule:
		var options rule.SectionExistsOptions
		options, err = rule.GetSectionExistsOptions(check.Options)
		if err != nil {
			return nil, err
		}
		result, err = rule.RunSectionExists(check.ID, check.Description, options, system, current, change, recommendAction, llmOpts)
	default:
		err = errors.New("unknown rule " + check.Rule)
	}

	return
}
