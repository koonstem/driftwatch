package output

import "github.com/spf13/cobra"

// BindTransformFlags registers transform-related flags on the given command.
func BindTransformFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("transform-uppercase-service", false, "uppercase service names in output")
	cmd.Flags().Bool("transform-trim-image-tag", false, "strip image tags from image field values")
	cmd.Flags().StringToString("transform-rename-fields", map[string]string{}, "rename output fields (e.g. image=container_image)")
}

// TransformOptionsFromFlags builds TransformOptions from bound cobra flags.
func TransformOptionsFromFlags(cmd *cobra.Command) (TransformOptions, error) {
	upper, err := cmd.Flags().GetBool("transform-uppercase-service")
	if err != nil {
		return TransformOptions{}, err
	}
	trim, err := cmd.Flags().GetBool("transform-trim-image-tag")
	if err != nil {
		return TransformOptions{}, err
	}
	rename, err := cmd.Flags().GetStringToString("transform-rename-fields")
	if err != nil {
		return TransformOptions{}, err
	}
	return TransformOptions{
		UppercaseService: upper,
		TrimImageTag:     trim,
		RenameFields:     rename,
	}, nil
}
