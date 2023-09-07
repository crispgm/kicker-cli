package cmd

import (
	"os"
	"strconv"

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

func convertToFloat(in string) float64 {
	out, _ := strconv.Atoi(in)
	return float64(out)
}
