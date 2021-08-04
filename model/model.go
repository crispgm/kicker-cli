package model

// Modes
const (
	ModeMonsterDYP = "monster_dyp"
)

// Tournament .
type Tournament struct {
	Name      string   `json:"name"`
	Mode      string   `json:"mode"`
	NumRounds int      `json:"numRounds"`
	Players   []Player `json:"players,omitempty"`
	Teams     []Team
	Rounds    []Round `json:"rounds,omitempty"`
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
	Name        string `json:"name"`
	Deactivated bool   `json:"deactivated"`
	Skipped     bool   `json:"skipped"`
	Plays       []Play `json:"plays,omitempty"`
}

// Play .
type Play struct {
	Valid     bool `json:"valid"`
	TimeStart int  `json:"timeStart"`
	TimeEnd   int  `json:"timeEnd"`
	Team1     struct {
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
