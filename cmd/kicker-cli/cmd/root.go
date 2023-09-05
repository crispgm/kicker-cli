// Package cmd .
package cmd

import (
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
)

var rootCmd = &cobra.Command{
	Use:     "kicker-cli",
	Long:    "A Foosball data aggregator, analyzers, and manager based on Kickertool.",
	Version: app.Version,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {
	cobra.OnInitialize()
}

// Execute .
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorMessageAndExit(err)
	}
}
