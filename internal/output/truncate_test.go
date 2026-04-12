package output

import (
	"strings"
	"testing"
)

func newTestTruncator(max int) *Truncator {
	return NewTruncator(TruncateOptions{MaxLength: max, Suffix: "..."})
}

func TestTruncate_ShortString(t *testing.T) {
	tr := newTestTruncator(20)
	out := tr.Truncate("hello")
	if out != "hello" {
		t.Fatalf("expected 'hello', got %q", out)
	}
}

func TestTruncate_ExactLength(t *testing.T) {
	tr := newTestTruncator(5)
	out := tr.Truncate("hello")
	if out != "hello" {
		t.Fatalf("expected 'hello', got %q", out)
	}
}

func TestTruncate_LongString(t *testing.T) {
	tr := newTestTruncator(10)
	input := "this is a very long string"
	out := tr.Truncate(input)
	if len([]rune(out)) > 10 {
		t.Fatalf("expected truncated length <= 10, got %d: %q", len(out), out)
	}
	if !strings.HasSuffix(out, "...") {
		t.Fatalf("expected suffix '...', got %q", out)
	}
}

func TestTruncate_UnicodeSafe(t *testing.T) {
	tr := newTestTruncator(5)
	// each character is multi-byte
	input := "日本語テスト"
	out := tr.Truncate(input)
	runes := []rune(out)
	if len(runes) > 5 {
		t.Fatalf("expected <= 5 runes, got %d: %q", len(runes), out)
	}
}

func TestWasTruncated_True(t *testing.T) {
	tr := newTestTruncator(5)
	if !tr.WasTruncated("longer than five") {
		t.Fatal("expected WasTruncated to return true")
	}
}

func TestWasTruncated_False(t *testing.T) {
	tr := newTestTruncator(20)
	if tr.WasTruncated("short") {
		t.Fatal("expected WasTruncated to return false")
	}
}

func TestTruncateField_TrimsSpace(t *testing.T) {
	tr := newTestTruncator(20)
	out := tr.TruncateField("  padded  ")
	if out != "padded" {
		t.Fatalf("expected 'padded', got %q", out)
	}
}

func TestDefaultTruncateOptions(t *testing.T) {
	opts := DefaultTruncateOptions()
	if opts.MaxLength != 64 {
		t.Fatalf("expected MaxLength 64, got %d", opts.MaxLength)
	}
	if opts.Suffix != "..." {
		t.Fatalf("expected Suffix '...', got %q", opts.Suffix)
	}
}

func TestNewTruncator_ZeroMaxLength_UsesDefault(t *testing.T) {
	tr := NewTruncator(TruncateOptions{MaxLength: 0, Suffix: "..."})
	if tr.opts.MaxLength != 64 {
		t.Fatalf("expected default MaxLength 64, got %d", tr.opts.MaxLength)
	}
}
