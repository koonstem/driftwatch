package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile    string
	manifest   string
	outputFmt  string
	failOnDrift bool
)

var rootCmd = &cobra.Command{
	Use:   "driftwatch",
	Short: "Detect configuration drift between running services and IaC definitions",
	Long: `driftwatch compares running container state against declared service
definitions and reports any configuration drift detected.`,
	RunE: runDetect,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "driftwatch.yaml", "path to config file")
	rootCmd.PersistentFlags().StringVarP(&manifest, "manifest", "m", "", "path to service manifest file (overrides config)")
	rootCmd.PersistentFlags().StringVarP(&outputFmt, "output", "o", "text", "output format: text or json")
	rootCmd.PersistentFlags().BoolVar(&failOnDrift, "fail-on-drift", false, "exit with non-zero code when drift is detected")
}
