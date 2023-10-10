package operator

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var _ Operator = (*DoubleTeamRival)(nil)

// DoubleTeamRival generate statistics data of multiple monster DYP tournaments by team
type DoubleTeamRival struct {
	options     Option
	tournaments []entity.Tournament
}

// SupportedFormats .
func (o DoubleTeamRival) SupportedFormats(trn *model.Tournament) bool {
	return openDoubleTournament(trn)
}

// Input .
func (o *DoubleTeamRival) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.options = options
}

// Output .
func (o *DoubleTeamRival) Output() {
	data := make(map[string]entity.Rival)
	for _, trn := range o.tournaments {
		for _, g := range trn.Converted.AllGames {
			t1p1Name := g.Team1[0]
			t1p2Name := g.Team1[1]
			t2p1Name := g.Team2[0]
			t2p2Name := g.Team2[1]
			team1Name := fmt.Sprintf("%s/%s", t1p1Name, t1p2Name)
			if t1p1Name > t1p2Name {
				team1Name = fmt.Sprintf("%s/%s", t1p2Name, t1p1Name)
			}
			team2Name := fmt.Sprintf("%s/%s", t2p1Name, t2p2Name)
			if t2p1Name > t2p2Name {
				team2Name = fmt.Sprintf("%s/%s", t2p2Name, t2p1Name)
			}

			rivalName := fmt.Sprintf("%s_vs_%s", team1Name, team2Name)
			rivalNameAlt := fmt.Sprintf("%s_vs_%s", team2Name, team1Name)
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
						Player1: t1p1Name,
						Player2: t1p2Name,
					},
					Team2: entity.Team{
						Player1: t2p1Name,
						Player2: t2p2Name,
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
			fmt.Sprintf("%s/%s", d.Team1.Player1, d.Team1.Player2),
			fmt.Sprintf("%s/%s", d.Team2.Player1, d.Team2.Player2),
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
		}
		table = append(table, item)
	}

	_ = pterm.DefaultTable.WithHasHeader(o.options.WithHeader).WithData(table).WithBoxed(o.options.WithBoxes).Render()
}
