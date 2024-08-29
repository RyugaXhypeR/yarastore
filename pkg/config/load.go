package config

import (
	"path/filepath"

	"github.com/BurntSushi/toml"

	"yarastore/pkg/utils"
)

type ConfigValues struct {
	Dirs           []string `toml:"dirs"`
	Files          []string `toml:"files"`
	Exclude        []string `toml:"exclude"`
	IncludePattern string   `toml:"include_pattern"`
	ExcludePattern string   `toml:"exclude_pattern"`
}

type Config struct {
	Rules  ConfigValues `toml:"rules"`
	Target ConfigValues `toml:"target"`
}

func LoadConfig(filename string) (*Config, error) {
	var config Config
	_, err := toml.DecodeFile(filename, &config)
	return &config, err
}

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

	if utils.SliceContains(c.Exclude, filename) {
		return false
	}

	return true
}
