package cmd

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/internal/operator/double"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	rankGameMode  string
	rankNameType  string
	rankAllEvents bool
)

func init() {
	rankCmd.Flags().BoolVarP(&rankAllEvents, "all", "a", false, "rank all events")
	eventCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Get rank",
	Long:  "Get rank",
	Run: func(cmd *cobra.Command, args []string) {
		instance := initInstanceAndLoadConf()

		var files []string
		if rankAllEvents {
			for _, e := range instance.Conf.Events {
				files = append(files, e.Path)
			}
		} else {
			e := instance.GetEvent(eventIDOrName)
			if e == nil {
				errorMessageAndExit("No event(s) found.")
			}
		}

		// load tournaments
		var tournaments []model.Tournament
		pterm.Info.Println("Loading tournaments ...")
		for _, p := range files {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), p))
			if err != nil {
				errorMessageAndExit(err)
			}
			if t.NameType != rankNameType {
				continue
			}
			tournaments = append(tournaments, *t)
		}
		c := converter.NewConverter()
		games, err := c.Normalize(tournaments, instance.Conf.Players)
		if err != nil {
			errorMessageAndExit(err)
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
		}
		if rankGameMode == entity.ModeDoublePlayerRanks {
			statOperator = double.NewPlayerStats(games, instance.Conf.Players, option)
		} else if rankGameMode == entity.ModeDoubleTeamRanks {
			statOperator = double.NewTeamStats(games, option)
		}
		pterm.Info.Println("Briefing:", c.Briefing())
		table := statOperator.Output()
		pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
	},
}
