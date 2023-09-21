// Package single is operators for single games
package single

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
	if trn.IsSingle() {
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
			p1Data := data[g.Team1[0]]
			p2Data := data[g.Team2[0]]
			p1Data.Name = g.Team1[0]
			p2Data.Name = g.Team2[0]
			p1Data.GamesPlayed++
			p2Data.GamesPlayed++
			p1Data.TimePlayed += g.TimePlayed
			p2Data.TimePlayed += g.TimePlayed
			if g.Point1 > g.Point2 {
				p1Data.Win++
				p2Data.Loss++
				p1Data.HomeWin++
				p2Data.AwayLoss++
				p1Data.GoalsWin += (g.Point1 - g.Point2)
				p2Data.GoalsInLoss += (g.Point1 - g.Point2)
			} else if g.Point2 > g.Point1 {
				p1Data.Loss++
				p2Data.Win++
				p1Data.HomeLoss++
				p2Data.AwayWin++
				p2Data.GoalsWin += (g.Point2 - g.Point1)
				p1Data.GoalsInLoss += (g.Point2 - g.Point1)
			} else {
				// basically not approachable
				p1Data.Draw++
				p2Data.Draw++
			}
			p1Data.Goals += g.Point1
			p2Data.Goals += g.Point2
			p1Data.GoalsIn += g.Point2
			p2Data.GoalsIn += g.Point1
			elo := rating.Elo{}
			p1Elo := elo.InitialScore()
			p2Elo := elo.InitialScore()
			if p1Data.EloRating != 0 {
				p1Elo = p1Data.EloRating
			}
			if p2Data.EloRating != 0 {
				p2Elo = p2Data.EloRating
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
			p1Data.EloRating = p.calculateELO(p1Data.GamesPlayed, p1Elo, p2Elo, sa)
			p2Data.EloRating = p.calculateELO(p2Data.GamesPlayed, p2Elo, p1Elo, sb)
			data[g.Team1[0]] = p1Data
			data[g.Team2[0]] = p2Data
		}
		// ranking points
		curRank := 0
		for i := len(t.Converted.Ranks) - 1; i >= 0; i-- {
			rank := t.Converted.Ranks[i]
			curRank += len(rank)
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
		if d.GamesPlayed != 0 {
			d.GoalDiff = d.Goals - d.GoalsIn
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
