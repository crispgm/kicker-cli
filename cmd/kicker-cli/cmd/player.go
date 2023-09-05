package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(playerCmd)
}

var playerCmd = &cobra.Command{
	Use:     "player",
	Aliases: []string{"players"},
	Short:   "Manage players in your organization",
	Long:    "Manage players in your organization",
	Run: func(cmd *cobra.Command, args []string) {
		instance := initInstanceAndLoadConf()
		// load tournaments
		var table [][]string
		header := []string{"ID", "Name", "Points", "Played", "Won", "Lost", "W%", "ELO"}
		table = append(table, header)
		for _, p := range instance.Conf.Players {
			table = append(table, []string{
				p.ID,
				p.Name,
				"",
				fmt.Sprintf("%d", p.Played),
				fmt.Sprintf("%d", p.Won),
				fmt.Sprintf("%d", p.Lost),
				fmt.Sprintf("%.0f%%", p.WinRate),
				fmt.Sprintf("%.0f", p.EloRating),
			})
		}
		pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
	},
}
