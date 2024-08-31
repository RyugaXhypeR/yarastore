package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"yarastore/pkg/utils"
)

// ConfigValues Common configure options to decide the files and directories to scan.
type ConfigValues struct {
	// The directories to consider for scan.
	Dirs []string `toml:"dirs"`
	// The files to consider for scan.
	Files []string `toml:"files"`
	// The files and directories to exclude. Directories must be suffixed with `/`, e.g. `.git/`.
	Exclude []string `toml:"exclude"`
	// Consider files if they have this pattern.
    // Note: Only works against the filename.
	IncludePattern string `toml:"include_pattern"`
	// Discard files if they have this pattern.
    // Note: Only works against the filename.
	ExcludePattern string `toml:"exclude_pattern"`
}

// Config Toml config options for the app.
type Config struct {
	// ConfigValues for yara-rules.
	Rules ConfigValues `toml:"rules"`
	// ConfigValues for the files to match yara-rules against.
	Target ConfigValues `toml:"target"`
}

// LoadConfig Load toml config from a file.
func LoadConfig(filename string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	return &config, err
}

// IsFilenameValid A predicate which applies filters to check given `path` is valid.
func (c *ConfigValues) IsFilenameValid(path string) bool {
	filename := filepath.Base(path)
	if c.IncludePattern != "" {
		match, err := filepath.Match(c.IncludePattern, filename)
		if !match || err != nil {
			return false
		}
	}
	if c.ExcludePattern != "" {
		match, err := filepath.Match(c.ExcludePattern, filename)
		if match || err != nil {
			return false
		}
	}

    // Checks if any component in `c.Exclude` is present in `path`.
	if utils.FileContains(c.Exclude, path) {
		return false
	}
	return true
}
