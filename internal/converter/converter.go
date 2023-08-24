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
	eGames   []entity.Game
	briefing string

	mu *sync.RWMutex
}

// NewConverter .
func NewConverter() *Converter {
	return &Converter{
		mu: &sync.RWMutex{},
	}
}

// Normalize .
func (c *Converter) Normalize(tournaments []model.Tournament, ePlayers []entity.Player) ([]entity.Game, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

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

		for _, r := range t.Rounds {
			for _, p := range r.Plays {
				if !p.Valid || p.Deactivated || p.Skipped {
					continue
				}
				team1 := teams[p.Team1.ID]
				team2 := teams[p.Team2.ID]
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

				c.eGames = append(c.eGames, game)
			}
		}
	}
	return c.eGames, nil
}

// Briefing .
func (c *Converter) Briefing() string {
	players := make(map[string]bool)
	for _, g := range c.eGames {
		players[g.Team1[0]] = true
		players[g.Team1[1]] = true
		players[g.Team2[0]] = true
		players[g.Team2[1]] = true
	}
	c.briefing = fmt.Sprintf("%d player(s) played %d game(s)", len(players), len(c.eGames))
	return c.briefing
}
