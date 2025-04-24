package config

// Note: If there are multiple ruleDocuments in different rule sets, this function
// will return the last one found based on the order of the passed in rule IDs.
func GetRuleDocument(rules []Rule, ruleIDs []string, path string) (found bool, ruleDoc RuleDocument) {
	for i := len(ruleIDs) - 1; i >= 0; i-- {
		rule := getRule(rules, ruleIDs[i])
		if rule != nil {
			for _, doc := range rule.Documents {
				if doc.Path == path {
					return true, doc
				}
			}
		}
	}

	return
}

// Return the rule identified by the id, or nil if not found
func getRule(rules []Rule, ruleID string) *Rule {
	for _, rule := range rules {
		if rule.ID == ruleID {
			return &rule
		}
	}

	return nil
}
