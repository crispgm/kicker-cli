package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	playerCmd.AddCommand(playerDeleteCmd)
}

var playerDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete a player",
	Long:    "Delete a player",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errorMessageAndExit("Please present a player ID")
		}
		playersDeleted := 0
		instance := initInstanceAndLoadConf()
		for _, arg := range args {
			if p := instance.GetPlayer(arg); p != nil {
				err := instance.DeletePlayer(arg)
				if err != nil {
					errorMessageAndExit(err)
				}
				playersDeleted++
			} else {
				errorMessageAndExit("Player not found:", arg)
			}
		}
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
		pterm.Printfln("%d player deleted", playersDeleted)
	},
}
