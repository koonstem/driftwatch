package filter

import (
	"fmt"
	"strings"
)

// LabelSelector represents a parsed key=value label selector.
type LabelSelector struct {
	Key   string
	Value string
}

// ParseLabelSelector parses a label selector string of the form "key=value".
// Returns an error if the format is invalid.
func ParseLabelSelector(raw string) (LabelSelector, error) {
	if raw == "" {
		return LabelSelector{}, fmt.Errorf("label selector must not be empty")
	}

	parts := strings.SplitN(raw, "=", 2)
	if len(parts) != 2 {
		return LabelSelector{}, fmt.Errorf("invalid label selector %q: expected key=value format", raw)
	}

	key := strings.TrimSpace(parts[0])
	value := strings.TrimSpace(parts[1])

	if key == "" {
		return LabelSelector{}, fmt.Errorf("label selector key must not be empty in %q", raw)
	}

	return LabelSelector{Key: key, Value: value}, nil
}

// ParseLabelSelectors parses multiple label selector strings.
func ParseLabelSelectors(raws []string) ([]LabelSelector, error) {
	selectors := make([]LabelSelector, 0, len(raws))
	for _, raw := range raws {
		sel, err := ParseLabelSelector(raw)
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, sel)
	}
	return selectors, nil
}

// String returns the canonical string representation of a LabelSelector.
func (l LabelSelector) String() string {
	return l.Key + "=" + l.Value
}

// Matches reports whether the given labels map satisfies this selector.
func (l LabelSelector) Matches(labels map[string]string) bool {
	v, ok := labels[l.Key]
	if !ok {
		return false
	}
	return v == l.Value
}
