package cmd

import (
	"github.com/skratchdot/open-golang/open"
	"github.com/spf13/cobra"
)

func init() {
	eventCmd.AddCommand(eventOpenCmd)
}

var eventOpenCmd = &cobra.Command{
	Use:   "open",
	Short: "Open event URL",
	Long:  "Open event URL",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			errorMessageAndExit("Please present an event ID or name")
		}
		instance := initInstanceAndLoadConf()
		if e := instance.GetEvent(args[0]); e != nil {
			if len(e.URL) == 0 {
				errorMessageAndExit("URL is not set for this event")
			}
			err := open.Run(e.URL)
			if err != nil {
				errorMessageAndExit(err)
			}
		} else {
			errorMessageAndExit("No event(s) found")
		}
	},
}
