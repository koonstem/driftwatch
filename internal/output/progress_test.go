package output

import (
	"bytes"
	"strings"
	"testing"
)

func newTestProgress(verbose bool, colorEnabled bool) (*ProgressWriter, *bytes.Buffer) {
	buf := &bytes.Buffer{}
	c := NewColorizer(colorEnabled)
	return NewProgressWriter(buf, verbose, c), buf
}

func TestProgressWriter_StepStart_Verbose(t *testing.T) {
	pw, buf := newTestProgress(true, false)
	pw.StepStart("loading manifest")
	if !strings.Contains(buf.String(), "loading manifest") {
		t.Errorf("expected step name in output, got: %q", buf.String())
	}
}

func TestProgressWriter_StepStart_Silent(t *testing.T) {
	pw, buf := newTestProgress(false, false)
	pw.StepStart("loading manifest")
	if buf.Len() != 0 {
		t.Errorf("expected no output in non-verbose mode, got: %q", buf.String())
	}
}

func TestProgressWriter_StepDone_Verbose(t *testing.T) {
	pw, buf := newTestProgress(true, false)
	pw.StepDone("loading manifest")
	if !strings.Contains(buf.String(), "loading manifest") {
		t.Errorf("expected step name in done output, got: %q", buf.String())
	}
}

func TestProgressWriter_StepWarn_WithReason(t *testing.T) {
	pw, buf := newTestProgress(true, false)
	pw.StepWarn("container check", "timeout")
	out := buf.String()
	if !strings.Contains(out, "container check") {
		t.Errorf("expected step name in warn output, got: %q", out)
	}
	if !strings.Contains(out, "timeout") {
		t.Errorf("expected reason in warn output, got: %q", out)
	}
}

func TestProgressWriter_StepWarn_NoReason(t *testing.T) {
	pw, buf := newTestProgress(true, false)
	pw.StepWarn("container check", "")
	out := buf.String()
	if !strings.Contains(out, "container check") {
		t.Errorf("expected step name in warn output, got: %q", out)
	}
}

func TestProgressWriter_Summary_NoDrift(t *testing.T) {
	pw, buf := newTestProgress(false, false)
	pw.Summary(5, 0)
	out := buf.String()
	if !strings.Contains(out, "5") || !strings.Contains(out, "0") {
		t.Errorf("expected counts in summary, got: %q", out)
	}
}

func TestProgressWriter_Summary_WithDrift(t *testing.T) {
	pw, buf := newTestProgress(false, false)
	pw.Summary(5, 3)
	out := buf.String()
	if !strings.Contains(out, "5") || !strings.Contains(out, "3") {
		t.Errorf("expected counts in summary, got: %q", out)
	}
}
