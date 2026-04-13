package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newTemplateCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindTemplateFlags(cmd)
	return cmd
}

func TestBindTemplateFlags_Defaults(t *testing.T) {
	cmd := newTemplateCmd()
	opts, err := TemplateOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.TemplatePath != "" {
		t.Errorf("expected empty TemplatePath, got %q", opts.TemplatePath)
	}
	if opts.TemplateStr != "" {
		t.Errorf("expected empty TemplateStr, got %q", opts.TemplateStr)
	}
}

func TestBindTemplateFlags_InlineTemplate(t *testing.T) {
	cmd := newTemplateCmd()
	_ = cmd.Flags().Set("template", "hello={{.Total}}")
	opts, err := TemplateOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.TemplateStr != "hello={{.Total}}" {
		t.Errorf("unexpected TemplateStr: %q", opts.TemplateStr)
	}
}

func TestBindTemplateFlags_FileTemplate(t *testing.T) {
	cmd := newTemplateCmd()
	_ = cmd.Flags().Set("template-file", "/tmp/my.tmpl")
	opts, err := TemplateOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.TemplatePath != "/tmp/my.tmpl" {
		t.Errorf("unexpected TemplatePath: %q", opts.TemplatePath)
	}
}

func TestBindTemplateFlags_MutuallyExclusive(t *testing.T) {
	cmd := newTemplateCmd()
	_ = cmd.Flags().Set("template", "hello")
	_ = cmd.Flags().Set("template-file", "/tmp/tmpl")
	_, err := TemplateOptionsFromFlags(cmd)
	if err == nil {
		t.Fatal("expected error for mutually exclusive flags")
	}
}

func TestIsTemplateEnabled_False(t *testing.T) {
	cmd := newTemplateCmd()
	if IsTemplateEnabled(cmd) {
		t.Error("expected IsTemplateEnabled to be false with no flags set")
	}
}

func TestIsTemplateEnabled_True(t *testing.T) {
	cmd := newTemplateCmd()
	_ = cmd.Flags().Set("template", "{{.Total}}")
	if !IsTemplateEnabled(cmd) {
		t.Error("expected IsTemplateEnabled to be true")
	}
}
