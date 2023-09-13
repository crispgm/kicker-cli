package cmd

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

func init() {
	eventCmd.AddCommand(eventDeleteCmd)
}

var eventDeleteCmd = &cobra.Command{
	Use:     "delete",
	Aliases: []string{"remove", "rm"},
	Short:   "Delete an event",
	Long:    "Delete an event",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errorMessageAndExit("Please present an event ID")
		}
		eventsDeleted := 0
		instance := initInstanceAndLoadConf()
		for _, arg := range args {
			if e := instance.GetEvent(arg); e != nil {
				path := filepath.Join(instance.DataPath(), e.Path)
				err := instance.DeleteEvent(arg)
				if err != nil {
					errorMessageAndExit(err)
				}
				err = os.Remove(path)
				if err != nil {
					pterm.Error.Println(err)
				}
				eventsDeleted++
			} else {
				errorMessageAndExit("Event not found:", arg)
			}
		}
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
		pterm.Printfln("%d event deleted", eventsDeleted)
	},
}
