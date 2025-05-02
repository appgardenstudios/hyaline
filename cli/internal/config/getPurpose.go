package config

func GetPurpose(systemID string, documentationSourceID string, documentID string, section []string, cfg *Config) (purpose string, found bool) {
	// Get system
	system, found := cfg.GetSystem(systemID)
	if !found {
		return
	}

	// Get documentationSource
	documentationSource, found := system.GetDocumentationSource(documentationSourceID)
	if !found {
		return
	}

	// Get ruleDoc
	ruleDoc, found := GetRuleDocument(cfg.Rules, documentationSource.Rules, documentID)
	if !found {
		return
	}

	// If this is a document, return
	if len(section) == 0 {
		return ruleDoc.Purpose, true
	}

	// Recurse through sections
	// var sectionPtr *RuleDocumentSection
	// for {

	// }
	// TODO

	return
}
