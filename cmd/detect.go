package cmd

import (
	"fmt"
	"os"

	"github.com/driftwatch/internal/config"
	"github.com/driftwatch/internal/drift"
	"github.com/driftwatch/internal/output"
	"github.com/driftwatch/internal/runner"
	"github.com/driftwatch/internal/source"
	"github.com/spf13/cobra"
)

var (
	configFile   string
	manifestFile string
	formatFlag   string
	failOnDrift  bool
)

func runDetect(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}

	manifestPath := manifestFile
	if manifestPath == "" {
		manifestPath = cfg.ManifestPath
	}

	loader := source.NewLoader()
	manifest, err := loader.Load(manifestPath)
	if err != nil {
		return fmt.Errorf("loading manifest: %w", err)
	}

	r, err := runner.New(cfg)
	if err != nil {
		return fmt.Errorf("creating runner: %w", err)
	}

	containers, err := r.ListContainers(cmd.Context())
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}

	detector := drift.NewDetector()
	results := detector.Detect(manifest, containers)

	report := drift.NewReporter().Build(results)

	fmt := output.Format(formatFlag)
	formatter := output.NewFormatter(fmt, os.Stdout)
	if err := formatter.Write(report); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}

	exitCoder := output.NewExitCoder(failOnDrift)
	os.Exit(exitCoder.Code(report))
	return nil
}
