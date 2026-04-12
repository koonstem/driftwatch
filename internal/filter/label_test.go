package filter

import (
	"testing"
)

func TestParseLabelSelector_Valid(t *testing.T) {
	sel, err := ParseLabelSelector("env=production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Key != "env" {
		t.Errorf("expected key %q, got %q", "env", sel.Key)
	}
	if sel.Value != "production" {
		t.Errorf("expected value %q, got %q", "production", sel.Value)
	}
}

func TestParseLabelSelector_ValueWithEquals(t *testing.T) {
	sel, err := ParseLabelSelector("tag=v1=2")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if sel.Key != "tag" {
		t.Errorf("expected key %q, got %q", "tag", sel.Key)
	}
	if sel.Value != "v1=2" {
		t.Errorf("expected value %q, got %q", "v1=2", sel.Value)
	}
}

func TestParseLabelSelector_Empty(t *testing.T) {
	_, err := ParseLabelSelector("")
	if err == nil {
		t.Fatal("expected error for empty selector")
	}
}

func TestParseLabelSelector_NoEquals(t *testing.T) {
	_, err := ParseLabelSelector("justkey")
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestParseLabelSelector_EmptyKey(t *testing.T) {
	_, err := ParseLabelSelector("=value")
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestLabelSelector_String(t *testing.T) {
	sel := LabelSelector{Key: "app", Value: "nginx"}
	if sel.String() != "app=nginx" {
		t.Errorf("expected %q, got %q", "app=nginx", sel.String())
	}
}

func TestLabelSelector_Matches_True(t *testing.T) {
	sel := LabelSelector{Key: "env", Value: "staging"}
	labels := map[string]string{"env": "staging", "team": "platform"}
	if !sel.Matches(labels) {
		t.Error("expected selector to match labels")
	}
}

func TestLabelSelector_Matches_WrongValue(t *testing.T) {
	sel := LabelSelector{Key: "env", Value: "staging"}
	labels := map[string]string{"env": "production"}
	if sel.Matches(labels) {
		t.Error("expected selector not to match labels with different value")
	}
}

func TestLabelSelector_Matches_MissingKey(t *testing.T) {
	sel := LabelSelector{Key: "env", Value: "staging"}
	labels := map[string]string{"team": "platform"}
	if sel.Matches(labels) {
		t.Error("expected selector not to match labels missing the key")
	}
}

func TestParseLabelSelectors_Multiple(t *testing.T) {
	selectors, err := ParseLabelSelectors([]string{"env=prod", "team=infra"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(selectors) != 2 {
		t.Fatalf("expected 2 selectors, got %d", len(selectors))
	}
}

func TestParseLabelSelectors_Invalid(t *testing.T) {
	_, err := ParseLabelSelectors([]string{"env=prod", "badformat"})
	if err == nil {
		t.Fatal("expected error for invalid selector in list")
	}
}
