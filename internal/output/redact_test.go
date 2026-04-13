package output

import (
	"testing"
)

func TestRedactor_Disabled_PassesThrough(t *testing.T) {
	opts := DefaultRedactOptions()
	r, err := NewRedactor(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Redact("super-secret-value")
	if got != "super-secret-value" {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestRedactor_Enabled_MasksMatch(t *testing.T) {
	opts := RedactOptions{
		Enabled:  true,
		Patterns: []string{`password=\S+`},
		Mask:     "[REDACTED]",
	}
	r, err := NewRedactor(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Redact("password=hunter2")
	if got != "[REDACTED]" {
		t.Errorf("expected [REDACTED], got %q", got)
	}
}

func TestRedactor_Enabled_NoMatch_Unchanged(t *testing.T) {
	opts := RedactOptions{
		Enabled:  true,
		Patterns: []string{`token=\S+`},
		Mask:     "[REDACTED]",
	}
	r, err := NewRedactor(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := r.Redact("image=nginx:latest")
	if got != "image=nginx:latest" {
		t.Errorf("expected value unchanged, got %q", got)
	}
}

func TestRedactor_InvalidPattern_ReturnsError(t *testing.T) {
	opts := RedactOptions{
		Enabled:  true,
		Patterns: []string{`[invalid`},
		Mask:     "[REDACTED]",
	}
	_, err := NewRedactor(opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern")
	}
}

func TestRedactFields_MasksMatchingValues(t *testing.T) {
	opts := RedactOptions{
		Enabled:  true,
		Patterns: []string{`secret`},
		Mask:     "***",
	}
	r, err := NewRedactor(opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fields := map[string]string{
		"image":    "nginx:latest",
		"env_pass": "mysecretvalue",
	}
	out := r.RedactFields(fields)
	if out["image"] != "nginx:latest" {
		t.Errorf("image should be unchanged, got %q", out["image"])
	}
	if out["env_pass"] != "my***value" {
		t.Errorf("env_pass should be masked, got %q", out["env_pass"])
	}
}

func TestContainsSensitiveKey(t *testing.T) {
	cases := []struct {
		key      string
		expected bool
	}{
		{"DB_PASSWORD", true},
		{"api_key", true},
		{"auth_token", true},
		{"image", false},
		{"replicas", false},
		{"SECRET_VALUE", true},
	}
	for _, tc := range cases {
		got := ContainsSensitiveKey(tc.key)
		if got != tc.expected {
			t.Errorf("ContainsSensitiveKey(%q) = %v, want %v", tc.key, got, tc.expected)
		}
	}
}
