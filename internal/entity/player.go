package entity

import (
	"strings"
)

// Player is a real world player
type Player struct {
	ID      string   `yaml:"id"`
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`

	Played           int     `yaml:"played"`
	TimePlayed       int     `yaml:"time_played"`
	LongestGameTime  int     `yaml:"longest_game_time"`
	ShortestGameTime int     `yaml:"shortest_game_time"`
	TimePerGame      int     `yaml:"time_per_game"`
	Won              int     `yaml:"won"`
	Lost             int     `yaml:"lost"`
	Draws            int     `yaml:"draws"`
	HomeWon          int     `yaml:"home_won"`
	HomeLost         int     `yaml:"home_lost"`
	HomeWonRate      float32 `yaml:"home_won_rate"`
	AwayWon          int     `yaml:"away_won"`
	AwayLost         int     `yaml:"away_lost"`
	AwayWonRate      float32 `yaml:"away_won_rate"`
	WinRate          float32 `yaml:"win_rate"`
	EloRating        float64 `yaml:"elo_rating"`
	Goals            int     `yaml:"goals"`
	GoalsIn          int     `yaml:"goals_in"`
	GoalDiff         int     `yaml:"goal_diff"`
	PointsPerGame    float32 `yaml:"points_per_game"`
	PointsInPerGame  float32 `yaml:"points_in_per_game"`
	GoalsWon         int     `yaml:"goals_won"`
	DiffPerWon       float32 `yaml:"diff_per_won"`
	GoalsInLost      int     `yaml:"goals_in_lost"`
	DiffPerLost      float32 `yaml:"diff_per_lost"`
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
