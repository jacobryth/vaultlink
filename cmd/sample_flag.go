package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"vaultlink/internal/sample"
)

func registerSampleFlags(cmd *cobra.Command) {
	cmd.Flags().String("sample", "none", `sampling strategy: none | random | nth`)
	cmd.Flags().Int("sample-n", 10, "number of secrets to sample (random) or step size (nth)")
}

func resolvedSampleLevel(cmd *cobra.Command) (sample.Level, error) {
	raw, err := cmd.Flags().GetString("sample")
	if err != nil {
		return sample.LevelNone, fmt.Errorf("sample flag: %w", err)
	}
	lvl := sample.Level(raw)
	switch lvl {
	case sample.LevelNone, sample.LevelRandom, sample.LevelNth:
		return lvl, nil
	default:
		return sample.LevelNone, fmt.Errorf("sample: unknown level %q", raw)
	}
}

func resolvedSampleN(cmd *cobra.Command) (int, error) {
	n, err := cmd.Flags().GetInt("sample-n")
	if err != nil {
		return 0, fmt.Errorf("sample-n flag: %w", err)
	}
	if n < 1 {
		return 0, fmt.Errorf("sample-n must be >= 1, got %d", n)
	}
	return n, nil
}
