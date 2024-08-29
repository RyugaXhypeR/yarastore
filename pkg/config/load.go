package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"yarastore/pkg/utils"
)

// ConfigValues Common configure options for `Rules` and `Target` files
type ConfigValues struct {
	// The directories to scan.
	Dirs []string `toml:"dirs"`
	// The files to scan.
	Files []string `toml:"files"`
	// The files to exclude.
	Exclude []string `toml:"exclude"`
	// Consider files if they have this pattern.
	IncludePattern string `toml:"include_pattern"`
	// Discard files if they have this pattern.
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
	if err != nil {
		return nil, err
	}
	if err := config.Rules.makePathAbs(); err != nil {
		return nil, err
	}
	if err := config.Target.makePathAbs(); err != nil {
		return nil, err
	}
	return &config, nil
}

// makePathAbs Convert all patterns and exclude files to absolute paths
// to make matching patterns easier.
func (c *ConfigValues) makePathAbs() error {
	// Convert all files in `Exclude` to absolute paths.
	for i, file := range c.Exclude {
		absPath, err := filepath.Abs(file)
		if err != nil {
			return err
		}
		c.Exclude[i] = absPath
	}

	if c.IncludePattern != "" {
		absPath, err := filepath.Abs(c.IncludePattern)
		if err != nil {
			return err
		}
		c.IncludePattern = absPath
	}

	if c.ExcludePattern != "" {
		absPath, err := filepath.Abs(c.ExcludePattern)
		if err != nil {
			return err
		}
		c.ExcludePattern = absPath
	}

	return nil
}

// IsFilenameValid A predicate which applies filters to check given `path` is valid.
func (c *ConfigValues) IsFilenameValid(path string) bool {
	if c.IncludePattern != "" {
		match, err := filepath.Match(c.IncludePattern, path)
		if !match || err != nil {
			return false
		}
	}

	if c.ExcludePattern != "" {
		match, err := filepath.Match(c.ExcludePattern, path)
		if match || err != nil {
			return false
		}
	}

	if utils.FileHasPrefix(c.Exclude, path) {
		return false
	}

	return true
}
