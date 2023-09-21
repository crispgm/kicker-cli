// Package double is operators for double games
package double

import (
	"fmt"
	"sort"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

var _ operator.Operator = (*PlayerRank)(nil)

// PlayerRank generate statistics data of double tournaments by player
type PlayerRank struct {
	options     operator.Option
	tournaments []entity.Tournament
	players     []entity.Player
}

// SupportedFormats .
func (p PlayerRank) SupportedFormats(trn *model.Tournament) bool {
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
func (p *PlayerRank) Input(tournaments []entity.Tournament, players []entity.Player, options operator.Option) {
	p.tournaments = tournaments
	p.players = players
	p.options = options
}

// Output .
func (p *PlayerRank) Output() [][]string {
	data := make(map[string]entity.Player)
	for _, p := range p.players {
		data[p.Name] = p
	}
	for _, t := range p.tournaments {
		for _, g := range t.Converted.AllGames {
			t1p1Data := data[g.Team1[0]]
			t1p2Data := data[g.Team1[1]]
			t2p1Data := data[g.Team2[0]]
			t2p2Data := data[g.Team2[1]]
			t1p1Data.Name = g.Team1[0]
			t1p2Data.Name = g.Team1[1]
			t2p1Data.Name = g.Team2[0]
			t2p2Data.Name = g.Team2[1]
			t1p1Data.GamesPlayed++
			t1p2Data.GamesPlayed++
			t2p1Data.GamesPlayed++
			t2p2Data.GamesPlayed++
			t1p1Data.TimePlayed += g.TimePlayed
			t1p2Data.TimePlayed += g.TimePlayed
			t2p1Data.TimePlayed += g.TimePlayed
			t2p2Data.TimePlayed += g.TimePlayed
			if g.Point1 > g.Point2 {
				t1p1Data.Win++
				t1p2Data.Win++
				t2p1Data.Loss++
				t2p2Data.Loss++
				t1p1Data.HomeWin++
				t1p2Data.HomeWin++
				t2p1Data.AwayLoss++
				t2p2Data.AwayLoss++
				t1p1Data.GoalsWin += (g.Point1 - g.Point2)
				t1p2Data.GoalsWin += (g.Point1 - g.Point2)
				t2p1Data.GoalsInLoss += (g.Point1 - g.Point2)
				t2p2Data.GoalsInLoss += (g.Point1 - g.Point2)
			} else if g.Point2 > g.Point1 {
				t1p1Data.Loss++
				t1p2Data.Loss++
				t2p1Data.Win++
				t2p2Data.Win++
				t1p1Data.HomeLoss++
				t1p2Data.HomeLoss++
				t2p1Data.AwayWin++
				t2p2Data.AwayWin++
				t2p1Data.GoalsWin += (g.Point2 - g.Point1)
				t2p2Data.GoalsWin += (g.Point2 - g.Point1)
				t1p1Data.GoalsInLoss += (g.Point2 - g.Point1)
				t1p2Data.GoalsInLoss += (g.Point2 - g.Point1)
			} else {
				// basically not approachable
				t1p1Data.Draw++
				t1p2Data.Draw++
				t2p1Data.Draw++
				t2p2Data.Draw++
			}
			t1p1Data.Goals += g.Point1
			t1p2Data.Goals += g.Point1
			t2p1Data.Goals += g.Point2
			t2p2Data.Goals += g.Point2
			t1p1Data.GoalsIn += g.Point2
			t1p2Data.GoalsIn += g.Point2
			t2p1Data.GoalsIn += g.Point1
			t2p2Data.GoalsIn += g.Point1
			// ELO
			elo := rating.Elo{}
			t1p1Elo := elo.InitialScore()
			t1p2Elo := elo.InitialScore()
			t2p1Elo := elo.InitialScore()
			t2p2Elo := elo.InitialScore()
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
			sa := rating.Win
			sb := rating.Loss
			if g.Point1 == g.Point2 {
				sa = rating.Draw
				sb = rating.Draw
			} else if g.Point1 < g.Point2 {
				sa = rating.Loss
				sb = rating.Win
			}
			team1elo := (t1p1Elo + t1p2Elo) / 2
			team2elo := (t2p1Elo + t2p2Elo) / 2
			t1p1Data.EloRating = p.calculateELO(t1p1Data.GamesPlayed, t1p1Elo, team2elo, sa)
			t1p2Data.EloRating = p.calculateELO(t1p2Data.GamesPlayed, t1p2Elo, team2elo, sa)
			t2p1Data.EloRating = p.calculateELO(t2p1Data.GamesPlayed, t2p1Elo, team1elo, sb)
			t2p2Data.EloRating = p.calculateELO(t2p2Data.GamesPlayed, t2p2Elo, team1elo, sb)

			data[g.Team1[0]] = t1p1Data
			data[g.Team1[1]] = t1p2Data
			data[g.Team2[0]] = t2p1Data
			data[g.Team2[1]] = t2p2Data
		}
		// ranking points
		curRank := 0
		for i := len(t.Converted.Ranks) - 1; i >= 0; i-- {
			rank := t.Converted.Ranks[i]
			curRank += len(rank) / 2
			factors := rating.Factor{
				Place: curRank,
			}
			for _, r := range rank {
				ranker := rating.Rank{}
				d := data[r.Name]
				if len(t.Event.KickerLevel) > 0 {
					factors.PlayerScore = float64(d.KickerPoints)
					factors.Level = t.Event.KickerLevel
					d.KickerPoints = int(ranker.Calculate(factors))
				}
				if len(t.Event.ATSALevel) > 0 {
					factors.Level = t.Event.ATSALevel
					factors.PlayerScore = float64(d.ATSAPoints)
					d.ATSAPoints = int(ranker.Calculate(factors))
				}
				if len(t.Event.ITSFLevel) > 0 {
					factors.PlayerScore = float64(d.ITSFPoints)
					factors.Level = t.Event.ITSFLevel
					d.ITSFPoints = int(ranker.Calculate(factors))
				}
				d.EventsPlayed++
				data[r.Name] = d
			}
		}
	}

	var sliceData []entity.Player
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.GamesPlayed != 0 {
			d.WinRate = float32(d.Win) / float32(d.GamesPlayed) * 100.0
			if d.HomeWin+d.HomeLoss > 0 {
				d.HomeWinRate = float32(d.HomeWin) / float32(d.HomeWin+d.HomeLoss) * 100.0
			}
			if d.AwayWin+d.AwayLoss > 0 {
				d.AwayWinRate = float32(d.AwayWin) / float32(d.AwayWin+d.AwayLoss) * 100.0
			}
			d.PointsPerGame = float32(d.Goals) / float32(d.GamesPlayed)
			d.PointsInPerGame = float32(d.GoalsIn) / float32(d.GamesPlayed)
			d.TimePerGame = d.TimePlayed / d.GamesPlayed / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			if d.Win > 0 {
				d.DiffPerWin = float32(d.GoalsWin) / float32(d.Win)
			}
			if d.Loss > 0 {
				d.DiffPerLoss = float32(d.GoalsInLoss) / float32(d.Loss)
			}
			sliceData = append(sliceData, d)
		}
	}
	p.players = sliceData
	sort.SliceStable(sliceData, func(i, j int) bool {
		if p.options.OrderBy == "wr" || p.options.OrderBy == "elo" {
			if sliceData[i].GamesPlayed >= p.options.MinimumPlayed && sliceData[j].GamesPlayed < p.options.MinimumPlayed {
				return true
			}
			if sliceData[i].GamesPlayed < p.options.MinimumPlayed && sliceData[j].GamesPlayed >= p.options.MinimumPlayed {
				return false
			}
		}

		if p.options.OrderBy == "krs" {
			if sliceData[i].KickerPoints > sliceData[j].KickerPoints {
				return true
			}
		} else if p.options.OrderBy == "atsa" {
			if sliceData[i].ATSAPoints > sliceData[j].ATSAPoints {
				return true
			}
		} else if p.options.OrderBy == "itsf" {
			if sliceData[i].ITSFPoints > sliceData[j].ITSFPoints {
				return true
			}
		} else if p.options.OrderBy == "elo" {
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

	if p.options.Head > 0 && len(sliceData) > p.options.Head {
		sliceData = sliceData[:p.options.Head]
	} else if p.options.Tail > 0 && len(sliceData) > p.options.Tail {
		sliceData = sliceData[len(sliceData)-p.options.Tail:]
	}

	header := []string{"#", "Name", "Events", "Games", "Win", "Loss", "Draw", "WR%", "Elo", "KRS", "ATSA", "ITSF"}
	pointHeader := []string{"G+", "G-", "G±", "PPG", "LPG", "DPW", "DPL"}
	if p.options.WithGoals {
		header = append(header, pointHeader...)
	}
	table := [][]string{}
	if p.options.WithHeader {
		table = append(table, header)
	}
	for i, d := range sliceData {
		item := []string{
			fmt.Sprintf("%d", i+1),
			d.Name,
			fmt.Sprintf("%d", d.EventsPlayed),
			fmt.Sprintf("%d", d.GamesPlayed),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
			fmt.Sprintf("%.0f%%", d.WinRate),
			fmt.Sprintf("%.0f", d.EloRating),
			fmt.Sprintf("%d", d.KickerPoints),
			fmt.Sprintf("%d", d.ATSAPoints),
			fmt.Sprintf("%d", d.ITSFPoints),
		}
		if p.options.WithGoals {
			item = append(item, []string{
				fmt.Sprintf("%d", d.Goals),
				fmt.Sprintf("%d", d.GoalsIn),
				fmt.Sprintf("%d", d.GoalDiff),
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

// calculateELO calculate ELO for player
func (p PlayerRank) calculateELO(played int, p1Elo, p2Elo float64, result int) float64 {
	eloCalc := rating.Elo{}
	factors := rating.Factor{
		Played:        played,
		PlayerScore:   p1Elo,
		OpponentScore: p2Elo,
		Result:        result,
	}
	return eloCalc.Calculate(factors)
}
