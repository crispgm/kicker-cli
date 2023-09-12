package model

// Round .
type Round struct {
	Model

	Name  string `json:"name"`
	Plays []Play `json:"plays,omitempty"`
}

// Play .
type Play struct {
	Model

	Valid       bool `json:"valid"`
	Deactivated bool `json:"deactivated"`
	Skipped     bool `json:"skipped"`
	TimeStart   int  `json:"timeStart"`
	TimeEnd     int  `json:"timeEnd"`
	Team1       struct {
		ID   string `json:"_id"`
		Type string `json:"type"`
	}
	Team2 struct {
		ID   string `json:"_id"`
		Type string `json:"type"`
	}
	Disciplines []Discipline `json:"disciplines"`
}

// Discipline .
type Discipline struct {
	Model

	Sets []Set `json:"sets,omitempty"`
}

// Set .
type Set struct {
	Model

	Team1 int `json:"team1"`
	Team2 int `json:"team2"`
}
