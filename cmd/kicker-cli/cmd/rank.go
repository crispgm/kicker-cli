package cmd

import (
	"os"
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/app"
	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	monsterdyp "github.com/crispgm/kicker-cli/internal/operator/monster_dyp"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	rankGameMode  string
	rankEventName string
	allEvents     bool
)

func init() {
	rankCmd.Flags().StringVarP(&rankGameMode, "mode", "m", "", "Rank mode")
	rankCmd.Flags().StringVarP(&rankEventName, "name", "n", "", "Event name")
	rankCmd.Flags().BoolVarP(&allEvents, "all", "a", false, "Rank all events")
	rootCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Get rank for target event",
	Long:  "Get rank for target event",
	Run: func(cmd *cobra.Command, args []string) {
		instance := app.NewApp(initPath, app.DefaultName)
		err := instance.LoadConf()
		if err != nil {
			pterm.Error.Println("Not a valid kicker workspace")
			os.Exit(1)
		}
		var files []string
		if allEvents {
			for _, e := range instance.Conf.Events {
				files = append(files, e.Path)
			}
		} else {
			e := instance.GetEvent(rankEventName)
			if e == nil {
				pterm.Error.Println("Event not found")
				os.Exit(1)
			}
		}

		// load tournaments
		var tournaments []model.Tournament
		pterm.Info.Println("Loading tournaments ...")
		for _, p := range files {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), p))
			if err != nil {
				pterm.Error.Println(err)
				os.Exit(1)
			}
			tournaments = append(tournaments, *t)
		}
		c := converter.NewConverter()
		games, err := c.Normalize(tournaments, instance.Conf.Players)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}

		// calculating
		var statOperator operator.BaseOperator
		option := operator.Option{
			OrderBy:          "wr",
			RankMinThreshold: 5,
			EloKFactor:       eloKFactor,
			WithTime:         false,
			WithHomeAway:     false,
			WithPoint:        false,
			Incremental:      false,
		}
		if rankGameMode == entity.ModeMonsterDYPPlayerStats {
			statOperator = monsterdyp.NewPlayerStats(games, instance.Conf.Players, option)
		} else if rankGameMode == entity.ModeMonsterDYPTeamStats {
			statOperator = monsterdyp.NewTeamStats(games, option)
		}
		pterm.Info.Println("Briefing:", c.Briefing())
		table := statOperator.Output()
		if !dryRun {
			pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
		}
	},
}
