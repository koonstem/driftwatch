package output

import (
	"strings"
	"testing"
)

func TestColorizer_Disabled(t *testing.T) {
	c := NewColorizer(false)

	cases := []struct {
		name string
		fn   func(string) string
	}{
		{"Red", c.Red},
		{"Green", c.Green},
		{"Yellow", c.Yellow},
		{"Cyan", c.Cyan},
		{"Bold", c.Bold},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn("hello")
			if got != "hello" {
				t.Errorf("expected plain text, got %q", got)
			}
		})
	}
}

func TestColorizer_Enabled_ContainsEscape(t *testing.T) {
	c := NewColorizer(true)

	cases := []struct {
		name string
		fn   func(string) string
	}{
		{"Red", c.Red},
		{"Green", c.Green},
		{"Yellow", c.Yellow},
		{"Cyan", c.Cyan},
		{"Bold", c.Bold},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn("hello")
			if !strings.Contains(got, "\033[") {
				t.Errorf("expected ANSI escape in output, got %q", got)
			}
			if !strings.Contains(got, "hello") {
				t.Errorf("expected original text preserved, got %q", got)
			}
			if !strings.HasSuffix(got, colorReset) {
				t.Errorf("expected reset suffix, got %q", got)
			}
		})
	}
}

func TestColorizer_Enabled_DifferentCodes(t *testing.T) {
	c := NewColorizer(true)

	red := c.Red("x")
	green := c.Green("x")

	if red == green {
		t.Error("expected Red and Green to produce different output")
	}
}
