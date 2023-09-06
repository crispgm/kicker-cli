package entity

import (
	"strings"

	"github.com/crispgm/kicker-cli/internal/util"
)

// Player is a real world player
type Player struct {
	ID        string   `yaml:"id"`
	ITSFID    string   `yaml:"itsf_id"`
	ATSAID    string   `yaml:"atsa_id"`
	Name      string   `yaml:"name"`
	Aliases   []string `yaml:"aliases"`
	Points    int      `yaml:"points"`
	Played    int      `yaml:"played"`
	Won       int      `yaml:"won"`
	Lost      int      `yaml:"lost"`
	Draws     int      `yaml:"draws"`
	EloRating float64  `yaml:"elo_rating"`
	WinRate   float32  `yaml:"win_rate"`

	TimePlayed       int     `yaml:"time_played,omitempty"`
	LongestGameTime  int     `yaml:"longest_game_time,omitempty"`
	ShortestGameTime int     `yaml:"shortest_game_time,omitempty"`
	TimePerGame      int     `yaml:"time_per_game,omitempty"`
	HomeWon          int     `yaml:"home_won,omitempty,omitempty"`
	HomeLost         int     `yaml:"home_lost,omitempty"`
	HomeWonRate      float32 `yaml:"home_won_rate,omitempty"`
	AwayWon          int     `yaml:"away_won,omitempty"`
	AwayLost         int     `yaml:"away_lost,omitempty"`
	AwayWonRate      float32 `yaml:"away_won_rate,omitempty"`
	Goals            int     `yaml:"goals,omitempty"`
	GoalsIn          int     `yaml:"goals_in,omitempty"`
	GoalDiff         int     `yaml:"goal_diff,omitempty"`
	PointsPerGame    float32 `yaml:"points_per_game,omitempty"`
	PointsInPerGame  float32 `yaml:"points_in_per_game,omitempty"`
	GoalsWon         int     `yaml:"goals_won,omitempty"`
	DiffPerWon       float32 `yaml:"diff_per_won,omitempty"`
	GoalsInLost      int     `yaml:"goals_in_lost,omitempty"`
	DiffPerLost      float32 `yaml:"diff_per_lost,omitempty"`
}

// NewPlayer creates a player
func NewPlayer(name string) *Player {
	return &Player{
		ID:   util.UUID(),
		Name: name,
	}
}

// IsPlayer .
func (p Player) IsPlayer(name string) bool {
	name = strings.ToLower(strings.Trim(name, " "))
	if name == strings.ToLower(p.Name) {
		return true
	}
	for _, alias := range p.Aliases {
		if strings.ToLower(alias) == name {
			return true
		}
	}

	return false
}

// AddAlias add an alias for player
func (p *Player) AddAlias(aliases ...string) {
	p.Aliases = append(p.Aliases, aliases...)
}
