package source

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ServiceDefinition represents a declared service configuration from IaC.
type ServiceDefinition struct {
	Name        string            `yaml:"name"`
	Image       string            `yaml:"image"`
	Environment map[string]string `yaml:"environment"`
	Ports       []string          `yaml:"ports"`
	Replicas    int               `yaml:"replicas"`
}

// Manifest holds all service definitions loaded from a source file.
type Manifest struct {
	Version  string              `yaml:"version"`
	Services []ServiceDefinition `yaml:"services"`
}

// Loader is responsible for reading IaC source definitions.
type Loader struct {
	basePath string
}

// NewLoader creates a new Loader rooted at basePath.
func NewLoader(basePath string) *Loader {
	return &Loader{basePath: basePath}
}

// LoadManifest reads and parses a YAML manifest file at the given relative path.
func (l *Loader) LoadManifest(relPath string) (*Manifest, error) {
	fullPath := filepath.Join(l.basePath, relPath)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("source: reading manifest %q: %w", fullPath, err)
	}

	var m Manifest
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("source: parsing manifest %q: %w", fullPath, err)
	}

	if err := m.validate(); err != nil {
		return nil, fmt.Errorf("source: invalid manifest %q: %w", fullPath, err)
	}

	return &m, nil
}

func (m *Manifest) validate() error {
	if m.Version == "" {
		return fmt.Errorf("version is required")
	}
	for i, svc := range m.Services {
		if svc.Name == "" {
			return fmt.Errorf("service[%d]: name is required", i)
		}
		if svc.Image == "" {
			return fmt.Errorf("service %q: image is required", svc.Name)
		}
	}
	return nil
}
