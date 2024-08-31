package utils

import (
	"io"
	"net/http"
	"strings"
)

// FileContains Return true if any path component in `excludeComponents`
// is present in `filepath`.
func FileContains(excludeComponents []string, filepath string) bool {
	for _, component := range excludeComponents {
		if strings.Contains(filepath, component) {
			return true
		}
	}
	return false
}

// FetchRemoteFile Download a file from a remote server.
func FetchRemoteFile(remoteURL string) ([]byte, error) {
	response, err := http.Get(remoteURL)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(response.Body)
}
