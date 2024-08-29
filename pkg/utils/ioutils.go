package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func ListDirWithPred(dirname string, predicate func(string, os.FileInfo) bool) ([]string, error) {
	var files []string
	err := filepath.Walk(dirname, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && predicate(path, info) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func FileContains(excludeComponents []string, filepath string) bool {
	for _, component := range excludeComponents {
		if strings.Contains(filepath, component) {
			return true
		}
	}
	return false
}
