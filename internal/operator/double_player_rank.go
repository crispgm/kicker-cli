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
func (o DoublePlayerRank) SupportedFormats(trn *model.Tournament) bool {
	return openDoubleTournament(trn)
}

// Input .
func (o *DoublePlayerRank) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.players = players
	o.options = options
}

// Output .
func (o *DoublePlayerRank) Output() {
	data := make(map[string]entity.Player)
	for _, p := range o.players {
		data[p.Name] = p
	}
	for _, t := range o.tournaments {
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
				if g.GameType == entity.GameTypeQualification {
					t1p1Data.QualificationWin++
					t1p2Data.QualificationWin++
					t2p1Data.QualificationLoss++
					t2p2Data.QualificationLoss++
				} else if g.GameType == entity.GameTypeElimination {
					t1p1Data.EliminationWin++
					t1p2Data.EliminationWin++
					t2p1Data.EliminationLoss++
					t2p2Data.EliminationLoss++
				}
			} else if g.Point2 > g.Point1 {
				t1p1Data.Loss++
				t1p2Data.Loss++
				t2p1Data.Win++
				t2p2Data.Win++
				if g.GameType == entity.GameTypeQualification {
					t1p1Data.QualificationLoss++
					t1p2Data.QualificationLoss++
					t2p1Data.QualificationWin++
					t2p2Data.QualificationWin++
				} else if g.GameType == entity.GameTypeElimination {
					t1p1Data.EliminationLoss++
					t1p2Data.EliminationLoss++
					t2p1Data.EliminationWin++
					t2p2Data.EliminationWin++
				}
			} else {
				t1p1Data.Draw++
				t1p2Data.Draw++
				t2p1Data.Draw++
				t2p2Data.Draw++
				if g.GameType == entity.GameTypeElimination {
					t1p1Data.QualificationDraw++
					t1p2Data.QualificationDraw++
					t2p1Data.QualificationDraw++
					t2p2Data.QualificationDraw++
				}
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
			t1p1Data.EloRating = calculateELO(t1p1Data.GamesPlayed, t1p1Elo, team2elo, sa)
			t1p2Data.EloRating = calculateELO(t1p2Data.GamesPlayed, t1p2Elo, team2elo, sa)
			t2p1Data.EloRating = calculateELO(t2p1Data.GamesPlayed, t2p1Elo, team1elo, sb)
			t2p2Data.EloRating = calculateELO(t2p2Data.GamesPlayed, t2p2Elo, team1elo, sb)
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
			d.WinRate = float64(d.Win) / float64(d.GamesPlayed) * 100.0
			d.QualificationWinRate = float64(d.QualificationWin) / float64(d.QualificationWin+d.QualificationDraw+d.QualificationLoss) * 100.0
			d.EliminationWinRate = float64(d.EliminationWin) / float64(d.EliminationWin+d.EliminationLoss) * 100.0
			sliceData = append(sliceData, d)
		}
	}
	o.players = sliceData
	// }}}
	// {{{ sort
	sort.SliceStable(sliceData, func(i, j int) bool {
		if o.options.OrderBy == rating.RSysWinRate || o.options.OrderBy == rating.RSysELO {
			if sliceData[i].GamesPlayed >= o.options.MinimumPlayed && sliceData[j].GamesPlayed < o.options.MinimumPlayed {
				return true
			}
			if sliceData[i].GamesPlayed < o.options.MinimumPlayed && sliceData[j].GamesPlayed >= o.options.MinimumPlayed {
				return false
			}
		}

		if o.options.OrderBy == rating.RSysKicker {
			if sliceData[i].KickerPoints > sliceData[j].KickerPoints {
				return true
			}
		} else if o.options.OrderBy == rating.RSysATSA {
			if sliceData[i].ATSAPoints > sliceData[j].ATSAPoints {
				return true
			}
		} else if o.options.OrderBy == rating.RSysITSF {
			if sliceData[i].ITSFPoints > sliceData[j].ITSFPoints {
				return true
			}
		} else if o.options.OrderBy == rating.RSysELO {
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
	if o.options.Head > 0 && len(sliceData) > o.options.Head {
		sliceData = sliceData[:o.options.Head]
	} else if o.options.Tail > 0 && len(sliceData) > o.options.Tail {
		sliceData = sliceData[len(sliceData)-o.options.Tail:]
	}

	header := []string{"#", "Name", "Events", "Games", "Win", "Loss", "Draw", "WR%", "QWR%", "ELR%", "ELO", "KRP", "ATSA", "ITSF"}
	table := [][]string{}
	index := 1
	for _, d := range sliceData {
		if !o.options.ShowInactive && d.Inactive {
			continue
		}
		item := []string{
			fmt.Sprintf("%d", index),
			d.Name,
			fmt.Sprintf("%d", d.EventsPlayed),
			fmt.Sprintf("%d", d.GamesPlayed),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
			fmt.Sprintf("%.0f%%", d.WinRate),
			fmt.Sprintf("%.0f%%", d.QualificationWinRate),
			fmt.Sprintf("%.0f%%", d.EliminationWinRate),
			fmt.Sprintf("%.0f", d.EloRating),
			fmt.Sprintf("%d", d.KickerPoints),
			fmt.Sprintf("%d", d.ATSAPoints),
			fmt.Sprintf("%d", d.ITSFPoints),
		}
		table = append(table, item)
		index++
	}
	// }}}

	output(o.options, header, table)
}
