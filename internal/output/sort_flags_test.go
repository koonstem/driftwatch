package output

import (
	"testing"

	"github.com/spf13/cobra"
)

func newSortCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindSortFlags(cmd)
	return cmd
}

func TestBindSortFlags_Defaults(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.ParseFlags([]string{})

	opts, err := SortOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.By != SortByService {
		t.Errorf("expected default SortByService, got %q", opts.By)
	}
	if opts.Reverse {
		t.Errorf("expected Reverse=false by default")
	}
}

func TestBindSortFlags_Status(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.ParseFlags([]string{"--sort-by", "status"})

	opts, err := SortOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.By != SortByStatus {
		t.Errorf("expected SortByStatus, got %q", opts.By)
	}
}

func TestBindSortFlags_Reverse(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.ParseFlags([]string{"--sort-reverse"})

	opts, err := SortOptionsFromFlags(cmd)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.Reverse {
		t.Errorf("expected Reverse=true")
	}
}

func TestBindSortFlags_InvalidField(t *testing.T) {
	cmd := newSortCmd()
	_ = cmd.ParseFlags([]string{"--sort-by", "unknown"})

	_, err := SortOptionsFromFlags(cmd)
	if err == nil {
		t.Fatal("expected error for unknown sort field")
	}
}
