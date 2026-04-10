package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config holds the top-level driftwatch configuration.
type Config struct {
	Version  string          `yaml:"version"`
	Defaults DefaultSettings `yaml:"defaults"`
	Sources  []SourceConfig  `yaml:"sources"`
}

// DefaultSettings contains global defaults applied to all checks.
type DefaultSettings struct {
	Interval time.Duration `yaml:"interval"`
	Timeout  time.Duration `yaml:"timeout"`
}

// SourceConfig describes a single IaC source to compare against.
type SourceConfig struct {
	Name    string            `yaml:"name"`
	Type    string            `yaml:"type"` // e.g. "terraform", "kubernetes"
	Path    string            `yaml:"path"`
	Labels  map[string]string `yaml:"labels,omitempty"`
}

// Load reads and parses a YAML config file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate performs basic sanity checks on the loaded config.
func (c *Config) validate() error {
	if c.Version == "" {
		return fmt.Errorf("version field is required")
	}
	for i, s := range c.Sources {
		if s.Name == "" {
			return fmt.Errorf("source[%d]: name is required", i)
		}
		if s.Type == "" {
			return fmt.Errorf("source %q: type is required", s.Name)
		}
		if s.Path == "" {
			return fmt.Errorf("source %q: path is required", s.Name)
		}
	}
	return nil
}
