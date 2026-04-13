package output

import (
	"regexp"
	"strings"
)

// RedactOptions controls which field values are masked in output.
type RedactOptions struct {
	Enabled  bool
	Patterns []string
	Mask     string
}

// DefaultRedactOptions returns safe defaults: redaction disabled.
func DefaultRedactOptions() RedactOptions {
	return RedactOptions{
		Enabled:  false,
		Patterns: []string{},
		Mask:     "[REDACTED]",
	}
}

// Redactor masks sensitive values in drift field strings.
type Redactor struct {
	opts     RedactOptions
	compiled []*regexp.Regexp
}

// NewRedactor compiles the provided patterns and returns a Redactor.
func NewRedactor(opts RedactOptions) (*Redactor, error) {
	var compiled []*regexp.Regexp
	for _, p := range opts.Patterns {
		re, err := regexp.Compile(p)
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return &Redactor{opts: opts, compiled: compiled}, nil
}

// Redact returns the value with any matching patterns replaced by the mask.
// If redaction is disabled the original value is returned unchanged.
func (r *Redactor) Redact(value string) string {
	if !r.opts.Enabled {
		return value
	}
	for _, re := range r.compiled {
		value = re.ReplaceAllString(value, r.opts.Mask)
	}
	return value
}

// RedactFields applies Redact to every field value in the provided map,
// returning a new map with masked values.
func (r *Redactor) RedactFields(fields map[string]string) map[string]string {
	out := make(map[string]string, len(fields))
	for k, v := range fields {
		out[k] = r.Redact(v)
	}
	return out
}

// ContainsSensitiveKey is a heuristic helper that returns true when a field
// key looks like it might hold a secret (password, token, key, secret).
func ContainsSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	for _, hint := range []string{"password", "token", "secret", "key", "apikey", "api_key"} {
		if strings.Contains(lower, hint) {
			return true
		}
	}
	return false
}
