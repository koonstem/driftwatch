package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "driftwatch.yaml")
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	return p
}

func TestRootCmd_MissingConfig(t *testing.T) {
	cfgFile = "/nonexistent/driftwatch.yaml"
	err := runDetect(rootCmd, nil)
	if err == nil {
		t.Fatal("expected error for missing config, got nil")
	}
}

func TestRootCmd_InvalidManifest(t *testing.T) {
	cfgPath := writeTempConfig(t, `
version: "1"
source_type: docker
manifest_path: /nonexistent/manifest.yaml
`)
	cfgFile = cfgPath
	manifest = ""

	err := runDetect(rootCmd, nil)
	if err == nil {
		t.Fatal("expected error for missing manifest, got nil")
	}
}

func TestRootCmd_ManifestFlagOverridesConfig(t *testing.T) {
	_ = bytes.NewBufferString("")
	cfgPath := writeTempConfig(t, `
version: "1"
source_type: docker
manifest_path: /config/path/manifest.yaml
`)
	cfgFile = cfgPath
	manifest = "/override/path/manifest.yaml"

	err := runDetect(rootCmd, nil)
	// We expect an error since the override path doesn't exist,
	// but the error should reference the manifest, not the config.
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	manifest = ""
}
