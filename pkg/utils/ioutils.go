package utils

import (
	"strings"
)

// FileContains Return true if any path component in `excludeComponents`
// is present in `filepath`
func FileContains(excludeComponents []string, filepath string) bool {
	for _, component := range excludeComponents {
		if strings.Contains(filepath, component) {
			return true
		}
	}
	return false
}
