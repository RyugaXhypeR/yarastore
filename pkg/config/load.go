package config

import (
	"os"
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

func (c *ConfigValues) GetFilteredFiles() ([]string, error) {
	var filteredFiles []string

	predicate := func(path string, _ os.FileInfo) bool {
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

	for _, dir := range c.Dirs {
		files, err := utils.ListDirWithPred(dir, predicate)
		if err != nil {
			return nil, err
		}
		filteredFiles = append(filteredFiles, files...)
	}

	for _, file := range c.Files {
		if predicate(file, nil) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles, nil
}

func (c *Config) GetRuleFiles() ([]string, error) {
	return c.Rules.GetFilteredFiles()
}

func (c *Config) GetTargetFiles() ([]string, error) {
	return c.Target.GetFilteredFiles()
}
