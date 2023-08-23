package entity

import (
	"strings"
)

// Player .
type Player struct {
	Name    string
	Aliases []string

	Played           int     `yaml:"played,omitempty"`
	TimePlayed       int     `yaml:"time_played,omitempty"`
	LongestGameTime  int     `yaml:"longest_game_time,omitempty"`
	ShortestGameTime int     `yaml:"shortest_game_time,omitempty"`
	TimePerGame      int     `yaml:"time_per_game,omitempty"`
	Won              int     `yaml:"won,omitempty"`
	Lost             int     `yaml:"lost,omitempty"`
	Draws            int     `yaml:"draws,omitempty"`
	HomeWon          int     `yaml:"home_won,omitempty"`
	HomeLost         int     `yaml:"home_lost,omitempty"`
	HomeWonRate      float32 `yaml:"home_won_rate,omitempty"`
	AwayWon          int     `yaml:"away_won,omitempty"`
	AwayLost         int     `yaml:"away_lost,omitempty"`
	AwayWonRate      float32 `yaml:"away_won_rate,omitempty"`
	WinRate          float32 `yaml:"win_rate,omitempty"`
	EloRating        float64 `yaml:"elo_rating,omitempty"`
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
