package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourorg/driftwatch/internal/config"
	"github.com/yourorg/driftwatch/internal/drift"
	"github.com/yourorg/driftwatch/internal/output"
	"github.com/yourorg/driftwatch/internal/runner"
	"github.com/yourorg/driftwatch/internal/source"
)

func runDetect(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	manifestPath := cfg.ManifestPath
	if manifest != "" {
		manifestPath = manifest
	}

	loader := source.NewLoader()
	svcManifest, err := loader.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	containers, err := r.ListContainers()
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	detector := drift.NewDetector()
	results := detector.Detect(svcManifest, containers)

	reporter := drift.NewReporter(results)
	report := reporter.Build()

	fmt := output.NewFormatter(outputFmt, os.Stdout)
	if err := fmt.Write(report); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	exitCoder := output.NewExitCoder(failOnDrift)
	os.Exit(exitCoder.Code(report))
	return nil
}
