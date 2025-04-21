package config

import "github.com/bmatcuk/doublestar/v4"

func PathIsIncluded(path string, includes []string, excludes []string) bool {
	for _, include := range includes {
		if doublestar.MatchUnvalidated(include, path) {
			for _, exclude := range excludes {
				if doublestar.MatchUnvalidated(exclude, path) {
					return false
				}
			}
			return true
		}
	}

	return false
}
