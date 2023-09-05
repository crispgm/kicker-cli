package cmd

import (
	"os"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/app"
)

func initInstanceAndLoadConf() *app.App {
	instance := app.NewApp(initPath, app.DefaultName)
	err := instance.LoadConf()
	if err != nil {
		errorMessageAndExit("Not a valid kicker workspace")
	}

	return instance
}

func errorMessageAndExit(a ...interface{}) {
	pterm.Error.Println(a...)
	os.Exit(1)
}
