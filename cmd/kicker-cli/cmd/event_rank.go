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
	"github.com/crispgm/kicker-cli/pkg/rating/elo"
)

var (
	rankMinPlayed  int
	rankELOKFactor int
	rankWithTime   bool
	rankWithGoals  bool
)

func init() {
	rankCmd.Flags().IntVarP(&rankMinPlayed, "minimum-played", "", 0, "minimum matches played")
	rankCmd.Flags().BoolVarP(&rankWithGoals, "with-goals", "", false, "rank with goals")
	rankCmd.Flags().BoolVarP(&rankWithTime, "with-time", "", false, "rank with time duration")
	rankCmd.Flags().IntVarP(&rankELOKFactor, "elo-k", "k", elo.K, "K factor")
	eventCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Get rank",
	Long:  "Get rank",
	Run: func(cmd *cobra.Command, args []string) {
		var op operator.Operator
		switch eventGameMode {
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

		instance := initInstanceAndLoadConf()

		var files []string
		if allEvents {
			for _, e := range instance.Conf.Events {
				files = append(files, e.Path)
			}
		} else if len(args) > 0 {
			for _, arg := range args {
				e := instance.GetEvent(arg)
				if e != nil {
					files = append(files, e.Path)
				} else {
					errorMessageAndExit("Event", arg, "not found")
				}
			}
		}

		// load tournaments
		var tournaments []model.Tournament
		for _, p := range files {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), p))
			if err != nil {
				errorMessageAndExit(err)
			}
			if eventNameType == "" {
				// choose the first file as name type if it's not set
				eventNameType = t.NameType
			}
			if t.NameType != eventNameType {
				continue
			}
			tournaments = append(tournaments, *t)
		}
		if len(tournaments) == 0 {
			pterm.Warning.Println("No matched tournament(s)")
			return
		}

		pterm.Println("Loading tournaments ...")
		c := converter.NewConverter()
		trn, err := c.Normalize(tournaments, instance.Conf.Players)
		if err != nil {
			errorMessageAndExit(err)
		}

		// calculating
		options := operator.Option{
			OrderBy:          "wr",
			RankMinThreshold: rankMinPlayed,
			EloKFactor:       rankELOKFactor,
			WithHeader:       !globalNoHeaders,
			WithHomeAway:     false,
			WithTime:         rankWithTime,
			WithGoals:        rankWithGoals,
		}

		pterm.Println("Briefing:", c.Briefing())
		pterm.Println()
		op.Input(trn.AllGames, instance.Conf.Players, options)
		table := op.Output()
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	},
}
