package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(orgCmd)
}

var orgCmd = &cobra.Command{
	Use:     "organization",
	Aliases: []string{"org"},
	Short:   "Manage your organization",
	Long:    "Manage your organization",
	Run: func(cmd *cobra.Command, args []string) {
		instance := initInstanceAndLoadConf()
		var table [][]string
		header := []string{"ID", "Name", "Players", "Events", "Kickertool ID"}
		if !globalNoHeaders {
			table = append(table, header)
		}
		table = append(table, []string{
			instance.Conf.Organization.ID,
			instance.Conf.Organization.Name,
			fmt.Sprintf("%d", len(instance.Conf.Players)),
			fmt.Sprintf("%d", len(instance.Conf.Events)),
			dashIfEmpty(instance.Conf.Organization.KickerToolID),
		})
		if len(table) <= 1 {
			errorMessageAndExit("No organization found")
		}
		_ = pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	},
}
