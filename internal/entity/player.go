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
	Win       int      `yaml:"win"`
	Loss      int      `yaml:"loss"`
	Draw      int      `yaml:"draw"`
	EloRating float64  `yaml:"elo_rating"`
	WinRate   float32  `yaml:"win_rate"`

	TimePlayed       int     `yaml:"time_played,omitempty"`
	LongestGameTime  int     `yaml:"longest_game_time,omitempty"`
	ShortestGameTime int     `yaml:"shortest_game_time,omitempty"`
	TimePerGame      int     `yaml:"time_per_game,omitempty"`
	HomeWin          int     `yaml:"home_win,omitempty,omitempty"`
	HomeLoss         int     `yaml:"home_loss,omitempty"`
	HomeWinRate      float32 `yaml:"home_win_rate,omitempty"`
	AwayWin          int     `yaml:"away_win,omitempty"`
	AwayLoss         int     `yaml:"away_loss,omitempty"`
	AwayWinRate      float32 `yaml:"away_win_rate,omitempty"`
	Goals            int     `yaml:"goals,omitempty"`
	GoalsIn          int     `yaml:"goals_in,omitempty"`
	GoalDiff         int     `yaml:"goal_diff,omitempty"`
	PointsPerGame    float32 `yaml:"points_per_game,omitempty"`
	PointsInPerGame  float32 `yaml:"points_in_per_game,omitempty"`
	GoalsWin         int     `yaml:"goals_win,omitempty"`
	DiffPerWin       float32 `yaml:"diff_per_win,omitempty"`
	GoalsInLoss      int     `yaml:"goals_in_loss,omitempty"`
	DiffPerLoss      float32 `yaml:"diff_per_loss,omitempty"`
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
