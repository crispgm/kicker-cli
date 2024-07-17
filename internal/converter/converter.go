// Package converter .
package converter

import (
	"fmt"
	"sort"
	"sync"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

// Converter .
type Converter struct {
	mu *sync.RWMutex

	record   entity.Record
	briefing string
	single   bool
}

// NewConverter .
func NewConverter() *Converter {
	return &Converter{
		mu:     &sync.RWMutex{},
		record: entity.Record{},
	}
}

// Normalize convert double games to entity formats
func (c *Converter) Normalize(orgPlayers []entity.Player, tournaments ...model.Tournament) (*entity.Record, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// players and teams
	for _, t := range tournaments {
		c.single = t.IsSingle()
		teams := make(map[string]model.Team)
		players := make(map[string]model.Player)
		idToOrgPlayers := make(map[string]entity.Player)
		for _, p := range t.Players {
			if !p.Removed {
				var found bool
				for _, op := range orgPlayers {
					if op.IsPlayer(p.Name) {
						found = true
						p.Name = op.Name
						players[p.ID] = p
						idToOrgPlayers[p.ID] = op
						break
					}
				}
				if !found {
					return nil, fmt.Errorf("%s not found", p.Name)
				}
			}
		}
		for _, t := range t.Teams {
			teams[t.ID] = t
		}
		for _, p := range idToOrgPlayers {
			c.record.Players = append(c.record.Players, p)
		}

		// convert rounds and knockoff games
		err := c.convertGames(t, teams, players)
		if err != nil {
			return nil, err
		}
		// calculate player rank based on results
		c.playerRank(t, teams, idToOrgPlayers)
	}

	// combine rounds/knockOffs
	c.combineAllGames()

	return &c.record, nil
}

func (c *Converter) convertGames(
	t model.Tournament,
	teams map[string]model.Team,
	players map[string]model.Player,
) error {
	rec := &c.record
	// preliminary rounds
	for _, r := range t.Rounds {
		games, err := c.convertPlayToGame(r.Name, r.Plays, teams, players)
		if err != nil {
			return err
		}
		rec.PreliminaryRounds = append(rec.PreliminaryRounds, games...)
	}

	for _, ko := range t.KnockOffs {
		for _, level := range ko.Levels {
			games, err := c.convertPlayToGame(level.Name, level.Plays, teams, players)
			if err != nil {
				return err
			}
			rec.WinnerBracket = append(rec.WinnerBracket, games...)
		}
		for _, level := range ko.LeftLevels {
			games, err := c.convertPlayToGame(level.Name, level.Plays, teams, players)
			if err != nil {
				return err
			}
			rec.LoserBracket = append(rec.LoserBracket, games...)
		}
		games, err := c.convertPlayToGame(ko.Third.Name, ko.Third.Plays, teams, players)
		if err != nil {
			return err
		}
		if len(games) > 0 {
			rec.ThirdPlace = &games[0]
		}
	}

	return nil
}

func (c *Converter) playerRank(
	t model.Tournament,
	teams map[string]model.Team,
	players map[string]entity.Player,
) {
	if len(t.KnockOffs) > 0 { // elimination
		ko := t.KnockOffs[0]
		if len(ko.LeftLevels) > 0 { // double elimination
			for i := len(ko.LeftLevels) - 1; i >= 0; i-- {
				level := ko.LeftLevels[i]
				var ranks []entity.Player
				for _, p := range level.Plays {
					rank := c.extractPlayerFromPlay(players, teams, p, rating.Loss) // extract losing players
					if len(rank) > 0 {
						ranks = append(ranks, rank...)
					}
				}
				if len(ranks) > 0 {
					c.record.Ranks = append(c.record.Ranks, ranks)
				}
			}
			// and final game participants to get champion and runner up
			{
				finalLevel := ko.Levels[len(ko.Levels)-1]
				finalPlay := finalLevel.Plays[0]
				runnerUp := c.extractPlayerFromPlay(players, teams, finalPlay, rating.Loss)
				champion := c.extractPlayerFromPlay(players, teams, finalPlay, rating.Win)
				c.record.Ranks = append(c.record.Ranks, runnerUp)
				c.record.Ranks = append(c.record.Ranks, champion)
			}
		} else if len(ko.Levels) > 0 { // single elimination
			// losing player of every levels
			for _, level := range ko.Levels {
				var ranks []entity.Player
				for _, p := range level.Plays {
					rank := c.extractPlayerFromPlay(players, teams, p, rating.Loss) // extract loss players
					if len(rank) > 0 {
						ranks = append(ranks, rank...)
					}
				}
				if len(ranks) > 0 {
					c.record.Ranks = append(c.record.Ranks, ranks)
				}
			}
			// third place winner if it's available
			if len(c.record.Players) > 2 { // prevent panic because there is *third* place even if there are only 2 players
				if len(ko.Third.Plays) > 0 {
					lastPos := len(c.record.Ranks) - 1
					runnerUp := c.record.Ranks[lastPos]
					third := c.extractPlayerFromPlay(players, teams, ko.Third.Plays[0], rating.Win)
					c.record.Ranks[lastPos] = third
					// remove third from fourth level
					fourthPos := lastPos - 1
					for i, p := range c.record.Ranks[fourthPos] {
						if len(third) > 0 && p.ID == third[0].ID {
							ranks := c.record.Ranks[fourthPos]
							ranks = append(ranks[:i], ranks[i+1:]...)
							c.record.Ranks[fourthPos] = ranks
							break
						}
					}
					c.record.Ranks = append(c.record.Ranks, runnerUp)
				}
			}
			// and final game winner
			{
				finalLevel := ko.Levels[len(ko.Levels)-1]
				finalPlay := finalLevel.Plays[0]
				champion := c.extractPlayerFromPlay(players, teams, finalPlay, rating.Win)
				c.record.Ranks = append(c.record.Ranks, champion)
			}
		}
	} else if len(t.Rounds) > 0 {
		// rounds only
		// cannot sort rounds only game atm
	} else {
		// no games
	}
}

func (c Converter) extractPlayerFromPlay(
	players map[string]entity.Player,
	teams map[string]model.Team,
	p model.Play,
	winDrawLoss int,
) []entity.Player {
	var ranks []entity.Player
	if !c.validPlay(p) {
		return ranks
	}
	if winDrawLoss == rating.Draw {
		// no idea right now
		return ranks
	}

	if winDrawLoss == rating.Loss {
		if p.Team1.Type == "Team" {
			team1 := teams[p.Team1.ID]
			team2 := teams[p.Team2.ID]
			if p.Winner == 1 {
				ranks = append(ranks, players[team2.Players[0].ID])
				ranks = append(ranks, players[team2.Players[1].ID])
			} else {
				ranks = append(ranks, players[team1.Players[0].ID])
				ranks = append(ranks, players[team1.Players[1].ID])
			}
		} else if p.Team1.Type == "Player" {
			if p.Winner == 1 {
				ranks = append(ranks, players[p.Team2.ID])
			} else {
				ranks = append(ranks, players[p.Team1.ID])
			}
		} else {
			return ranks
		}
	} else {
		if p.Team1.Type == "Team" {
			team1 := teams[p.Team1.ID]
			team2 := teams[p.Team2.ID]
			if p.Winner == 2 {
				ranks = append(ranks, players[team2.Players[0].ID])
				ranks = append(ranks, players[team2.Players[1].ID])
			} else {
				ranks = append(ranks, players[team1.Players[0].ID])
				ranks = append(ranks, players[team1.Players[1].ID])
			}
		} else if p.Team1.Type == "Player" {
			if p.Winner == 2 {
				ranks = append(ranks, players[p.Team2.ID])
			} else {
				ranks = append(ranks, players[p.Team1.ID])
			}
		} else {
			return ranks
		}

	}
	return ranks
}

func (c *Converter) combineAllGames() {
	rec := &c.record

	if len(rec.PreliminaryRounds) > 0 {
		rec.AllGames = append(rec.AllGames, rec.PreliminaryRounds...)
	}
	numOfWB, numOfLB := len(rec.WinnerBracket), len(rec.LoserBracket)
	if numOfWB > 0 && numOfLB == 0 {
		// single elimination
		rec.AllGames = append(rec.AllGames, rec.WinnerBracket...)
		// ThirdPlace is included in loser LoserBracket of double elimination
		// so only append in single elimination
		if rec.ThirdPlace != nil {
			rec.AllGames = append(rec.AllGames, *rec.ThirdPlace)
		}
	} else {
		rec.AllGames = append(rec.AllGames, rec.WinnerBracket...)
		rec.AllGames = append(rec.AllGames, rec.LoserBracket...)
		sort.SliceStable(
			rec.AllGames,
			func(i int, j int) bool {
				return rec.AllGames[i].TimeEnd < rec.AllGames[j].TimeEnd
			},
		)
	}
}

func (Converter) validPlay(p model.Play) bool {
	if !p.Valid || p.Deactivated || p.Skipped {
		return false
	}
	if p.Team1.ID == "" || p.Team2.ID == "" {
		return false
	}

	return true
}

func (c Converter) convertPlayToGame(
	name string,
	plays []model.Play,
	teams map[string]model.Team,
	players map[string]model.Player) ([]entity.Game, error) {
	var games []entity.Game
	for _, p := range plays {
		if !c.validPlay(p) {
			continue
		}
		var game entity.Game
		if p.Team1.Type == "Team" {
			team1 := teams[p.Team1.ID]
			team2 := teams[p.Team2.ID]
			t1p1 := players[team1.Players[0].ID]
			t1p2 := players[team1.Players[1].ID]
			t2p1 := players[team2.Players[0].ID]
			t2p2 := players[team2.Players[1].ID]
			game.Team1 = []string{t1p1.Name, t1p2.Name}
			game.Team2 = []string{t2p1.Name, t2p2.Name}
		} else if p.Team1.Type == "Player" {
			game.Team1 = []string{players[p.Team1.ID].Name}
			game.Team2 = []string{players[p.Team2.ID].Name}
		} else {
			continue
		}
		game.TimeStart = p.TimeStart
		game.TimeEnd = p.TimeEnd
		game.TimePlayed = (p.TimeEnd - p.TimeStart) / 1000
		game.Winner = p.Winner
		// team result may not exist
		// game.Point1 = p.Team1Result
		// game.Point2 = p.Team2Result
		var sets []entity.Set
		for _, d := range p.Disciplines {
			for _, s := range d.Sets {
				sets = append(sets, entity.Set{Point1: s.Team1, Point2: s.Team2})
				if s.Team1 > s.Team2 {
					game.Point1++
				} else if s.Team1 < s.Team2 {
					game.Point2++
				} else {
					// draw
					game.Point1++
					game.Point2++
				}
			}
			break // only support one discipline right now
		}
		game.Name = name
		game.Sets = sets
		games = append(games, game)
	}
	return games, nil
}

// Briefing .
func (c *Converter) Briefing() string {
	players := make(map[string]bool)
	for _, g := range c.record.PreliminaryRounds {
		players[g.Team1[0]] = true
		players[g.Team2[0]] = true
		if !c.single {
			players[g.Team2[1]] = true
			players[g.Team1[1]] = true
		}
	}
	for _, g := range c.record.WinnerBracket {
		players[g.Team1[0]] = true
		players[g.Team2[0]] = true
		if !c.single {
			players[g.Team2[1]] = true
			players[g.Team1[1]] = true
		}
	}
	for _, g := range c.record.LoserBracket {
		players[g.Team1[0]] = true
		players[g.Team2[0]] = true
		if !c.single {
			players[g.Team2[1]] = true
			players[g.Team1[1]] = true
		}
	}
	numOfGames := len(c.record.PreliminaryRounds) +
		len(c.record.WinnerBracket) +
		len(c.record.LoserBracket)
	if c.record.ThirdPlace != nil {
		numOfGames++
	}
	c.briefing = fmt.Sprintf("%d player(s) played %d game(s)", len(players), numOfGames)
	return c.briefing
}
