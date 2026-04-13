package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

// HistoryOptions holds CLI flags for the history feature.
type HistoryOptions struct {
	Enabled  bool
	FilePath string
	MaxSize  int
}

// BindHistoryFlags registers history-related flags onto cmd.
func BindHistoryFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("history", false, "append results to a history log file")
	cmd.Flags().String("history-file", ".driftwatch-history.json", "path to the history log file")
	cmd.Flags().Int("history-max", 50, "maximum number of history entries to retain (0 = unlimited)")
}

// HistoryOptionsFromFlags reads history flag values from cmd.
func HistoryOptionsFromFlags(cmd *cobra.Command) (HistoryOptions, error) {
	enabled, err := cmd.Flags().GetBool("history")
	if err != nil {
		return HistoryOptions{}, fmt.Errorf("history flag: %w", err)
	}
	filePath, err := cmd.Flags().GetString("history-file")
	if err != nil {
		return HistoryOptions{}, fmt.Errorf("history-file flag: %w", err)
	}
	maxSize, err := cmd.Flags().GetInt("history-max")
	if err != nil {
		return HistoryOptions{}, fmt.Errorf("history-max flag: %w", err)
	}
	if maxSize < 0 {
		return HistoryOptions{}, fmt.Errorf("history-max must be >= 0, got %d", maxSize)
	}
	return HistoryOptions{
		Enabled:  enabled,
		FilePath: filePath,
		MaxSize:  maxSize,
	}, nil
}
