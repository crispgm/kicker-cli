package operator

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var _ Operator = (*SinglePlayerRival)(nil)

// SinglePlayerRival generate statistics data of multiple monster DYP tournaments by team
type SinglePlayerRival struct {
	options     Option
	tournaments []entity.Tournament
}

// SupportedFormats .
func (o SinglePlayerRival) SupportedFormats(trn *model.Tournament) bool {
	return openSingleTournament(trn)
}

// Input .
func (o *SinglePlayerRival) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.options = options
}

// Output .
func (o *SinglePlayerRival) Output() {
	data := make(map[string]entity.Rival)
	for _, t := range o.tournaments {
		for _, g := range t.Converted.AllGames {
			p1Name := g.Team1[0]
			p2Name := g.Team2[0]
			rivalName := fmt.Sprintf("%s_vs_%s", p1Name, p2Name)
			rivalNameAlt := fmt.Sprintf("%s_vs_%s", p2Name, p1Name)
			reversed := false

			var rival entity.Rival
			if _, ok := data[rivalName]; ok {
				rival = data[rivalName]
			} else if _, ok := data[rivalNameAlt]; ok {
				rivalName = rivalNameAlt
				rival = data[rivalNameAlt]
				reversed = true
			} else {
				rival = entity.Rival{
					Team1: entity.Team{
						Player1: p1Name,
					},
					Team2: entity.Team{
						Player1: p2Name,
					},
				}
			}

			rival.Played++
			rival.Team1.Played++
			rival.Team2.Played++
			if !reversed {
				if g.Point1 > g.Point2 {
					rival.Win++
					rival.Team1.Win++
					rival.Team2.Loss++
				} else if g.Point1 < g.Point2 {
					rival.Loss++
					rival.Team2.Win++
					rival.Team1.Loss++
				} else {
					rival.Draw++
					rival.Team1.Draw++
					rival.Team2.Draw++
				}
			} else {
				if g.Point1 < g.Point2 {
					rival.Win++
					rival.Team1.Win++
					rival.Team2.Loss++
				} else if g.Point1 > g.Point2 {
					rival.Loss++
					rival.Team2.Win++
					rival.Team1.Loss++
				} else {
					rival.Draw++
					rival.Team1.Draw++
					rival.Team2.Draw++
				}
			}

			data[rivalName] = rival
		}
	}

	var sliceData []entity.Rival
	for _, d := range data {
		sliceData = append(sliceData, d)
	}

	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= o.options.MinimumPlayed && sliceData[j].Played < o.options.MinimumPlayed {
			return true
		}
		if sliceData[i].Played < o.options.MinimumPlayed && sliceData[j].Played >= o.options.MinimumPlayed {
			return false
		}

		if sliceData[i].Played > sliceData[j].Played {
			return true
		}
		return false
	})

	if o.options.Head > 0 && len(sliceData) > o.options.Head {
		sliceData = sliceData[:o.options.Head]
	} else if o.options.Tail > 0 && len(sliceData) > o.options.Tail {
		sliceData = sliceData[len(sliceData)-o.options.Tail:]
	}

	header := []string{"#", "Team1", "Team2", "Num", "Win", "Loss", "Draw"}
	table := [][]string{}
	if o.options.WithHeader {
		table = append(table, header)
	}
	for i, d := range sliceData {
		if d.Played == 0 {
			continue
		}
		item := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%s", d.Team1.Player1),
			fmt.Sprintf("%s", d.Team2.Player1),
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
		}
		table = append(table, item)
	}
	pterm.DefaultTable.WithHasHeader(o.options.WithHeader).WithData(table).WithBoxed(o.options.WithBoxes).Render()
}
