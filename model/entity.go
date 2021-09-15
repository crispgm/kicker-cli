package model

import "strings"

// Supported modes
const (
	ModeMonsterDYPPlayerStats = "mdp"
	ModeMonsterDYPTeamStats   = "mdt"
)

// EntityGame is stat for single games
type EntityGame struct {
	Team1 []string
	Team2 []string

	TimePlayed int
	Point1     int
	Point2     int
}

// EntityTeam .
type EntityTeam struct {
	Player1          string
	Player2          string
	Played           int
	TimePlayed       int
	Won              int
	Lost             int
	Draws            int
	Goals            int
	GoalsIn          int
	GoalDiff         int
	WinRate          float32
	TimePerGame      int
	PointsPerGame    float32
	PointsInPerGame  float32
	GoalsWon         int
	DiffPerWon       float32
	GoalsInLost      int
	DiffPerLost      float32
	LongestGameTime  int
	ShortestGameTime int
}

// EntityPlayer .
type EntityPlayer struct {
	Name    string
	Aliases []string

	Played int

	TimePlayed       int
	LongestGameTime  int
	ShortestGameTime int
	TimePerGame      int

	Won         int
	Lost        int
	Draws       int
	HomeWon     int
	HomeLost    int
	HomeWonRate float32
	AwayWon     int
	AwayLost    int
	AwayWonRate float32
	WinRate     float32

	Goals           int
	GoalsIn         int
	GoalDiff        int
	PointsPerGame   float32
	PointsInPerGame float32
	GoalsWon        int
	DiffPerWon      float32
	GoalsInLost     int
	DiffPerLost     float32
}

// IsPlayer .
func (p EntityPlayer) IsPlayer(name string) bool {
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
