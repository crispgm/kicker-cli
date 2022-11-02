package monsterdyp

import (
	"fmt"
	"sort"

	"github.com/crispgm/kickertool-analyzer/elo"
	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/operator"
)

// PlayerStats generate statistics data of multiple monster DYP tournaments
type PlayerStats struct {
	option  operator.Option
	games   []model.EntityGame
	players []model.EntityPlayer
}

// NewPlayerStats .
func NewPlayerStats(games []model.EntityGame, players []model.EntityPlayer, option operator.Option) *PlayerStats {
	return &PlayerStats{
		games:   games,
		players: players,
		option:  option,
	}
}

// ValidMode .
func (p PlayerStats) ValidMode(mode string) bool {
	return mode == model.ModeMonsterDYP
}

// Output .
func (p *PlayerStats) Output() [][]string {
	data := make(map[string]model.EntityPlayer)
	for _, p := range p.players {
		data[p.Name] = p
	}
	for _, g := range p.games {
		t1p1Data := data[g.Team1[0]]
		t1p2Data := data[g.Team1[1]]
		t2p1Data := data[g.Team2[0]]
		t2p2Data := data[g.Team2[1]]
		t1p1Data.Name = g.Team1[0]
		t1p2Data.Name = g.Team1[1]
		t2p1Data.Name = g.Team2[0]
		t2p2Data.Name = g.Team2[1]
		t1p1Data.Played++
		t1p2Data.Played++
		t2p1Data.Played++
		t2p2Data.Played++
		t1p1Data.TimePlayed += g.TimePlayed
		t1p2Data.TimePlayed += g.TimePlayed
		t2p1Data.TimePlayed += g.TimePlayed
		t2p2Data.TimePlayed += g.TimePlayed
		p.playedTimeStats(&t1p1Data, g.TimePlayed)
		p.playedTimeStats(&t1p2Data, g.TimePlayed)
		p.playedTimeStats(&t2p1Data, g.TimePlayed)
		p.playedTimeStats(&t2p2Data, g.TimePlayed)
		if g.Point1 > g.Point2 {
			t1p1Data.Won++
			t1p2Data.Won++
			t2p1Data.Lost++
			t2p2Data.Lost++
			t1p1Data.HomeWon++
			t1p2Data.HomeWon++
			t2p1Data.AwayLost++
			t2p2Data.AwayLost++
			t1p1Data.GoalsWon += (g.Point1 - g.Point2)
			t1p2Data.GoalsWon += (g.Point1 - g.Point2)
			t2p1Data.GoalsInLost += (g.Point1 - g.Point2)
			t2p2Data.GoalsInLost += (g.Point1 - g.Point2)
		} else if g.Point2 > g.Point1 {
			t1p1Data.Lost++
			t1p2Data.Lost++
			t2p1Data.Won++
			t2p2Data.Won++
			t1p1Data.HomeLost++
			t1p2Data.HomeLost++
			t2p1Data.AwayWon++
			t2p2Data.AwayWon++
			t2p1Data.GoalsWon += (g.Point2 - g.Point1)
			t2p2Data.GoalsWon += (g.Point2 - g.Point1)
			t1p1Data.GoalsInLost += (g.Point2 - g.Point1)
			t1p2Data.GoalsInLost += (g.Point2 - g.Point1)
		} else {
			// basically not approachable
			t1p1Data.Draws++
			t1p2Data.Draws++
			t2p1Data.Draws++
			t2p2Data.Draws++
		}
		t1p1Data.Goals += g.Point1
		t1p2Data.Goals += g.Point1
		t2p1Data.Goals += g.Point2
		t2p2Data.Goals += g.Point2
		t1p1Data.GoalsIn += g.Point2
		t1p2Data.GoalsIn += g.Point2
		t2p1Data.GoalsIn += g.Point1
		t2p2Data.GoalsIn += g.Point1
		t1p1Elo := elo.InitialScore
		t1p2Elo := elo.InitialScore
		t2p1Elo := elo.InitialScore
		t2p2Elo := elo.InitialScore
		if t1p1Data.EloRating != 0 {
			t1p1Elo = t1p1Data.EloRating
		}
		if t1p2Data.EloRating != 0 {
			t1p2Elo = t1p2Data.EloRating
		}
		if t2p1Data.EloRating != 0 {
			t2p1Elo = t2p1Data.EloRating
		}
		if t2p2Data.EloRating != 0 {
			t2p2Elo = t2p2Data.EloRating
		}
		rate := elo.Rate{
			T1P1Score: t1p1Elo,
			T1P2Score: t1p2Elo,
			T2P1Score: t2p1Elo,
			T2P2Score: t2p2Elo,
			HostWin:   g.Point1 > g.Point2,
			K:         float64(p.option.EloKFactor),
		}
		rate.CalcEloRating()
		t1p1Data.EloRating = rate.T1P1Score
		t1p2Data.EloRating = rate.T1P2Score
		t2p1Data.EloRating = rate.T2P1Score
		t2p2Data.EloRating = rate.T2P2Score
		data[g.Team1[0]] = t1p1Data
		data[g.Team1[1]] = t1p2Data
		data[g.Team2[0]] = t2p1Data
		data[g.Team2[1]] = t2p2Data
	}

	var sliceData []model.EntityPlayer
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Won) / float32(d.Played) * 100.0
			if d.HomeWon+d.HomeLost > 0 {
				d.HomeWonRate = float32(d.HomeWon) / float32(d.HomeWon+d.HomeLost) * 100.0
			}
			if d.AwayWon+d.AwayLost > 0 {
				d.AwayWonRate = float32(d.AwayWon) / float32(d.AwayWon+d.AwayLost) * 100.0
			}
			d.PointsPerGame = float32(d.Goals) / float32(d.Played)
			d.PointsInPerGame = float32(d.GoalsIn) / float32(d.Played)
			d.TimePerGame = d.TimePlayed / d.Played / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			if d.Won > 0 {
				d.DiffPerWon = float32(d.GoalsWon) / float32(d.Won)
			}
			if d.Lost > 0 {
				d.DiffPerLost = float32(d.GoalsInLost) / float32(d.Lost)
			}
		}
		sliceData = append(sliceData, d)
	}
	p.players = sliceData
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= p.option.RankMinThreshold && sliceData[j].Played < p.option.RankMinThreshold {
			return true
		}
		if sliceData[i].Played < p.option.RankMinThreshold && sliceData[j].Played >= p.option.RankMinThreshold {
			return false
		}

		if p.option.OrderBy == "elo" {
			if sliceData[i].EloRating > sliceData[j].EloRating {
				return true
			}
		} else {
			if sliceData[i].WinRate > sliceData[j].WinRate {
				return true
			}
		}
		return false
	})

	header := []string{"#", "Name", "Elo", "Num", "Won", "Lost", "G+", "G-", "G±", "WR%"}
	haHeader := []string{"HW", "HL", "HW%", "AW", "AL", "AW%"}
	timeHeader := []string{"TPG", "LGP", "SGP"}
	pointHeader := []string{"PPG", "LPG", "DPW", "DPL"}
	if p.option.WithHomeAway {
		header = append(header, haHeader...)
	}
	if p.option.WithTime {
		header = append(header, timeHeader...)
	}
	if p.option.WithPoint {
		header = append(header, pointHeader...)
	}
	table := [][]string{header}
	for i, d := range sliceData {
		item := []string{
			fmt.Sprintf("%d", i+1),
			d.Name,
			fmt.Sprintf("%.0f", d.EloRating),
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Won),
			fmt.Sprintf("%d", d.Lost),
			fmt.Sprintf("%d", d.Goals),
			fmt.Sprintf("%d", d.GoalsIn),
			fmt.Sprintf("%d", d.GoalDiff),
			fmt.Sprintf("%.0f%%", d.WinRate),
		}
		if p.option.WithHomeAway {
			item = append(item, []string{
				fmt.Sprintf("%d", d.HomeWon),
				fmt.Sprintf("%d", d.HomeLost),
				fmt.Sprintf("%.0f%%", d.HomeWonRate),
				fmt.Sprintf("%d", d.AwayWon),
				fmt.Sprintf("%d", d.HomeLost),
				fmt.Sprintf("%.0f%%", d.AwayWonRate),
			}...)
		}
		if p.option.WithTime {
			item = append(item, []string{
				fmt.Sprintf("%02d:%02d", d.TimePerGame/60, d.TimePerGame%60),
				fmt.Sprintf("%02d:%02d", d.LongestGameTime/60, d.LongestGameTime%60),
				fmt.Sprintf("%02d:%02d", d.ShortestGameTime/60, d.ShortestGameTime%60),
			}...)
		}
		if p.option.WithPoint {
			item = append(item, []string{
				fmt.Sprintf("%.2f", d.PointsPerGame),
				fmt.Sprintf("%.2f", d.PointsInPerGame),
				fmt.Sprintf("%.2f", d.DiffPerWon),
				fmt.Sprintf("%.2f", d.DiffPerLost),
			}...)
		}
		table = append(table, item)
	}
	return table
}

func (PlayerStats) playedTimeStats(data *model.EntityPlayer, timePlayed int) {
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

// Details .
func (p *PlayerStats) Details() []model.EntityPlayer {
	return p.players
}
