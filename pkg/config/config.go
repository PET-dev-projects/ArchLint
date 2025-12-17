package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/PET-dev-projects/ArchLint/pkg/engine"
)

// File describes the YAML configuration file layout.
type File struct {
	Rules []RuleEntry `yaml:"rules"`
}

// RuleEntry declares an individual rule override.
type RuleEntry struct {
	ID      string                 `yaml:"id"`
	Enabled *bool                  `yaml:"enabled"`
	Config  map[string]any         `yaml:"config"`
	Meta    map[string]interface{} `yaml:"-"`
}

// LoadOptionsFromFile parses the YAML config file into engine.Options.
func LoadOptionsFromFile(path string) (engine.Options, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return engine.Options{}, err
	}
	var file File
	if err := yaml.Unmarshal(data, &file); err != nil {
		return engine.Options{}, err
	}
	opts := engine.Options{
		RuleConfig:   map[string]map[string]any{},
		EnabledRules: []string{},
	}
	for idx, entry := range file.Rules {
		if entry.ID == "" {
			return engine.Options{}, fmt.Errorf("config rules[%d]: id is required", idx)
		}
		if entry.Enabled != nil && !*entry.Enabled {
			continue
		}
		opts.EnabledRules = append(opts.EnabledRules, entry.ID)
		if entry.Config != nil {
			opts.RuleConfig[entry.ID] = entry.Config
		}
	}
	if len(opts.RuleConfig) == 0 {
		opts.RuleConfig = nil
	}
	if len(opts.EnabledRules) == 0 {
		opts.EnabledRules = nil
	}
	return opts, nil
}
