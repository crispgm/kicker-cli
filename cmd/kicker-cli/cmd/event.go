package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	eventIDOrName string
	eventNameType string
	allEvents     bool
)

func init() {
	eventCmd.PersistentFlags().BoolVarP(&allEvents, "all", "a", false, "rank all events")
	eventCmd.PersistentFlags().StringVarP(&eventNameType, "name-type", "t", "", "name type (single, byp, dyp or monster_dyp)")
	eventCmd.AddCommand(eventListCmd)
	rootCmd.AddCommand(eventCmd)
}

var eventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"events"},
	Short:   "Manage events organized by your organization",
	Long:    "Manage events organized by your organization",
	Run:     eventListCommand,
}

var eventListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List events",
	Long:    "List events",
	Run:     eventListCommand,
}

func eventListCommand(cmd *cobra.Command, args []string) {
	instance := initInstanceAndLoadConf()
	// load tournaments
	var table [][]string
	header := []string{"ID", "Name", "Date Time", "Points", "Name Type", "Mode", "URL"}
	table = append(table, header)
	if len(args) > 0 {
		for _, arg := range args {
			if e := instance.GetEvent(arg); e != nil {
				t, err := parser.ParseFile(filepath.Join(instance.DataPath(), e.Path))
				if err != nil {
					errorMessageAndExit(err)
				}
				if len(t.Mode) > 0 && t.Mode == eventNameType {
					showInfo(&table, e, t)
				}
			}
		}
	} else {
		for _, e := range instance.Conf.Events {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), e.Path))
			if err != nil {
				errorMessageAndExit(err)
			}
			if len(t.Mode) > 0 && t.Mode == eventNameType {
				showInfo(&table, &e, t)
			}
		}
	}
	if len(table) == 1 {
		errorMessageAndExit("No event(s) found")
	}
	pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
}

func showInfo(table *[][]string, e *entity.Event, t *model.Tournament) {
	url := e.URL
	if url == "" {
		url = "-"
	}
	*table = append(*table, []string{
		e.ID,
		e.Name,
		t.Created.Format("2006-01-02 15:04"),
		fmt.Sprintf("%d", e.Points),
		t.NameType,
		t.Mode,
		url,
	})
}
