package double

import (
	"fmt"
	"sort"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var _ operator.Operator = (*TeamRank)(nil)

// TeamRank generate statistics data of multiple double tournaments by team
type TeamRank struct {
	options     operator.Option
	tournaments []entity.Tournament
}

// SupportedFormats .
func (t TeamRank) SupportedFormats(trn *model.Tournament) bool {
	if trn.IsDouble() {
		if trn.Mode == model.ModeMonsterDYP ||
			trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}

// Input .
func (t *TeamRank) Input(tournaments []entity.Tournament, players []entity.Player, options operator.Option) {
	t.tournaments = tournaments
	t.options = options
}

// Output .
func (t *TeamRank) Output() [][]string {
	data := make(map[string]entity.Team)
	for _, trn := range t.tournaments {
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
			timePlayed := g.TimePlayed
			et1.Played++
			et2.Played++
			et1.TimePlayed += timePlayed
			et2.TimePlayed += timePlayed

			if g.Point1 > g.Point2 {
				et1.Win++
				et2.Loss++
				et1.GoalsWin += (g.Point1 - g.Point2)
				et2.GoalsInLoss += (g.Point1 - g.Point2)
			} else if g.Point1 < g.Point2 {
				et1.Loss++
				et2.Win++
				et2.GoalsWin += (g.Point2 - g.Point1)
				et1.GoalsInLoss += (g.Point2 - g.Point1)
			} else {
				et1.Draw++
				et2.Draw++
			}
			et1.Goals += g.Point1
			et2.Goals += g.Point2
			et1.GoalsIn += g.Point2
			et2.GoalsIn += g.Point1

			data[team1Name] = et1
			data[team2Name] = et2
		}
	}

	var sliceData []entity.Team
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Win) / float32(d.Played) * 100.0
			d.PointsPerGame = float32(d.Goals) / float32(d.Played)
			d.PointsInPerGame = float32(d.GoalsIn) / float32(d.Played)
			d.TimePerGame = d.TimePlayed / d.Played / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			d.DiffPerWin = float32(d.GoalsWin) / float32(d.Win)
			d.DiffPerLoss = float32(d.GoalsInLoss) / float32(d.Loss)
			sliceData = append(sliceData, d)
		}
	}
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= t.options.MinimumPlayed && sliceData[j].Played < t.options.MinimumPlayed {
			return true
		}
		if sliceData[i].Played < t.options.MinimumPlayed && sliceData[j].Played >= t.options.MinimumPlayed {
			return false
		}

		if sliceData[i].WinRate > sliceData[j].WinRate {
			return true
		} else if sliceData[i].WinRate == sliceData[j].WinRate {
			iWinLoss := sliceData[i].Win - sliceData[i].Loss
			jWinLoss := sliceData[j].Win - sliceData[j].Loss
			if iWinLoss > jWinLoss {
				return true
			} else if iWinLoss == jWinLoss {
				if sliceData[i].GoalDiff > sliceData[j].GoalDiff {
					return true
				} else if sliceData[i].GoalDiff == sliceData[j].GoalDiff {
					return sliceData[i].Goals > sliceData[j].Goals
				}
			}
		}
		return false
	})

	if t.options.Head > 0 && len(sliceData) > t.options.Head {
		sliceData = sliceData[:t.options.Head]
	} else if t.options.Tail > 0 && len(sliceData) > t.options.Tail {
		sliceData = sliceData[len(sliceData)-t.options.Tail:]
	}

	header := []string{"#", "Name", "Num", "Win", "Loss", "Draw", "WR%"}
	pointHeader := []string{"G+", "G-", "G±", "PPG", "LPG", "DPW", "DPL"}
	if t.options.WithGoals {
		header = append(header, pointHeader...)
	}
	table := [][]string{}
	if t.options.WithHeader {
		table = append(table, header)
	}
	for i, d := range sliceData {
		if d.Played == 0 {
			continue
		}
		goalDiff := fmt.Sprintf("%d", d.GoalDiff)
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
		if t.options.WithGoals {
			item = append(item, []string{
				fmt.Sprintf("%d", d.Goals),
				fmt.Sprintf("%d", d.GoalsIn),
				goalDiff,
				fmt.Sprintf("%.2f", d.PointsPerGame),
				fmt.Sprintf("%.2f", d.PointsInPerGame),
				fmt.Sprintf("%.2f", d.DiffPerWin),
				fmt.Sprintf("%.2f", d.DiffPerLoss),
			}...)
		}
		table = append(table, item)
	}
	return table
}
