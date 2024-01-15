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
	"github.com/crispgm/kicker-cli/pkg/rating"
)

var (
	importEventName                 string
	importEventLevel                string
	importEventCreateUnknownPlayers bool
)

func init() {
	importCmd.Flags().StringVarP(&importEventName, "name", "n", "", "event name")
	importCmd.Flags().StringVarP(&importEventLevel, "points", "", rating.KLocal, "points for the event")
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

		var players []string
		for _, p := range instance.Conf.Players {
			players = append(players, p.Name)
		}
		eventsAdded := 0
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
					if !importEventCreateUnknownPlayers {
						text := fmt.Sprint("Unknown player: ", p.Name)
						selectedOption, _ := pterm.
							DefaultInteractiveSelect.
							WithOptions([]string{
								"Create a new player",
								"Attach to an existed player",
								"Cancel",
							}).
							Show(text)

						if selectedOption == "Create a new player" {
							instance.AddPlayer(*entity.NewPlayer(p.Name))
							pterm.Printfln("Creating player %s", p.Name)
						} else if selectedOption == "Attach to an existed player" {
							selectedPlayer, _ := pterm.
								DefaultInteractiveSelect.
								WithOptions(players).
								Show("Choose a player")
							sp := instance.FindPlayer(selectedPlayer)
							if sp == nil {
								errorMessageAndExit("Find player failed:", selectedPlayer)
							}
							sp.AddAlias(p.Name)
							res := instance.SetPlayer(sp)
							if res == nil {
								errorMessageAndExit("Attach player failed", selectedPlayer)
							} else {
								pterm.Printfln("Attaching alias %s to player %s", p.Name, sp.Name)
							}
						} else {
							errorMessageAndExit("Cancelled")
						}
					}
				}
			}

			pterm.Printfln("Importing \"%s\" `%s` ...", importEventName, importPath)
			event := *entity.NewEvent("temp", importEventName, importEventLevel)
			fn := fmt.Sprintf("%s.ktool", event.ID)
			event.Path = fn
			md5, err := util.MD5CheckSum(importPath)
			if err != nil {
				errorMessageAndExit(err)
			}
			for _, e := range instance.Conf.Events {
				if e.MD5 == md5 {
					errorMessageAndExit("Duplicated event found:", e.ID)
				}
			}
			event.MD5 = md5
			err = util.CopyFile(importPath, filepath.Join(instance.DataPath(), fn))
			if err != nil {
				errorMessageAndExit(err)
			}
			instance.AddEvent(event)
			eventsAdded++
		}
		err := instance.WriteConf()
		if err != nil {
			errorMessageAndExit(err)
		}
		pterm.Printfln("%d event(s) imported", eventsAdded)
	},
}
