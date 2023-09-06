// Package cmd .
package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
)

var globalNoColor bool

var rootCmd = &cobra.Command{
	Use:     "kicker-cli",
	Long:    "A Foosball data aggregator, analyzers, and manager based on Kickertool.",
	Version: app.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if globalNoColor {
			pterm.DisableColor()
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&globalNoColor, "no-colors", "", false, "show no colors")
}

// Execute .
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorMessageAndExit(err)
	}
}
