package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// LayerConfig defines a named layer and the packages it is forbidden to import.
type LayerConfig struct {
	Name   string   `yaml:"name"`
	Path   string   `yaml:"path"`
	Forbid []string `yaml:"forbid"`
}

// Config is the root configuration parsed from clinicius.yaml.
type Config struct {
	Layers []LayerConfig `yaml:"layers"`
}

// Load reads and parses a YAML configuration file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	return &cfg, nil
}
