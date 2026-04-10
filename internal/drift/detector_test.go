package drift_test

import (
	"testing"

	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/runner"
	"github.com/yourorg/driftwatch/internal/source"
)

func manifest(services ...source.ServiceSpec) *source.Manifest {
	return &source.Manifest{Version: "1", Services: services}
}

func svc(name, image, status string) source.ServiceSpec {
	return source.ServiceSpec{Name: name, Image: image, Status: status}
}

func container(name, image, status string) runner.ContainerInfo {
	return runner.ContainerInfo{Name: name, Image: image, Status: status}
}

func TestDetect_NoDrift(t *testing.T) {
	d := drift.NewDetector()
	results := d.Detect(
		manifest(svc("api", "nginx:1.25", "running")),
		[]runner.ContainerInfo{container("api", "nginx:1.25", "running")},
	)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Drifted {
		t.Errorf("expected no drift, got reasons: %v", results[0].Reasons)
	}
}

func TestDetect_ImageMismatch(t *testing.T) {
	d := drift.NewDetector()
	results := d.Detect(
		manifest(svc("api", "nginx:1.25", "running")),
		[]runner.ContainerInfo{container("api", "nginx:1.24", "running")},
	)
	if !results[0].Drifted {
		t.Fatal("expected drift due to image mismatch")
	}
	if len(results[0].Reasons) != 1 {
		t.Errorf("expected 1 reason, got %d", len(results[0].Reasons))
	}
}

func TestDetect_MissingContainer(t *testing.T) {
	d := drift.NewDetector()
	results := d.Detect(
		manifest(svc("worker", "myapp:latest", "running")),
		[]runner.ContainerInfo{},
	)
	if !results[0].Drifted {
		t.Fatal("expected drift for missing container")
	}
}

func TestDetect_StatusMismatch(t *testing.T) {
	d := drift.NewDetector()
	results := d.Detect(
		manifest(svc("db", "postgres:15", "running")),
		[]runner.ContainerInfo{container("db", "postgres:15", "exited")},
	)
	if !results[0].Drifted {
		t.Fatal("expected drift due to status mismatch")
	}
}

func TestDetect_MultipleServices(t *testing.T) {
	d := drift.NewDetector()
	results := d.Detect(
		manifest(
			svc("api", "nginx:1.25", "running"),
			svc("db", "postgres:15", "running"),
		),
		[]runner.ContainerInfo{
			container("api", "nginx:1.25", "running"),
			container("db", "postgres:14", "running"),
		},
	)
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].Drifted {
		t.Error("api should not have drifted")
	}
	if !results[1].Drifted {
		t.Error("db should have drifted")
	}
}
