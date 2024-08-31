package utils

import (
	"strings"
)

func FileContains(excludeComponents []string, filepath string) bool {
	for _, component := range excludeComponents {
		if strings.Contains(filepath, component) {
			return true
		}
	}
	return false
}
