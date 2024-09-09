package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"
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
	// Flag to match rules against files recursively
	Recursive bool `toml:"recursive"`

	// Internal field for faster comparison of filenames
	excludeFiles map[string]bool
	// Internal field for faster comparison of directories
	excludeDirs map[string]bool
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

	config.Rules.makeExcludeMaps()
	config.Target.makeExcludeMaps()

	return &config, err
}

// IsFilenameValid A predicate which applies filters to check given `path` is valid.
func (c *ConfigValues) IsFilenameValid(filename string) bool {
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

    // Valid if filename doesn't exist in `c.excludeFiles`
	return !c.excludeFiles[filename]
}

// makeExcludeMaps Seperates directories and files from `c.Exclude` into `c.excludeDirs` and `c.excludeFiles`.`
func (c *ConfigValues) makeExcludeMaps() {
	c.excludeFiles = map[string]bool{}
	c.excludeDirs = map[string]bool{}

	for _, component := range c.Exclude {
		if component[len(component)-1:] == "/" {
            c.excludeDirs[component[:len(component)-1]] = true
		} else {
			c.excludeFiles[component] = true
		}
	}
}

// IsDirExcluded Check if `dirname` is supposed to be excluded.
func (c *ConfigValues) IsDirExcluded(dirname string) bool {
	return c.excludeDirs[dirname]
}
