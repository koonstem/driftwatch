package filter

import (
	"strings"

	"github.com/spf13/cobra"
)

// BindFlags attaches filter-related flags to a cobra command and returns
// a function that builds an Options from the parsed flag values.
func BindFlags(cmd *cobra.Command) func() Options {
	var services string
	var onlyDrifted bool
	var labelSelector string

	cmd.Flags().StringVarP(&services, "services", "s", "", "Comma-separated list of service names to include")
	cmd.Flags().BoolVar(&onlyDrifted, "only-drifted", false, "Show only services with detected drift")
	cmd.Flags().StringVarP(&labelSelector, "label", "l", "", "Filter by field selector (e.g. image=nginx)")

	return func() Options {
		var svcList []string
		if services != "" {
			for _, s := range strings.Split(services, ",") {
				s = strings.TrimSpace(s)
				if s != "" {
					svcList = append(svcList, s)
				}
			}
		}
		return Options{
			Services:      svcList,
			OnlyDrifted:   onlyDrifted,
			LabelSelector: labelSelector,
		}
	}
}
