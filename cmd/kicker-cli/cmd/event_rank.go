package cmd

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/internal/operator/double"
	"github.com/crispgm/kicker-cli/internal/operator/single"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
	"github.com/crispgm/kicker-cli/pkg/rating/elo"
)

var (
	rankGameMode   string
	rankMinPlayed  int
	rankHead       int
	rankTail       int
	rankELOKFactor int
	rankOrderBy    string
	rankWithGoals  bool
)

func init() {
	rankCmd.Flags().StringVarP(&rankGameMode, "mode", "m", "", "rank mode")
	rankCmd.Flags().StringVarP(&rankOrderBy, "order-by", "o", "wr", "order by (wr/elo)")
	rankCmd.Flags().IntVarP(&rankMinPlayed, "minimum-played", "p", 0, "minimum matches played")
	rankCmd.Flags().BoolVarP(&rankWithGoals, "with-goals", "", false, "rank with goals")
	rankCmd.Flags().IntVarP(&rankELOKFactor, "elo-k", "k", elo.K, "K factor")
	rankCmd.Flags().IntVarP(&rankHead, "head", "", 0, "display the head part of rank")
	rankCmd.Flags().IntVarP(&rankTail, "tail", "", 0, "display the last part of rank")
	rankCmd.MarkFlagRequired("mode")
	rankCmd.MarkFlagsMutuallyExclusive("head", "tail")
	eventCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Get rank",
	Long:  "Get rank",
	Run: func(cmd *cobra.Command, args []string) {
		if rankHead < 0 || rankTail < 0 {
			errorMessageAndExit("Only non-negitive number is allowed for head or tail")
		}
		var op operator.Operator
		switch rankGameMode {
		case entity.ModeDoublePlayerRank:
			op = &double.PlayerRank{}
		case entity.ModeDoubleTeamRank:
			op = &double.TeamRank{}
		case entity.ModeDoubleTeamRival:
			op = &double.TeamRival{}
		case entity.ModeSinglePlayerRank:
			op = &single.PlayerRank{}
		case entity.ModeSinglePlayerRival:
			op = &single.PlayerRival{}
		default:
			errorMessageAndExit("Please present a valid rank mode")
		}

		instance := initInstanceAndLoadConf()

		var events []entity.Event
		if allEvents {
			for _, e := range instance.Conf.Events {
				events = append(events, e)
			}
		} else if len(args) > 0 {
			for _, arg := range args {
				e := instance.GetEvent(arg)
				if e != nil {
					events = append(events, *e)
				} else {
					errorMessageAndExit("Event", arg, "not found")
				}
			}
		}

		// load tournaments
		var tournaments []model.Tournament
		for _, e := range events {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), e.Path))
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
			if !op.SupportedFormats(t) {
				pterm.Warning.Println("Not supported by operator. Ignoring", e.ID)
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
			OrderBy:       rankOrderBy,
			MinimumPlayed: rankMinPlayed,
			Head:          rankHead,
			Tail:          rankTail,
			EloKFactor:    rankELOKFactor,
			WithHeader:    !globalNoHeaders,
			WithGoals:     rankWithGoals,
		}

		pterm.Println("Briefing:", c.Briefing())
		pterm.Println()
		op.Input(trn.AllGames, instance.Conf.Players, options)
		table := op.Output()
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	},
}
