package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(evaluateCmd)
}

var evaluateCmd = &cobra.Command{
	Use:     "evaluate",
	Aliases: []string{"eval"},
	Short:   "Simple tool to evaluate estimated changes between two teams/players",
	Long: `Simple tool to evaluate estimated changes between two teams/players.
$ kicker-cli evaluate -a elo 1100 1200
$ kicker-cli evaluate -a elo 1103 1203 1289 1013
$ kicker-cli evaluate -a elo -k 20 1103 1203 1289 1013
$ kicker-cli evaluate -a rank 1
$ kicker-cli evaluate -a rank -s KLocal -s ATSA50 3`,
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Usage()
	},
}
