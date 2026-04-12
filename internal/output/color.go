package output

import "fmt"

// ANSI color codes
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
	colorBold   = "\033[1m"
)

// Colorizer controls whether ANSI color codes are emitted.
type Colorizer struct {
	enabled bool
}

// NewColorizer returns a Colorizer. When enabled is false all methods
// return plain text without escape sequences.
func NewColorizer(enabled bool) *Colorizer {
	return &Colorizer{enabled: enabled}
}

// Red wraps s in red color codes.
func (c *Colorizer) Red(s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", colorRed, s, colorReset)
}

// Green wraps s in green color codes.
func (c *Colorizer) Green(s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", colorGreen, s, colorReset)
}

// Yellow wraps s in yellow color codes.
func (c *Colorizer) Yellow(s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", colorYellow, s, colorReset)
}

// Cyan wraps s in cyan color codes.
func (c *Colorizer) Cyan(s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", colorCyan, s, colorReset)
}

// Bold wraps s in bold escape codes.
func (c *Colorizer) Bold(s string) string {
	if !c.enabled {
		return s
	}
	return fmt.Sprintf("%s%s%s", colorBold, s, colorReset)
}
