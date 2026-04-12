package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/filter"
	"github.com/driftwatch/internal/output"
	"github.com/driftwatch/internal/runner"
	"github.com/driftwatch/internal/source"
	"github.com/spf13/cobra"
)

func runDetect(cmd *cobra.Command, args []string) error {
	cfgPath, _ := cmd.Flags().GetString("config")
	manifestOverride, _ := cmd.Flags().GetString("manifest")
	format, _ := cmd.Flags().GetString("output")
	failOnDrift, _ := cmd.Flags().GetBool("fail-on-drift")

	getFilterOpts := filter.BindFlags(cmd)

	cfg, err := config.Load(cfgPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	manifestPath := cfg.ManifestPath
	if manifestOverride != "" {
		manifestPath = manifestOverride
	}

	ldr := source.NewLoader()
	manifest, err := ldr.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	run, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	containers, err := run.ListContainers(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	detector := drift.NewDetector()
	results := detector.Detect(manifest, containers)

	filterOpts := getFilterOpts()
	results = filter.Filter(results, filterOpts)

	reporter := drift.NewReporter(results)
	report := reporter.Report()

	fmt_ := output.NewFormatter(os.Stdout, format)
	if err := fmt_.Write(report); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	exitCoder := output.NewExitCoder(failOnDrift)
	os.Exit(exitCoder.Code(report))
	return nil
}

var detectCmd = &cobra.Command{
	Use:   "detect",
	Short: "Detect configuration drift between running containers and manifest",
	RunE:  runDetect,
}

func init() {
	detectCmd.Flags().String("manifest", "", "Override manifest path from config")
	detectCmd.Flags().String("output", "text", "Output format: text, json, table")
	detectCmd.Flags().Bool("fail-on-drift", false, "Exit with non-zero code if drift is detected")
	rootCmd.AddCommand(detectCmd)
}
