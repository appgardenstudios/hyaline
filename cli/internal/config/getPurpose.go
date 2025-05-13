package config

func GetPurpose(systemID string, documentationSourceID string, documentID string, sectionIDs []string, cfg *Config) (purpose string, found bool) {
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

	// Get desired document
	desiredDoc, found := documentationSource.GetDocument(cfg, documentID)
	if !found {
		return
	}

	// If this is a document, return
	if len(sectionIDs) == 0 {
		return desiredDoc.Purpose, true
	}

	// Otherwise, get purpose from section
	return getPurposeFromSection(desiredDoc.Sections, sectionIDs)
}

func getPurposeFromSection(sections []DocumentSection, sectionIDs []string) (purpose string, found bool) {
	var currentID string
	var remainder []string

	switch len(sectionIDs) {
	case 0:
		return
	case 1:
		currentID = sectionIDs[0]
		remainder = []string{}
	default:
		currentID = sectionIDs[0]
		remainder = sectionIDs[1:]
	}

	for _, sec := range sections {
		if sec.Name == currentID {
			if len(remainder) == 0 {
				return sec.Purpose, true
			} else {
				return getPurposeFromSection(sec.Sections, remainder)
			}
		}
	}

	return
}
