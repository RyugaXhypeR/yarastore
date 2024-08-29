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

func FileHasPrefix(prefixes []string, filepath string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(filepath, prefix) {
			return true
		}
	}
	return false
}
