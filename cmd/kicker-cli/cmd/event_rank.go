package cmd

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/internal/operator/double"
	"github.com/crispgm/kicker-cli/pkg/class/elo"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	rankGameMode   string
	rankNameType   string
	rankAllEvents  bool
	rankMinPlayed  int
	rankELOKFactor int
	rankWithTime   bool
	rankWithGoals  bool
)

func init() {
	rankCmd.Flags().BoolVarP(&rankAllEvents, "all", "a", false, "rank all events")
	rankCmd.Flags().IntVarP(&rankMinPlayed, "minimum-played", "m", 5, "minimum matches played")
	rankCmd.Flags().IntVarP(&rankELOKFactor, "elo-k", "k", elo.K, "K factor")
	rankCmd.Flags().BoolVarP(&rankWithTime, "with-time", "", false, "rank with time duration")
	rankCmd.Flags().BoolVarP(&rankWithGoals, "with-goals", "", false, "rank with goals")
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
		var op operator.Operator
		options := operator.Option{
			OrderBy:          "wr",
			RankMinThreshold: rankMinPlayed,
			EloKFactor:       rankELOKFactor,
			WithHomeAway:     false,
			WithTime:         rankWithTime,
			WithGoals:        rankWithGoals,
		}

		switch rankGameMode {
		case entity.ModeDoublePlayerRanks:
			op = &double.PlayerRanks{}
		case entity.ModeDoubleTeamRanks:
			op = &double.TeamRanks{}
		// case entity.ModeDoubleTeamRivals:
		// case entity.ModeSinglePlayerRanks:
		// case entity.ModeSinglePlayerRivals:
		default:
			errorMessageAndExit("Please present a valid rank mode")
		}

		op.Input(games, instance.Conf.Players, options)
		pterm.Info.Println("Briefing:", c.Briefing())
		table := op.Output()
		pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
	},
}
