package source

// ServiceSpec describes the desired state of a single service
// as declared in the manifest file.
type ServiceSpec struct {
	Name   string `yaml:"name"`
	Image  string `yaml:"image"`
	Status string `yaml:"status"`
}

// Manifest is the top-level structure parsed from a driftwatch manifest YAML.
type Manifest struct {
	Version  string        `yaml:"version"`
	Services []ServiceSpec `yaml:"services"`
}
