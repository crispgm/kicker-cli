// Package converter .
package converter

import (
	"fmt"
	"sync"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

// Converter .
type Converter struct {
	mu *sync.RWMutex

	record   entity.Record
	briefing string
}

// NewConverter .
func NewConverter() *Converter {
	return &Converter{
		mu:     &sync.RWMutex{},
		record: entity.Record{},
	}
}

// Normalize convert double games to entity formats
func (c *Converter) Normalize(tournaments []model.Tournament, ePlayers []entity.Player) (*entity.Record, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	rec := &entity.Record{}

	// players and teams
	for _, t := range tournaments {
		teams := make(map[string]model.Team)
		players := make(map[string]model.Player)
		for _, p := range t.Players {
			if !p.Removed {
				var found bool
				for _, ep := range ePlayers {
					if ep.IsPlayer(p.Name) {
						found = true
						p.Name = ep.Name
						players[p.ID] = p
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

		// preliminary rounds
		for _, r := range t.Rounds {
			games, err := c.convertPlayToGame(r.Plays, teams, players)
			if err != nil {
				return nil, err
			}
			rec.PreliminaryRounds = append(rec.PreliminaryRounds, games...)
		}

		for _, ko := range t.KnockOffs {
			for _, level := range ko.Levels {
				games, err := c.convertPlayToGame(level.Plays, teams, players)
				if err != nil {
					return nil, err
				}
				rec.WinnerBracket = append(rec.WinnerBracket, games...)
			}
			for _, level := range ko.LeftLevels {
				games, err := c.convertPlayToGame(level.Plays, teams, players)
				if err != nil {
					return nil, err
				}
				rec.LoserBracket = append(rec.LoserBracket, games...)
			}
			games, err := c.convertPlayToGame(ko.Third.Plays, teams, players)
			if err != nil {
				return nil, err
			}
			if len(games) > 0 {
				rec.ThirdPlace = &games[0]
			}
		}
	}

	if len(rec.PreliminaryRounds) > 0 {
		rec.AllGames = append(rec.AllGames, rec.PreliminaryRounds...)
	}
	if len(rec.WinnerBracket) > 0 {
		rec.AllGames = append(rec.AllGames, rec.WinnerBracket...)
	}
	if len(rec.LoserBracket) > 0 {
		rec.AllGames = append(rec.AllGames, rec.LoserBracket...)
	}
	if rec.ThirdPlace != nil {
		rec.AllGames = append(rec.AllGames, *rec.ThirdPlace)
	}
	c.record = *rec
	return rec, nil
}

func (Converter) convertPlayToGame(
	plays []model.Play,
	teams map[string]model.Team,
	players map[string]model.Player) ([]entity.Game, error) {
	var games []entity.Game
	for _, p := range plays {
		if !p.Valid || p.Deactivated || p.Skipped {
			continue
		}

		team1 := teams[p.Team1.ID]
		team2 := teams[p.Team2.ID]
		if p.Team2.ID == "" {
			// pass the game
			continue
		}
		t1p1 := players[team1.Players[0].ID]
		t1p2 := players[team1.Players[1].ID]
		t2p1 := players[team2.Players[0].ID]
		t2p2 := players[team2.Players[1].ID]
		var game entity.Game
		game.Team1 = []string{t1p1.Name, t1p2.Name}
		game.Team2 = []string{t2p1.Name, t2p2.Name}
		game.TimePlayed = (p.TimeEnd - p.TimeStart) / 1000
		for _, d := range p.Disciplines {
			for _, s := range d.Sets {
				game.Point1 = s.Team1
				game.Point2 = s.Team2
			}
		}
		games = append(games, game)
	}
	return games, nil
}

// Briefing .
func (c *Converter) Briefing() string {
	players := make(map[string]bool)
	for _, g := range c.record.PreliminaryRounds {
		players[g.Team1[0]] = true
		players[g.Team1[1]] = true
		players[g.Team2[0]] = true
		players[g.Team2[1]] = true
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
