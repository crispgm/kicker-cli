package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
)

func init() {
	eventCmd.AddCommand(eventInfoCmd)
}

var eventInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"show"},
	Short:   "Show event details",
	Long:    "Show event details",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errorMessageAndExit("Please present an event ID")
		}
		arg := args[0]
		instance := initInstanceAndLoadConf()
		e := instance.GetEvent(arg)
		if e == nil {
			errorMessageAndExit("No event(s) found")
		}
		table := initEventInfoHeader()
		_, r, _ := loadAndShowEventInfo(&table, instance.DataPath(), instance.Conf.Players, e)
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
		table = showGames(r.PreliminaryRounds)
		if len(table) > 0 {
			pterm.Println("Rounds:")
			pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
		table = showGames(r.LoserBracket)
		if len(table) > 0 {
			pterm.Println("Loser Bracket:")
			pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
		table = showGames(r.WinnerBracket)
		if len(table) > 0 {
			pterm.Println("Winner Bracket:")
			pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
	},
}

func showGames(games []entity.Game) [][]string {
	var roundTable [][]string
	for _, g := range games {
		if len(g.Team1) == 1 {
			roundTable = append(roundTable, []string{
				fmt.Sprintf("%s", g.Team1[0]),
				fmt.Sprintf("%d:%d", g.Point1, g.Point2),
				fmt.Sprintf("%s", g.Team2[0]),
			})
		} else {
			roundTable = append(roundTable, []string{
				fmt.Sprintf("%s/%s", g.Team1[0], g.Team1[1]),
				fmt.Sprintf("%d:%d", g.Point1, g.Point2),
				fmt.Sprintf("%s/%s", g.Team2[0], g.Team2[1]),
			})
		}
	}
	return roundTable
}
