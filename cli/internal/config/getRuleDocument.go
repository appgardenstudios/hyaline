package config

import "slices"

func GetRuleDocument(rules []Rule, ruleIDs []string, path string) (found bool, ruleDoc RuleDocument) {
	for _, rule := range rules {
		if slices.Contains(ruleIDs, rule.ID) {
			for _, doc := range rule.Documents {
				if doc.Path == path {
					return true, doc
				}
			}
		}
	}

	return
}
