package code

import "github.com/bmatcuk/doublestar/v4"

func isIncluded(name string, includes []string, excludes []string) bool {
	for _, include := range includes {
		if doublestar.MatchUnvalidated(include, name) {
			for _, exclude := range excludes {
				if doublestar.MatchUnvalidated(exclude, name) {
					return false
				}
			}
			return true
		}
	}

	return false
}
