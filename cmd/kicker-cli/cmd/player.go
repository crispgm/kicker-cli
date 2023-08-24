package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(eloCmd)
}

var playerCmd = &cobra.Command{
	Use:   "player",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		pterm.Info.Println("")
	},
}
