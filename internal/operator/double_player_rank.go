package operator

import (
	"fmt"
	"sort"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

var _ Operator = (*DoublePlayerRank)(nil)

// DoublePlayerRank generate statistics data of double tournaments by player
type DoublePlayerRank struct {
	options     Option
	tournaments []entity.Tournament
	players     []entity.Player
}

// SupportedFormats .
func (p DoublePlayerRank) SupportedFormats(trn *model.Tournament) bool {
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
func (p *DoublePlayerRank) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	p.tournaments = tournaments
	p.players = players
	p.options = options
}

// Output .
func (p *DoublePlayerRank) Output() [][]string {
	data := make(map[string]entity.Player)
	for _, p := range p.players {
		data[p.Name] = p
	}
	for _, t := range p.tournaments {
		var played = make(map[string]bool)
		for _, g := range t.Converted.AllGames {
			t1p1Data := data[g.Team1[0]]
			t1p2Data := data[g.Team1[1]]
			t2p1Data := data[g.Team2[0]]
			t2p2Data := data[g.Team2[1]]
			t1p1Data.Name = g.Team1[0]
			t1p2Data.Name = g.Team1[1]
			t2p1Data.Name = g.Team2[0]
			t2p2Data.Name = g.Team2[1]

			// {{{ game data
			t1p1Data.GamesPlayed++
			t1p2Data.GamesPlayed++
			t2p1Data.GamesPlayed++
			t2p2Data.GamesPlayed++
			if g.Point1 > g.Point2 {
				t1p1Data.Win++
				t1p2Data.Win++
				t2p1Data.Loss++
				t2p2Data.Loss++
				t1p1Data.HomeWin++
				t1p2Data.HomeWin++
				t2p1Data.AwayLoss++
				t2p2Data.AwayLoss++
			} else if g.Point2 > g.Point1 {
				t1p1Data.Loss++
				t1p2Data.Loss++
				t2p1Data.Win++
				t2p2Data.Win++
				t1p1Data.HomeLoss++
				t1p2Data.HomeLoss++
				t2p1Data.AwayWin++
				t2p2Data.AwayWin++
			} else {
				// basically not approachable
				t1p1Data.Draw++
				t1p2Data.Draw++
				t2p1Data.Draw++
				t2p2Data.Draw++
			}
			// }}}
			// {{{ ELO
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
			// }}}
			// {{{ mark tournament played
			if _, ok := played[t1p1Data.Name]; !ok {
				t1p1Data.EventsPlayed++
				played[t1p1Data.Name] = true
			}
			if _, ok := played[t1p2Data.Name]; !ok {
				t1p2Data.EventsPlayed++
				played[t1p2Data.Name] = true
			}
			if _, ok := played[t2p1Data.Name]; !ok {
				t2p1Data.EventsPlayed++
				played[t2p1Data.Name] = true
			}
			if _, ok := played[t2p2Data.Name]; !ok {
				t2p2Data.EventsPlayed++
				played[t2p2Data.Name] = true
			}
			// }}}

			data[g.Team1[0]] = t1p1Data
			data[g.Team1[1]] = t1p2Data
			data[g.Team2[0]] = t2p1Data
			data[g.Team2[1]] = t2p2Data
		}
		// {{{ ranking points
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
				data[r.Name] = d
			}
		}
		// }}}
	}

	// {{{ map to slice
	var sliceData []entity.Player
	for _, d := range data {
		if d.GamesPlayed != 0 {
			d.WinRate = float32(d.Win) / float32(d.GamesPlayed) * 100.0
			if d.HomeWin+d.HomeLoss > 0 {
				d.HomeWinRate = float32(d.HomeWin) / float32(d.HomeWin+d.HomeLoss) * 100.0
			}
			if d.AwayWin+d.AwayLoss > 0 {
				d.AwayWinRate = float32(d.AwayWin) / float32(d.AwayWin+d.AwayLoss) * 100.0
			}
			sliceData = append(sliceData, d)
		}
	}
	p.players = sliceData
	// }}}
	// {{{ sort
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
	// }}}

	// {{{ build result
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
	// }}}
}

// calculateELO calculate ELO for player
func (p DoublePlayerRank) calculateELO(played int, p1Elo, p2Elo float64, result int) float64 {
	eloCalc := rating.Elo{}
	factors := rating.Factor{
		Played:        played,
		PlayerScore:   p1Elo,
		OpponentScore: p2Elo,
		Result:        result,
	}
	return eloCalc.Calculate(factors)
}
