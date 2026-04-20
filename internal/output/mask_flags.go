package output

import (
	"strings"

	"github.com/spf13/cobra"
)

// BindMaskFlags registers mask-related flags on the given command.
func BindMaskFlags(cmd *cobra.Command) {
	defaults := DefaultMaskOptions()
	cmd.Flags().Bool("mask", defaults.Enabled, "enable field masking in output")
	cmd.Flags().String("mask-char", defaults.MaskChar, "character used for masking")
	cmd.Flags().Int("mask-length", defaults.MaskLength, "number of mask characters to use")
	cmd.Flags().StringSlice("mask-fields", defaults.Fields, "comma-separated list of fields to mask")
}

// MaskOptionsFromFlags reads mask options from the command's flags.
func MaskOptionsFromFlags(cmd *cobra.Command) MaskOptions {
	defaults := DefaultMaskOptions()
	enabled, err := cmd.Flags().GetBool("mask")
	if err != nil {
		enabled = defaults.Enabled
	}
	maskChar, err := cmd.Flags().GetString("mask-char")
	if err != nil || strings.TrimSpace(maskChar) == "" {
		maskChar = defaults.MaskChar
	}
	maskLen, err := cmd.Flags().GetInt("mask-length")
	if err != nil || maskLen <= 0 {
		maskLen = defaults.MaskLength
	}
	fields, err := cmd.Flags().GetStringSlice("mask-fields")
	if err != nil {
		fields = defaults.Fields
	}
	return MaskOptions{
		Enabled:    enabled,
		MaskChar:   maskChar,
		MaskLength: maskLen,
		Fields:     fields,
	}
}
