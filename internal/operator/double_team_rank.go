package operator

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var _ Operator = (*DoubleTeamRank)(nil)

// DoubleTeamRank generate statistics data of multiple double tournaments by team
type DoubleTeamRank struct {
	options     Option
	tournaments []entity.Tournament
}

// SupportedFormats .
func (o DoubleTeamRank) SupportedFormats(trn *model.Tournament) bool {
	return openDoubleTournament(trn)
}

// Input .
func (o *DoubleTeamRank) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.options = options
}

// Output .
func (o *DoubleTeamRank) Output() {
	data := make(map[string]entity.Team)
	for _, trn := range o.tournaments {
		for _, g := range trn.Converted.AllGames {
			t1p1Name := g.Team1[0]
			t1p2Name := g.Team1[1]
			t2p1Name := g.Team2[0]
			t2p2Name := g.Team2[1]
			team1Name := fmt.Sprintf("%s_%s", t1p1Name, t1p2Name)
			if t1p1Name > t1p2Name {
				team1Name = fmt.Sprintf("%s_%s", t1p2Name, t1p1Name)
			}
			team2Name := fmt.Sprintf("%s_%s", t2p1Name, t2p2Name)
			if t2p1Name > t2p2Name {
				team2Name = fmt.Sprintf("%s_%s", t2p2Name, t2p1Name)
			}
			var et1, et2 entity.Team
			if t, ok := data[team1Name]; ok {
				et1 = t
			} else {
				et1 = entity.Team{
					Player1: t1p1Name,
					Player2: t1p2Name,
				}
			}
			if t, ok := data[team2Name]; ok {
				et2 = t
			} else {
				et2 = entity.Team{
					Player1: t2p1Name,
					Player2: t2p2Name,
				}
			}
			et1.Played++
			et2.Played++

			if g.Point1 > g.Point2 {
				et1.Win++
				et2.Loss++
			} else if g.Point1 < g.Point2 {
				et1.Loss++
				et2.Win++
			} else {
				et1.Draw++
				et2.Draw++
			}

			data[team1Name] = et1
			data[team2Name] = et2
		}
	}

	var sliceData []entity.Team
	for _, d := range data {
		if d.Played != 0 {
			d.WinRate = float32(d.Win) / float32(d.Played) * 100.0
			sliceData = append(sliceData, d)
		}
	}
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= o.options.MinimumPlayed && sliceData[j].Played < o.options.MinimumPlayed {
			return true
		}
		if sliceData[i].Played < o.options.MinimumPlayed && sliceData[j].Played >= o.options.MinimumPlayed {
			return false
		}

		if sliceData[i].WinRate > sliceData[j].WinRate {
			return true
		} else if sliceData[i].WinRate == sliceData[j].WinRate {
			iWinLoss := sliceData[i].Win - sliceData[i].Loss
			jWinLoss := sliceData[j].Win - sliceData[j].Loss
			if iWinLoss >= jWinLoss {
				return true
			}
		}
		return false
	})

	if o.options.Head > 0 && len(sliceData) > o.options.Head {
		sliceData = sliceData[:o.options.Head]
	} else if o.options.Tail > 0 && len(sliceData) > o.options.Tail {
		sliceData = sliceData[len(sliceData)-o.options.Tail:]
	}

	header := []string{"#", "Name", "Num", "Win", "Loss", "Draw", "WR%"}
	table := [][]string{}
	if o.options.WithHeader {
		table = append(table, header)
	}
	for i, d := range sliceData {
		if d.Played == 0 {
			continue
		}
		winRate := fmt.Sprintf("%.0f%%", d.WinRate)
		item := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%s/%s", d.Player1, d.Player2),
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
			winRate,
		}
		table = append(table, item)
	}
	pterm.DefaultTable.WithHasHeader(o.options.WithHeader).WithData(table).WithBoxed(o.options.WithBoxes).Render()
}
