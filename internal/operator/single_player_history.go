package operator

import (
	"fmt"

	"github.com/guptarohit/asciigraph"
	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

var _ Operator = (*SinglePlayerHistory)(nil)

// SinglePlayerHistory generate data details of single tournaments by player
type SinglePlayerHistory struct {
	options     Option
	tournaments []entity.Tournament
	players     []entity.Player
}

// SupportedFormats .
func (o SinglePlayerHistory) SupportedFormats(trn *model.Tournament) bool {
	return openSingleTournament(trn)
}

// Input .
func (o *SinglePlayerHistory) Input(tournaments []entity.Tournament, players []entity.Player, options Option) {
	o.tournaments = tournaments
	o.players = players
	o.options = options
}

// Output .
func (o *SinglePlayerHistory) Output() {
	found := false
	header := []string{"#", "Event", "Player 1", "Player 2", "Result", "WR%", "ELO", "KRP", "ATSA", "ITSF"}
	table := [][]string{}
	eloChart := []float64{}
	winRateChart := []float64{}
	data := make(map[string]entity.Player)
	for _, p := range o.players {
		data[p.Name] = p
	}
	seq := 1
	firstGamePlayed := false
	for _, t := range o.tournaments {
		var played = make(map[string]bool)
		for _, g := range t.Converted.AllGames {
			p1Data := data[g.Team1[0]]
			p2Data := data[g.Team2[0]]
			p1Data.Name = g.Team1[0]
			p2Data.Name = g.Team2[0]

			playerBefore := o.choosePlayerData(p1Data, p2Data)
			if playerBefore != nil {
				playerBefore.WinRate = float64(playerBefore.Win) / float64(playerBefore.GamesPlayed)
			}

			// {{{ game data
			p1Data.GamesPlayed++
			p2Data.GamesPlayed++
			if g.Point1 > g.Point2 {
				p1Data.Win++
				p2Data.Loss++
			} else if g.Point2 > g.Point1 {
				p1Data.Loss++
				p2Data.Win++
			} else {
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

			if playerBefore != nil {
				if !firstGamePlayed && playerBefore.EloRating == 0.0 {
					elo := rating.Elo{}
					playerBefore.EloRating = elo.InitialScore()
				}
				player := o.choosePlayerData(p1Data, p2Data)
				player.WinRate = float64(player.Win) / float64(player.GamesPlayed)
				pointText := fmt.Sprintf("%d:%d", g.Point1, g.Point2)
				winRateText := fmt.Sprintf("%.2f%% -> %.2f%%", playerBefore.WinRate*100, player.WinRate*100)
				eloText := fmt.Sprintf("%.0f -> %.0f", playerBefore.EloRating, player.EloRating)
				if player.EloRating < playerBefore.EloRating {
					eloText = pterm.FgRed.Sprintf("%.0f -> %.0f (%.0f)", playerBefore.EloRating, player.EloRating, player.EloRating-playerBefore.EloRating)
				} else if player.EloRating > playerBefore.EloRating {
					eloText = pterm.FgGreen.Sprintf("%.0f -> %.0f (+%.0f)", playerBefore.EloRating, player.EloRating, player.EloRating-playerBefore.EloRating)
				}
				if !firstGamePlayed || player.WinRate == playerBefore.WinRate {
					winRateText = pterm.Sprintf("%.2f%%", player.WinRate*100)
				} else {
					if player.WinRate < playerBefore.WinRate {
						winRateText = pterm.FgRed.Sprintf("%.2f%% -> %.2f%% (%.2f%%)", playerBefore.WinRate*100, player.WinRate*100, (player.WinRate-playerBefore.WinRate)*100)
					} else if player.WinRate > playerBefore.WinRate {
						winRateText = pterm.FgGreen.Sprintf("%.2f%% -> %.2f%% (+%.2f%%)", playerBefore.WinRate*100, player.WinRate*100, (player.WinRate-playerBefore.WinRate)*100)
					}
				}
				table = append(table, []string{
					fmt.Sprintf("%d", seq),
					t.Event.Name,
					p1Data.Name,
					p2Data.Name,
					pointText,
					winRateText,
					eloText,
					fmt.Sprintf("%0d", player.KickerPoints),
					fmt.Sprintf("%0d", player.ATSAPoints),
					fmt.Sprintf("%0d", player.ITSFPoints),
				})
				eloChart = append(eloChart, player.EloRating)
				winRateChart = append(winRateChart, player.WinRate*100)
				found = true
				firstGamePlayed = true
				seq++
			}
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

	if found {
		output(o.options, header, table)

		if o.options.OutputFormat == "default" {
			if len(eloChart) > 0 {
				fmt.Println()
				eloGraph := asciigraph.Plot(eloChart, asciigraph.Caption("ELO"), asciigraph.Height(20), asciigraph.Width(80))
				fmt.Println(eloGraph)
			}
			if len(winRateChart) > 0 {
				fmt.Println()
				winRateGraph := asciigraph.Plot(winRateChart, asciigraph.Caption("Win Rate"), asciigraph.Height(20), asciigraph.Width(80))
				fmt.Println(winRateGraph)
			}
		}
	}
}

func (o SinglePlayerHistory) choosePlayerData(p1Data, p2Data entity.Player) *entity.Player {
	name := o.options.PlayerName
	if p1Data.Name == name {
		return &p1Data
	} else if p2Data.Name == name {
		return &p2Data
	}
	return nil
}
