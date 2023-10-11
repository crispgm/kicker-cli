package cmd

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
)

func init() {
	playerCmd.AddCommand(playerCreateCmd)
}

var playerCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new"},
	Short:   "Create a player",
	Long:    "Create a player",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errorMessageAndExit("Please present at least one player name")
		}
		instance := initInstanceAndLoadConf()
		for _, arg := range args {
			for _, p := range instance.Conf.Players {
				if p.IsPlayer(arg) {
					errorMessageAndExit("Duplicated player name found:", arg)
				}
			}
		}
		np := entity.NewPlayer(args[0])
		if len(args) > 1 {
			np.AddAlias(args[1:]...)
		}
		instance.Conf.Players = append(instance.Conf.Players, *np)
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
		pterm.Println("1 player created")
	},
}
