package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/util"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	importEventName                 string
	importEventPoints               int
	importEventCreateUnknownPlayers bool
)

func init() {
	importCmd.Flags().StringVarP(&importEventName, "name", "n", "", "event name")
	importCmd.Flags().IntVarP(&importEventPoints, "points", "", entity.DefaultPoints, "points for the event")
	importCmd.Flags().BoolVarP(&importEventCreateUnknownPlayers, "create-unknown-players", "c", false, "create unknown players during importing")
	rootCmd.AddCommand(importCmd)
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import one .ktool file to create an event",
	Long:  "Import one .ktool file and it will be converted to an event automatically",
	Run: func(cmd *cobra.Command, args []string) {
		instance := initInstanceAndLoadConf()
		if len(args) == 0 {
			return
		}

		for _, importPath := range args {
			if strings.HasPrefix(importPath, "~") {
				importPath = util.ExpandUserHomeDir(importPath)
			}
			t, err := parser.ParseFile(importPath)
			if err != nil {
				errorMessageAndExit("Unable to parse file:", err)
			}
			if importEventName == "" {
				importEventName = t.Name
			}

			for _, p := range t.Players {
				found := false
				for _, cp := range instance.Conf.Players {
					if cp.IsPlayer(p.Name) {
						found = true
					}
				}
				if !found {
					var needCreate bool
					if !importEventCreateUnknownPlayers {
						needCreate, _ = pterm.DefaultInteractiveConfirm.
							WithDefaultText(fmt.Sprintf("Create a new player with name `%s`", p.Name)).
							WithDefaultValue(true).
							Show()
					} else {
						needCreate = true
					}
					if needCreate {
						instance.AddPlayer(*entity.NewPlayer(p.Name))
					}
				}
			}

			pterm.Printfln("Importing \"%s\" `%s` ...", importEventName, importPath)
			event := *entity.NewEvent("temp", importEventName, importEventPoints)
			fn := fmt.Sprintf("%s.ktool", event.ID)
			event.Path = fn
			err = util.CopyFile(importPath, filepath.Join(instance.DataPath(), fn))
			if err != nil {
				errorMessageAndExit(err)
			}
			instance.AddEvent(event)
		}
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
		pterm.Printfln("%d event(s) imported", len(instance.Conf.Events))
	},
}
