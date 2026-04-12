package filter_test

import (
	"testing"

	"github.com/driftwatch/internal/filter"
	"github.com/spf13/cobra"
)

func newTestCmd(args []string) (*cobra.Command, func() filter.Options) {
	cmd := &cobra.Command{Use: "test", RunE: func(cmd *cobra.Command, args []string) error { return nil }}
	getOpts := filter.BindFlags(cmd)
	cmd.SetArgs(args)
	return cmd, getOpts
}

func TestBindFlags_Defaults(t *testing.T) {
	cmd, getOpts := newTestCmd([]string{})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	opts := getOpts()
	if len(opts.Services) != 0 {
		t.Errorf("expected no services, got %v", opts.Services)
	}
	if opts.OnlyDrifted {
		t.Error("expected OnlyDrifted=false by default")
	}
	if opts.LabelSelector != "" {
		t.Errorf("expected empty label selector, got %q", opts.LabelSelector)
	}
}

func TestBindFlags_Services(t *testing.T) {
	cmd, getOpts := newTestCmd([]string{"--services", "api, worker"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	opts := getOpts()
	if len(opts.Services) != 2 {
		t.Fatalf("expected 2 services, got %d", len(opts.Services))
	}
	if opts.Services[0] != "api" || opts.Services[1] != "worker" {
		t.Errorf("unexpected services: %v", opts.Services)
	}
}

func TestBindFlags_OnlyDrifted(t *testing.T) {
	cmd, getOpts := newTestCmd([]string{"--only-drifted"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if !getOpts().OnlyDrifted {
		t.Error("expected OnlyDrifted=true")
	}
}

func TestBindFlags_LabelSelector(t *testing.T) {
	cmd, getOpts := newTestCmd([]string{"--label", "image=nginx"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
	if getOpts().LabelSelector != "image=nginx" {
		t.Errorf("unexpected label selector: %q", getOpts().LabelSelector)
	}
}
