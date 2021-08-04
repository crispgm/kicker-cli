package model

import "strings"

// EntityPlayer .
type EntityPlayer struct {
	Name    string
	Aliases []string

	Played        int
	TimePlayed    int
	Won           int
	Lost          int
	Draws         int
	Goals         int
	GoalsIn       int
	GoalDiff      int
	WinRate       float32
	TimePerGame   int
	PointsPerGame float32
	GoalsWon      int
	DiffPerWon    float32
	GoalsInLost   int
	DiffPerLost   float32
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
