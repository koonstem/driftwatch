package output

import "strings"

// TruncateOptions controls how long string values are truncated in output.
type TruncateOptions struct {
	MaxLength int
	Suffix    string
}

// DefaultTruncateOptions returns sensible defaults for truncation.
func DefaultTruncateOptions() TruncateOptions {
	return TruncateOptions{
		MaxLength: 64,
		Suffix:    "...",
	}
}

// NewTruncator creates a Truncator with the given options.
func NewTruncator(opts TruncateOptions) *Truncator {
	if opts.MaxLength <= 0 {
		opts.MaxLength = DefaultTruncateOptions().MaxLength
	}
	if opts.Suffix == "" {
		opts.Suffix = DefaultTruncateOptions().Suffix
	}
	return &Truncator{opts: opts}
}

// Truncator shortens strings that exceed a configured maximum length.
type Truncator struct {
	opts TruncateOptions
}

// Truncate returns s shortened to MaxLength runes, appending Suffix if cut.
func (t *Truncator) Truncate(s string) string {
	runes := []rune(s)
	if len(runes) <= t.opts.MaxLength {
		return s
	}
	cutAt := t.opts.MaxLength - len([]rune(t.opts.Suffix))
	if cutAt < 0 {
		cutAt = 0
	}
	return string(runes[:cutAt]) + t.opts.Suffix
}

// TruncateField truncates a named field value and returns a display-ready string.
// If the value was not shortened the original is returned unchanged.
func (t *Truncator) TruncateField(value string) string {
	return t.Truncate(strings.TrimSpace(value))
}

// WasTruncated reports whether Truncate would shorten s.
func (t *Truncator) WasTruncated(s string) bool {
	return len([]rune(s)) > t.opts.MaxLength
}
