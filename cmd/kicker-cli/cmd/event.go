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
	eventCmd.MarkFlagsMutuallyExclusive("all", "name-type")
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
				showEvent(instance.DataPath(), e, &table)
			}
		}
	} else {
		for _, e := range instance.Conf.Events {
			showEvent(instance.DataPath(), &e, &table)
		}
	}
	if len(table) == 1 {
		errorMessageAndExit("No event(s) found")
	}
	pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
}

func showEvent(dataPath string, e *entity.Event, table *[][]string) {
	t, err := parser.ParseFile(filepath.Join(dataPath, e.Path))
	if err != nil {
		errorMessageAndExit(err)
	}
	if !allEvents && (len(eventNameType) > 0 && t.NameType != eventNameType) {
		return
	}

	showInfo(table, e, t)
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
