package cmd

import (
	"path/filepath"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/converter"
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

var (
	rankGameMode  string
	rankMinPlayed int
	rankHead      int
	rankTail      int
	rankSortBy    string
)

func init() {
	rankCmd.Flags().StringVarP(&rankGameMode, "mode", "m", "", "rank mode")
	rankCmd.Flags().StringVarP(&rankSortBy, "sort-by", "o", "krs", "sort by (krs/itsf/atsa/elo/wr)")
	rankCmd.Flags().IntVarP(&rankMinPlayed, "minimum-played", "p", 0, "minimum matches played")
	rankCmd.Flags().IntVarP(&rankHead, "head", "", 0, "display the head part of rank")
	rankCmd.Flags().IntVarP(&rankTail, "tail", "", 0, "display the last part of rank")
	rankCmd.MarkFlagRequired("mode")
	rankCmd.MarkFlagsMutuallyExclusive("head", "tail")
	eventCmd.AddCommand(rankCmd)
}

var rankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Get player ranks",
	Long:  "Get player ranks of played tournaments and games",
	Run: func(cmd *cobra.Command, args []string) {
		if rankHead < 0 || rankTail < 0 {
			errorMessageAndExit("Only non-negitive number is allowed for head or tail")
		}
		var op operator.Operator
		switch rankGameMode {
		case entity.ModeDoublePlayerRank:
			op = &operator.DoublePlayerRank{}
		case entity.ModeDoubleTeamRank:
			op = &operator.DoubleTeamRank{}
		case entity.ModeDoubleTeamRival:
			op = &operator.DoubleTeamRival{}
		case entity.ModeSinglePlayerRank:
			op = &operator.SinglePlayerRank{}
		case entity.ModeSinglePlayerRival:
			op = &operator.SinglePlayerRival{}
		default:
			errorMessageAndExit("Please present a valid rank mode")
		}

		instance := initInstanceAndLoadConf()

		var events []entity.Event
		if len(args) > 0 {
			for _, arg := range args {
				e := instance.GetEvent(arg)
				if e != nil {
					events = append(events, *e)
				} else {
					errorMessageAndExit("Event", arg, "not found")
				}
			}
		} else {
			for _, e := range instance.Conf.Events {
				events = append(events, e)
			}
		}

		// load tournaments
		var tournaments []model.Tournament
		var filteredEvents []entity.Event
		for _, e := range events {
			t, err := parser.ParseFile(filepath.Join(instance.DataPath(), e.Path))
			if err != nil {
				errorMessageAndExit(err)
			}
			if len(eventNameTypes) == 0 {
				// choose the first file as name type if it's not set
				eventNameTypes = []string{t.NameType}
			}
			if !nameTypeIncluded(t.NameType) {
				continue
			}
			if !createdBetween(t.Created) {
				continue
			}
			if !op.SupportedFormats(t) {
				pterm.Warning.Println("Not supported by operator. Ignoring", e.ID)
				continue
			}
			tournaments = append(tournaments, *t)
			filteredEvents = append(filteredEvents, e)
		}
		if len(tournaments) == 0 {
			pterm.Warning.Println("No matched tournament(s)")
			return
		}

		var eTournaments []entity.Tournament
		for i, t := range tournaments {
			c := converter.NewConverter()
			trn, err := c.Normalize(instance.Conf.Players, t)
			if err != nil {
				errorMessageAndExit(err)
			}

			eTournaments = append(eTournaments, entity.Tournament{
				Event:     filteredEvents[i],
				Raw:       t,
				Converted: *trn,
			})
		}
		options := operator.Option{
			OrderBy:       rankSortBy,
			MinimumPlayed: rankMinPlayed,
			Head:          rankHead,
			Tail:          rankTail,
			WithHeader:    !globalNoHeaders,
		}
		op.Input(eTournaments, instance.Conf.Players, options)
		table := op.Output()
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	},
}
