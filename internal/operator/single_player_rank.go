package operator

import (
	"fmt"
	"sort"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

var _ Operator = (*SinglePlayerRank)(nil)

// SinglePlayerRank generate statistics data of double tournaments by player
type SinglePlayerRank struct {
	options     Option
	tournaments []entity.Tournament
	players     []entity.Player
}

// SupportedFormats .
func (o SinglePlayerRank) SupportedFormats(trn *model.Tournament) bool {
	return openSingleTournament(trn)
}

// Input .
func (o *SinglePlayerRank) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.players = players
	o.options = options
}

// Output .
func (o *SinglePlayerRank) Output() {
	data := make(map[string]entity.Player)
	for _, p := range o.players {
		data[p.Name] = p
	}
	for _, t := range o.tournaments {
		var played = make(map[string]bool)
		for _, g := range t.Converted.AllGames {
			p1Data := data[g.Team1[0]]
			p2Data := data[g.Team2[0]]
			p1Data.Name = g.Team1[0]
			p2Data.Name = g.Team2[0]

			// {{{ game data
			p1Data.GamesPlayed++
			p2Data.GamesPlayed++
			if g.Point1 > g.Point2 {
				p1Data.Win++
				p2Data.Loss++
				p1Data.HomeWin++
				p2Data.AwayLoss++
			} else if g.Point2 > g.Point1 {
				p1Data.Loss++
				p2Data.Win++
				p1Data.HomeLoss++
				p2Data.AwayWin++
			} else {
				// basically not approachable
				p1Data.Draw++
				p2Data.Draw++
			}
			// }}}
			// {{{ ELO
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
			p1Data.EloRating = calculateELO(p1Data.GamesPlayed, p1Elo, p2Elo, sa)
			p2Data.EloRating = calculateELO(p2Data.GamesPlayed, p2Elo, p1Elo, sb)
			// }}}
			// {{{ mark tournament played
			if _, ok := played[p1Data.Name]; !ok {
				p1Data.EventsPlayed++
				played[p1Data.Name] = true
			}
			if _, ok := played[p2Data.Name]; !ok {
				p2Data.EventsPlayed++
				played[p2Data.Name] = true
			}
			// }}}

			data[g.Team1[0]] = p1Data
			data[g.Team2[0]] = p2Data
		}
		// {{{ ranking points
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
			if d.HomeWin+d.HomeLoss > 0 {
				d.HomeWinRate = float64(d.HomeWin) / float64(d.HomeWin+d.HomeLoss) * 100.0
			}
			if d.AwayWin+d.AwayLoss > 0 {
				d.AwayWinRate = float64(d.AwayWin) / float64(d.AwayWin+d.AwayLoss) * 100.0
			}
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

	header := []string{"#", "Name", "Events", "Games", "Win", "Loss", "Draw", "WR%", "ELO", "KRP", "ATSA", "ITSF"}
	table := [][]string{}
	if o.options.WithHeader {
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
	pterm.DefaultTable.WithHasHeader(o.options.WithHeader).WithData(table).WithBoxed(o.options.WithBoxes).Render()
	// }}}
}
