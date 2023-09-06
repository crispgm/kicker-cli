package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	eventIDOrName string
	eventGameMode string
	eventNameType string
	allEvents     bool
)

func init() {
	eventCmd.PersistentFlags().BoolVarP(&allEvents, "all", "a", false, "rank all events")
	eventCmd.PersistentFlags().StringVarP(&eventIDOrName, "name", "n", "", "event ID or name")
	eventCmd.PersistentFlags().StringVarP(&eventGameMode, "mode", "m", "", "rank mode")
	eventCmd.PersistentFlags().StringVarP(&eventNameType, "name-type", "t", "", "name type (single, byp, dyp or monster_dyp)")
	eventCmd.AddCommand(eventListCmd)
	rootCmd.AddCommand(eventCmd)
}

var eventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"events"},
	Short:   "Manage events organized by your organization",
	Long:    "Manage events organized by your organization",
	Run:     eventListCommand,
}

var eventListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List events",
	Long:    "List events",
	Run:     eventListCommand,
}

func eventListCommand(cmd *cobra.Command, args []string) {
	instance := initInstanceAndLoadConf()
	// load tournaments
	var table [][]string
	header := []string{"ID", "Name", "Points", "URL"}
	table = append(table, header)
	if len(eventIDOrName) > 0 {
		if e := instance.GetEvent(eventIDOrName); e != nil {
			table = append(table, []string{
				e.ID,
				e.Name,
				fmt.Sprintf("%d", e.Points),
				"-",
			})
		}
	} else {
		for _, e := range instance.Conf.Events {
			url := e.URL
			if url == "" {
				url = "-"
			}
			table = append(table, []string{
				e.ID,
				e.Name,
				fmt.Sprintf("%d", e.Points),
				url,
			})
		}
	}
	if len(table) == 1 {
		errorMessageAndExit("No event(s) found.")
	}
	pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
}
