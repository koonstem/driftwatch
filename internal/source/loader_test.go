package source_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/driftwatch/driftwatch/internal/source"
)

func writeTempManifest(t *testing.T, content string) (dir string, filename string) {
	t.Helper()
	dir = t.TempDir()
	filename = "manifest.yaml"
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp manifest: %v", err)
	}
	return dir, filename
}

func TestLoadManifest_Valid(t *testing.T) {
	content := `version: "1.0"
services:
  - name: api
    image: myapp:latest
    replicas: 2
    ports:
      - "8080:8080"
    environment:
      ENV: production
`
	dir, file := writeTempManifest(t, content)
	loader := source.NewLoader(dir)

	m, err := loader.LoadManifest(file)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(m.Services) != 1 {
		t.Fatalf("expected 1 service, got %d", len(m.Services))
	}
	if m.Services[0].Name != "api" {
		t.Errorf("expected service name 'api', got %q", m.Services[0].Name)
	}
	if m.Services[0].Replicas != 2 {
		t.Errorf("expected replicas 2, got %d", m.Services[0].Replicas)
	}
}

func TestLoadManifest_MissingVersion(t *testing.T) {
	content := `services:
  - name: api
    image: myapp:latest
`
	dir, file := writeTempManifest(t, content)
	loader := source.NewLoader(dir)

	_, err := loader.LoadManifest(file)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestLoadManifest_MissingServiceImage(t *testing.T) {
	content := `version: "1.0"
services:
  - name: worker
`
	dir, file := writeTempManifest(t, content)
	loader := source.NewLoader(dir)

	_, err := loader.LoadManifest(file)
	if err == nil {
		t.Fatal("expected error for missing image, got nil")
	}
}

func TestLoadManifest_FileNotFound(t *testing.T) {
	loader := source.NewLoader(t.TempDir())
	_, err := loader.LoadManifest("nonexistent.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
