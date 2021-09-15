package monsterdyp

import (
	"fmt"
	"sort"

	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/stat"
)

// MultipleTournamentTeamStats generate statistics data of multiple monster DYP tournaments by team
type MultipleTournamentTeamStats struct {
	option stat.Option
	games  []model.EntityGame
}

// NewMultipleTournamentTeamStats .
func NewMultipleTournamentTeamStats(games []model.EntityGame, option stat.Option) *MultipleTournamentTeamStats {
	return &MultipleTournamentTeamStats{
		option: option,
		games:  games,
	}
}

// ValidMode .
func (m MultipleTournamentTeamStats) ValidMode(mode string) bool {
	return mode == model.ModeMonsterDYP
}

// Output .
func (m *MultipleTournamentTeamStats) Output() [][]string {
	data := make(map[string]model.EntityTeam)
	for _, g := range m.games {
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
		var et1, et2 model.EntityTeam
		if t, ok := data[team1Name]; ok {
			et1 = t
		} else {
			et1 = model.EntityTeam{
				Player1: t1p1Name,
				Player2: t1p2Name,
			}
		}
		if t, ok := data[team2Name]; ok {
			et2 = t
		} else {
			et2 = model.EntityTeam{
				Player1: t2p1Name,
				Player2: t2p2Name,
			}
		}
		timePlayed := g.TimePlayed
		et1.Played++
		et2.Played++
		et1.TimePlayed += timePlayed
		et2.TimePlayed += timePlayed
		m.playedTimeStats(&et1, timePlayed)
		m.playedTimeStats(&et2, timePlayed)

		if g.Point1 > g.Point2 {
			et1.Won++
			et2.Lost++
			et1.GoalsWon += (g.Point1 - g.Point2)
			et2.GoalsInLost += (g.Point1 - g.Point2)
		} else if g.Point1 < g.Point2 {
			et1.Lost++
			et2.Won++
			et2.GoalsWon += (g.Point2 - g.Point1)
			et1.GoalsInLost += (g.Point2 - g.Point1)
		} else {
			et1.Draws++
			et2.Draws++
		}
		et1.Goals += g.Point1
		et2.Goals += g.Point2
		et1.GoalsIn += g.Point2
		et2.GoalsIn += g.Point1

		data[team1Name] = et1
		data[team2Name] = et2
	}

	var sliceData []model.EntityTeam
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Won) / float32(d.Played) * 100.0
			d.PointsPerGame = float32(d.Goals) / float32(d.Played)
			d.PointsInPerGame = float32(d.GoalsIn) / float32(d.Played)
			d.TimePerGame = d.TimePlayed / d.Played / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			d.DiffPerWon = float32(d.GoalsWon) / float32(d.Won)
			d.DiffPerLost = float32(d.GoalsInLost) / float32(d.Lost)
		}
		sliceData = append(sliceData, d)
	}
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= m.option.RankMinThreshold && sliceData[j].Played < m.option.RankMinThreshold {
			return true
		}
		if sliceData[i].Played < m.option.RankMinThreshold && sliceData[j].Played >= m.option.RankMinThreshold {
			return false
		}

		if sliceData[i].WinRate > sliceData[j].WinRate {
			return true
		} else if sliceData[i].WinRate == sliceData[j].WinRate {
			iWinLost := sliceData[i].Won - sliceData[i].Lost
			jWinLost := sliceData[j].Won - sliceData[j].Lost
			if iWinLost > jWinLost {
				return true
			} else if iWinLost == jWinLost {
				return sliceData[i].Won > sliceData[j].Won
			}
		}
		return false
	})

	header := []string{"#", "Name", "Num", "Won", "Lost", "G+", "G-", "G±", "WR%", "PPG", "LPG", "DPW", "DPL"}
	timeHeader := []string{"TPG", "LGP", "SGP"}
	if m.option.WithTime {
		header = append(header, timeHeader...)
	}
	table := [][]string{header}
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
			fmt.Sprintf("%d", d.Won),
			fmt.Sprintf("%d", d.Lost),
			fmt.Sprintf("%d", d.Goals),
			fmt.Sprintf("%d", d.GoalsIn),
			goalDiff,
			winRate,
			fmt.Sprintf("%.2f", d.PointsPerGame),
			fmt.Sprintf("%.2f", d.PointsInPerGame),
			fmt.Sprintf("%.2f", d.DiffPerWon),
			fmt.Sprintf("%.2f", d.DiffPerLost),
		}
		if m.option.WithTime {
			item = append(item, []string{
				fmt.Sprintf("%02d:%02d", d.TimePerGame/60, d.TimePerGame%60),
				fmt.Sprintf("%02d:%02d", d.LongestGameTime/60, d.LongestGameTime%60),
				fmt.Sprintf("%02d:%02d", d.ShortestGameTime/60, d.ShortestGameTime%60),
			}...)
		}
		table = append(table, item)
	}
	return table
}

func (MultipleTournamentTeamStats) playedTimeStats(data *model.EntityTeam, timePlayed int) {
	if timePlayed < 0 || timePlayed > 1000*60*15 {
		// consider illegal
		return
	}
	if data.LongestGameTime < timePlayed || data.LongestGameTime == 0 {
		data.LongestGameTime = timePlayed
	}
	if data.ShortestGameTime > timePlayed || data.ShortestGameTime == 0 {
		data.ShortestGameTime = timePlayed
	}
}
