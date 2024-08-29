package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	eventNameTypes []string
	eventAfter     string
	eventBefore    string
)

func init() {
	eventCmd.PersistentFlags().StringArrayVarP(&eventNameTypes, "name-type", "t", []string{}, "name type (single, byp, dyp or monster_dyp)")
	eventCmd.PersistentFlags().StringVarP(&eventAfter, "after", "a", "", "show events created after a specific date")
	eventCmd.PersistentFlags().StringVarP(&eventBefore, "before", "b", "", "show events created before a specific date")
	eventCmd.AddCommand(eventListCmd)
	rootCmd.AddCommand(eventCmd)
}

var eventCmd = &cobra.Command{
	Use:     "event",
	Aliases: []string{"events", "ev"},
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
	dataPath := instance.DataPath()
	// load tournaments
	table := initEventInfoHeader()
	if len(args) > 0 {
		for _, arg := range args {
			if e := instance.GetEvent(arg); e != nil {
				loadAndShowEventInfo(&table, dataPath, instance.Conf.Players, e, instance.Conf.Organization.Timezone)
			}
		}
	} else {
		for _, e := range instance.Conf.Events {
			loadAndShowEventInfo(&table, dataPath, instance.Conf.Players, &e, instance.Conf.Organization.Timezone)
		}
	}
	if len(table) <= 1 {
		errorMessageAndExit("No event(s) found")
	}
	_ = pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
}

func nameTypeIncluded(input string) bool {
	for _, t := range eventNameTypes {
		if t == input {
			return true
		}
	}
	return false
}

func createdBetween(created time.Time) bool {
	if len(eventAfter) > 0 {
		after, err := dateparse.ParseLocal(eventAfter)
		if err != nil {
			errorMessageAndExit(err)
		}
		if created.Before(after) {
			return false
		}
	}
	if len(eventBefore) > 0 {
		before, err := dateparse.ParseLocal(eventBefore)
		if err != nil {
			errorMessageAndExit(err)
		}
		if created.After(before) {
			return false
		}
	}

	return true
}

func initEventInfoHeader() [][]string {
	var table [][]string
	header := []string{"ID", "Name", "Date Time", "Level", "Players", "Games", "Rounds", "Name Type", "Mode", "URL"}
	if !globalNoHeaders {
		table = append(table, header)
	}
	return table
}

func loadAndShowEventInfo(table *[][]string, dataPath string, players []entity.Player, e *entity.Event, tz string) (*model.Tournament, *entity.Record) {
	t, r, err := loadEventInfo(dataPath, players, e)
	if err != nil {
		errorMessageAndExit(err)
	}

	showEvent(table, e, t, r, tz)
	return t, r
}

func loadEventInfo(dataPath string, players []entity.Player, e *entity.Event) (*model.Tournament, *entity.Record, error) {
	t, err := parser.ParseFile(filepath.Join(dataPath, e.Path))
	if err != nil {
		return nil, nil, err
	}
	c := converter.NewConverter()
	trn, err := c.Normalize(players, *t)
	if err != nil {
		return nil, nil, err
	}

	return t, trn, nil
}

func showEvent(table *[][]string, e *entity.Event, t *model.Tournament, r *entity.Record, tz string) {
	if len(eventNameTypes) > 0 && !nameTypeIncluded(t.NameType) {
		return
	}
	if !createdBetween(t.Created) {
		return
	}

	showInfo(table, e, t, r, tz)
}

func showInfo(table *[][]string, e *entity.Event, t *model.Tournament, r *entity.Record, tz string) {
	var levels []string
	if len(e.ITSFLevel) > 0 {
		levels = append(levels, e.ITSFLevel)
	}
	if len(e.ATSALevel) > 0 {
		levels = append(levels, e.ATSALevel)
	}
	if len(e.KickerLevel) > 0 {
		levels = append(levels, e.KickerLevel)
	}
	loc, _ := time.LoadLocation(tz)
	*table = append(*table, []string{
		e.ID,
		e.Name,
		t.Created.In(loc).Format("2006-01-02 15:04"),
		strings.Join(levels, "|"),
		fmt.Sprintf("%d", len(r.Players)),
		fmt.Sprintf("%d", len(r.AllGames)),
		fmt.Sprintf("%d", len(t.Rounds)),
		t.NameType,
		t.TournamentMode(),
		dashIfEmpty(e.URL),
	})
}
