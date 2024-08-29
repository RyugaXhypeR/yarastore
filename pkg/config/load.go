package config

import (
	"github.com/BurntSushi/toml"
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
