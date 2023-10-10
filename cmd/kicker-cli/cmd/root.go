// Package cmd .
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
)

var (
	globalNoColors  bool
	globalNoHeaders bool
	globalNoBoxes   bool
)

var rootCmd = &cobra.Command{
	Use:     "kicker-cli",
	Long:    "A Foosball data aggregator, analyzers, and manager based on Kickertool.",
	Version: app.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if globalNoColors {
			pterm.DisableColor()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&globalNoColors, "no-colors", "", false, "show no colors")
	rootCmd.PersistentFlags().BoolVarP(&globalNoHeaders, "no-headers", "", false, "no table headers")
	rootCmd.PersistentFlags().BoolVarP(&globalNoBoxes, "no-boxes", "", false, "no table boxes")
}

// Execute .
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorMessageAndExit(err)
	}
}
