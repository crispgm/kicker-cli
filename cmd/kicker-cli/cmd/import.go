package cmd

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/util"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	importPath string
	eventName  string
	points     int
)

func init() {
	importCmd.Flags().StringVarP(&importPath, "path", "p", "", "Path to imported files")
	importCmd.Flags().StringVarP(&eventName, "name", "n", "", "Event name")
	importCmd.Flags().IntVarP(&points, "points", "", entity.DefaultPoints, "Points for the event")
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import one .ktool file to create an event",
	Long:  "Import one .ktool file and it will be converted to an event automatically",
	Run: func(cmd *cobra.Command, args []string) {
		if len(initPath) == 0 || initPath == "." {
			pterm.Info.Println("Initializing current folder...")
		}
		// check `.kicker.yaml`
		instance := app.NewApp(initPath, app.DefaultName)
		err := instance.LoadConf()
		if err != nil {
			pterm.Error.Println("Not a valid kicker workspace")
			os.Exit(1)
		}

		t, err := parser.ParseFile(importPath)
		if err != nil {
			pterm.Error.Println("Unable to parse file")
			os.Exit(1)
		}
		if eventName == "" {
			eventName = t.Name
		}

		pterm.Info.Printfln("Importing \"%s\" `%s` ...", eventName, importPath)
		fn := filepath.Base(importPath)
		event := *entity.NewEvent(fn, eventName, points)
		err = util.CopyFile(importPath, filepath.Join(instance.DataPath(), fn))
		instance.Conf.Events = append(instance.Conf.Events, event)
		err = instance.WriteConf()
		if err != nil {
			pterm.Error.Println(err)
		}
	},
}
