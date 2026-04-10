package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/driftwatch/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestLoad_ValidConfig(t *testing.T) {
	raw := `
version: "1"
defaults:
  interval: 30s
  timeout: 10s
sources:
  - name: prod-k8s
    type: kubernetes
    path: ./manifests
    labels:
      env: production
`
	p := writeTemp(t, raw)
	cfg, err := config.Load(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Version != "1" {
		t.Errorf("expected version 1, got %q", cfg.Version)
	}
	if len(cfg.Sources) != 1 {
		t.Fatalf("expected 1 source, got %d", len(cfg.Sources))
	}
	if cfg.Sources[0].Name != "prod-k8s" {
		t.Errorf("unexpected source name: %q", cfg.Sources[0].Name)
	}
}

func TestLoad_MissingVersion(t *testing.T) {
	raw := `
sources:
  - name: infra
    type: terraform
    path: ./tf
`
	p := writeTemp(t, raw)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing version, got nil")
	}
}

func TestLoad_MissingSourceType(t *testing.T) {
	raw := `
version: "1"
sources:
  - name: infra
    path: ./tf
`
	p := writeTemp(t, raw)
	_, err := config.Load(p)
	if err == nil {
		t.Fatal("expected error for missing source type, got nil")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := config.Load("/nonexistent/path/driftwatch.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
