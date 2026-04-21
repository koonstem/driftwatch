package output

import (
	"testing"
	"time"

	"github.com/spf13/cobra"
)

func newCircuitCmd() *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	BindCircuitFlags(cmd)
	return cmd
}

func TestBindCircuitFlags_Defaults(t *testing.T) {
	cmd := newCircuitCmd()
	_ = cmd.ParseFlags([]string{})
	opts := CircuitOptionsFromFlags(cmd)

	if opts.Enabled {
		t.Error("expected circuit disabled by default")
	}
	if opts.MaxFailures != 3 {
		t.Errorf("expected MaxFailures=3, got %d", opts.MaxFailures)
	}
	if opts.ResetTimeout != 30*time.Second {
		t.Errorf("expected ResetTimeout=30s, got %s", opts.ResetTimeout)
	}
}

func TestBindCircuitFlags_Enabled(t *testing.T) {
	cmd := newCircuitCmd()
	_ = cmd.ParseFlags([]string{"--circuit-enabled"})
	opts := CircuitOptionsFromFlags(cmd)

	if !opts.Enabled {
		t.Error("expected circuit enabled")
	}
}

func TestBindCircuitFlags_CustomValues(t *testing.T) {
	cmd := newCircuitCmd()
	_ = cmd.ParseFlags([]string{
		"--circuit-enabled",
		"--circuit-max-failures=5",
		"--circuit-reset-timeout=1m",
	})
	opts := CircuitOptionsFromFlags(cmd)

	if opts.MaxFailures != 5 {
		t.Errorf("expected MaxFailures=5, got %d", opts.MaxFailures)
	}
	if opts.ResetTimeout != time.Minute {
		t.Errorf("expected ResetTimeout=1m, got %s", opts.ResetTimeout)
	}
}

func TestIsCircuitEnabled_False(t *testing.T) {
	cmd := newCircuitCmd()
	_ = cmd.ParseFlags([]string{})
	if IsCircuitEnabled(cmd) {
		t.Error("expected circuit to be disabled")
	}
}

func TestIsCircuitEnabled_True(t *testing.T) {
	cmd := newCircuitCmd()
	_ = cmd.ParseFlags([]string{"--circuit-enabled"})
	if !IsCircuitEnabled(cmd) {
		t.Error("expected circuit to be enabled")
	}
}

func TestDefaultCircuitOptions_Values(t *testing.T) {
	opts := DefaultCircuitOptions()
	if opts.Enabled {
		t.Error("default should be disabled")
	}
	if opts.MaxFailures <= 0 {
		t.Error("MaxFailures must be positive")
	}
	if opts.ResetTimeout <= 0 {
		t.Error("ResetTimeout must be positive")
	}
}
