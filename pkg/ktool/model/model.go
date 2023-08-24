// Package model of kickertool
package model

import "time"

// Modes
const (
	ModeSwissSystem       = "swiss_system"
	ModeRounds            = "rounds"
	ModeRoundRobin        = "round_robin"
	ModeMonsterDYP        = "monster_dyp"
	ModeDoubleElimination = "double_elimination"
	ModeElimination       = "elimination"
)

// Tournament .
type Tournament struct {
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Version  string    `json:"version"`
	Mode     string    `json:"mode"`
	NameType string    `json:"nameType"`

	NumRounds int      `json:"numRounds"`
	Players   []Player `json:"players"`
	Teams     []Team   `json:"teams"`

	// for pre-eliminations
	Rounds []Round `json:"rounds"`

	// for eliminations
	KnockOffs []KnockOff `json:"ko"`
}

// Team .
type Team struct {
	ID      string `json:"_id"`
	Players []struct {
		ID string `json:"_id"`
	}
}

// Player .
type Player struct {
	ID               string `json:"_id"`
	Name             string `json:"_name"`
	Removed          bool   `json:"removed"`
	Deactivated      bool   `json:"deactivated"`
	MarkedForRemoval bool   `json:"markedForRemoval"`
}

// Round .
type Round struct {
	Name  string `json:"name"`
	Plays []Play `json:"plays,omitempty"`
}

// Play .
type Play struct {
	Valid       bool `json:"valid"`
	Deactivated bool `json:"deactivated"`
	Skipped     bool `json:"skipped"`
	TimeStart   int  `json:"timeStart"`
	TimeEnd     int  `json:"timeEnd"`
	Team1       struct {
		ID string `json:"_id"`
	}
	Team2 struct {
		ID string `json:"_id"`
	}
	Disciplines []Discipline `json:"disciplines"`
}

// Discipline .
type Discipline struct {
	Sets []Set `json:"sets,omitempty"`
}

// Set .
type Set struct {
	Team1 int `json:"team1"`
	Team2 int `json:"team2"`
}

// KnockOff .
type KnockOff struct {
	Levels     []Play `json:"levels"`
	LeftLevels []Play `json:"leftLevels"`
}
