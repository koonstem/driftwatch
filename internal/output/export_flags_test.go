package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newExportCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindExportFlags(cmd)
	return cmd
}

func TestBindExportFlags_Defaults(t *testing.T) {
	cmd := newExportCmd()
	_ = cmd.ParseFlags([]string{})
	opts, path, err := ExportOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "csv" {
		t.Errorf("Format = %q, want \"csv\"", opts.Format)
	}
	if !opts.Timestamp {
		t.Error("Timestamp should default to true")
	}
	if path != "" {
		t.Errorf("path = %q, want empty", path)
	}
}

func TestBindExportFlags_JSONFormat(t *testing.T) {
	cmd := newExportCmd()
	_ = cmd.ParseFlags([]string{"--export-format", "json"})
	opts, _, err := ExportOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Format != "json" {
		t.Errorf("Format = %q, want \"json\"", opts.Format)
	}
}

func TestBindExportFlags_NoTimestamp(t *testing.T) {
	cmd := newExportCmd()
	_ = cmd.ParseFlags([]string{"--export-timestamp=false"})
	opts, _, err := ExportOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.Timestamp {
		t.Error("Timestamp should be false")
	}
}

func TestBindExportFlags_InvalidFormat(t *testing.T) {
	cmd := newExportCmd()
	_ = cmd.ParseFlags([]string{"--export-format", "xml"})
	_, _, err := ExportOptionsFromFlags(cmd)
	if err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestBindExportFlags_Path(t *testing.T) {
	cmd := newExportCmd()
	_ = cmd.ParseFlags([]string{"--export-path", "/tmp/out.csv"})
	_, path, err := ExportOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if path != "/tmp/out.csv" {
		t.Errorf("path = %q, want \"/tmp/out.csv\"", path)
	}
}
