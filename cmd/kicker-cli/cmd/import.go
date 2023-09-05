package cmd

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/util"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	importPath        string
	importEventName   string
	importEventPoints int
)

func init() {
	importCmd.Flags().StringVarP(&importPath, "path", "p", "", "path to imported files")
	importCmd.Flags().StringVarP(&importEventName, "name", "n", "", "event name")
	importCmd.Flags().IntVarP(&importEventPoints, "points", "", entity.DefaultPoints, "points for the event")
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import one .ktool file to create an event",
	Long:  "Import one .ktool file and it will be converted to an event automatically",
	Run: func(cmd *cobra.Command, args []string) {
		instance := initInstanceAndLoadConf()

		t, err := parser.ParseFile(importPath)
		if err != nil {
			errorMessageAndExit("Unable to parse file")
		}
		if importEventName == "" {
			importEventName = t.Name
		}

		pterm.Info.Printfln("Importing \"%s\" `%s` ...", importEventName, importPath)
		fn := filepath.Base(importPath)
		event := *entity.NewEvent(fn, importEventName, importEventPoints)
		err = util.CopyFile(importPath, filepath.Join(instance.DataPath(), fn))
		instance.Conf.Events = append(instance.Conf.Events, event)
		err = instance.WriteConf()
		if err != nil {
			pterm.Error.Println(err)
		}
	},
}
