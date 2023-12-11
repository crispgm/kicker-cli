package cmd

import (
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/util"
)

func init() {
	playerCmd.AddCommand(playerListCmd)
	rootCmd.AddCommand(playerCmd)
}

var playerCmd = &cobra.Command{
	Use:     "player",
	Aliases: []string{"players"},
	Short:   "Manage players in your organization",
	Long:    "Manage players in your organization",
	Run:     listPlayerCommand,
}

var playerListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List players",
	Long:    "List players",
	Run:     listPlayerCommand,
}

func listPlayerCommand(cmd *cobra.Command, args []string) {
	instance := initInstanceAndLoadConf()
	// load tournaments
	var table [][]string
	header := []string{"ID", "Name", "ITSF_ID", "ATSA_ID", "A.K.A"}
	if !globalNoHeaders {
		table = append(table, header)
	}
	needWrite := false
	for i, p := range instance.Conf.Players {
		if p.ID == "" {
			p.ID = util.UUID()
			instance.Conf.Players[i].ID = p.ID
			needWrite = true
		}
		table = append(table, []string{
			p.ID,
			p.Name,
			dashIfEmpty(p.ITSFID),
			dashIfEmpty(p.ATSAID),
			dashIfEmpty(strings.Join(p.Aliases, ", ")),
		})
	}
	_ = pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	if needWrite {
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
	}
}
