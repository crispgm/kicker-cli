// Package cmd .
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
)

var (
	dryRun bool
)

var rootCmd = &cobra.Command{
	Use:     "kicker-cli",
	Long:    `A Foosball data aggregator, analyzers, and manager based on Kickertool.`,
	Version: app.Version,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "", false, "Dry Run")
}

// Execute .
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
