package utils

import (
	"os"
	"path/filepath"
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

func SliceContains(haystack []string, needle string) bool {
	for _, straw := range haystack {
		if straw == needle {
			return true
		}
	}
	return false
}
