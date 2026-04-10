package drift

import (
	"fmt"

	"github.com/yourorg/driftwatch/internal/runner"
	"github.com/yourorg/driftwatch/internal/source"
)

// Result holds the drift comparison result for a single service.
type Result struct {
	ServiceName string
	Expected    source.ServiceSpec
	Actual      runner.ContainerInfo
	Drifted     bool
	Reasons     []string
}

// Detector compares running containers against declared service specs.
type Detector struct{}

// NewDetector returns a new Detector instance.
func NewDetector() *Detector {
	return &Detector{}
}

// Detect compares a list of running containers against the declared manifest.
// It returns one Result per declared service.
func (d *Detector) Detect(manifest *source.Manifest, containers []runner.ContainerInfo) []Result {
	results := make([]Result, 0, len(manifest.Services))

	containerMap := indexContainers(containers)

	for _, svc := range manifest.Services {
		result := Result{
			ServiceName: svc.Name,
			Expected:    svc,
		}

		container, found := containerMap[svc.Name]
		if !found {
			result.Drifted = true
			result.Reasons = append(result.Reasons, fmt.Sprintf("service %q not found in running containers", svc.Name))
			results = append(results, result)
			continue
		}

		result.Actual = container
		result.Reasons = compareService(svc, container)
		result.Drifted = len(result.Reasons) > 0
		results = append(results, result)
	}

	return results
}

func indexContainers(containers []runner.ContainerInfo) map[string]runner.ContainerInfo {
	m := make(map[string]runner.ContainerInfo, len(containers))
	for _, c := range containers {
		m[c.Name] = c
	}
	return m
}

func compareService(svc source.ServiceSpec, c runner.ContainerInfo) []string {
	var reasons []string

	if svc.Image != "" && svc.Image != c.Image {
		reasons = append(reasons, fmt.Sprintf("image mismatch: expected %q, got %q", svc.Image, c.Image))
	}

	if svc.Status != "" && svc.Status != c.Status {
		reasons = append(reasons, fmt.Sprintf("status mismatch: expected %q, got %q", svc.Status, c.Status))
	}

	return reasons
}
