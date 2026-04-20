package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newMaskCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindMaskFlags(cmd)
	return cmd
}

func TestBindMaskFlags_Defaults(t *testing.T) {
	cmd := newMaskCmd()
	opts := MaskOptionsFromFlags(cmd)
	if opts.Enabled {
		t.Error("expected mask disabled by default")
	}
	if opts.MaskChar != "*" {
		t.Errorf("expected default mask char '*', got %q", opts.MaskChar)
	}
	if opts.MaskLength != 6 {
		t.Errorf("expected default mask length 6, got %d", opts.MaskLength)
	}
	if len(opts.Fields) != 1 || opts.Fields[0] != "image" {
		t.Errorf("expected default fields [image], got %v", opts.Fields)
	}
}

func TestBindMaskFlags_Enabled(t *testing.T) {
	cmd := newMaskCmd()
	_ = cmd.Flags().Set("mask", "true")
	opts := MaskOptionsFromFlags(cmd)
	if !opts.Enabled {
		t.Error("expected mask to be enabled")
	}
}

func TestBindMaskFlags_CustomChar(t *testing.T) {
	cmd := newMaskCmd()
	_ = cmd.Flags().Set("mask-char", "#")
	opts := MaskOptionsFromFlags(cmd)
	if opts.MaskChar != "#" {
		t.Errorf("expected mask char '#', got %q", opts.MaskChar)
	}
}

func TestBindMaskFlags_CustomLength(t *testing.T) {
	cmd := newMaskCmd()
	_ = cmd.Flags().Set("mask-length", "10")
	opts := MaskOptionsFromFlags(cmd)
	if opts.MaskLength != 10 {
		t.Errorf("expected mask length 10, got %d", opts.MaskLength)
	}
}

func TestBindMaskFlags_CustomFields(t *testing.T) {
	cmd := newMaskCmd()
	_ = cmd.Flags().Set("mask-fields", "image,env")
	opts := MaskOptionsFromFlags(cmd)
	if len(opts.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(opts.Fields))
	}
	if opts.Fields[0] != "image" || opts.Fields[1] != "env" {
		t.Errorf("unexpected fields: %v", opts.Fields)
	}
}

func TestDefaultMaskOptions_Values(t *testing.T) {
	opts := DefaultMaskOptions()
	if opts.MaskChar == "" {
		t.Error("default mask char should not be empty")
	}
	if opts.MaskLength <= 0 {
		t.Error("default mask length should be positive")
	}
}
