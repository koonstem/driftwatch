package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// BindSchemaFlags attaches schema-related flags to cmd.
func BindSchemaFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("schema", false, "emit a JSON schema summary of observed drift fields")
	cmd.Flags().Bool("schema-pretty", true, "pretty-print schema JSON output")
	cmd.Flags().String("schema-version", "v1", "schema document version tag")
}

// SchemaOptionsFromFlags builds SchemaOptions from bound cobra flags.
func SchemaOptionsFromFlags(cmd *cobra.Command) (SchemaOptions, error) {
	pretty, err := cmd.Flags().GetBool("schema-pretty")
	if err != nil {
		return SchemaOptions{}, fmt.Errorf("schema-flags: pretty: %w", err)
	}

	version, err := cmd.Flags().GetString("schema-version")
	if err != nil {
		return SchemaOptions{}, fmt.Errorf("schema-flags: version: %w", err)
	}

	if version == "" {
		return SchemaOptions{}, fmt.Errorf("schema-flags: --schema-version must not be empty")
	}

	return SchemaOptions{
		Pretty:  pretty,
		Version: version,
	}, nil
}

// IsSchemaEnabled returns true when the --schema flag is set.
func IsSchemaEnabled(cmd *cobra.Command) bool {
	v, err := cmd.Flags().GetBool("schema")
	if err != nil {
		return false
	}
	return v
}
