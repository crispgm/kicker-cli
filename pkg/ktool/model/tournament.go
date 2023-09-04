package model

import "time"

// Tournament .
type Tournament struct {
	Model

	Sport    Sport     `json:"sport"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Version  string    `json:"version"`
	Mode     string    `json:"mode"`
	NameType string    `json:"nameType"`
	Options  Options   `json:"options"`

	NumRounds int      `json:"numRounds"`
	Players   []Player `json:"players"`
	Teams     []Team   `json:"teams"`

	// for pre-eliminations
	Rounds []Round `json:"rounds"`

	// for eliminations
	KnockOffs []KnockOff `json:"ko"`
}

// Sport .
type Sport struct {
	Name string `json:"name"`
}

// Options .
type Options struct {
	Model

	Name       string `json:"name"`
	NumPoints  int    `json:"numPoints"`
	NumSets    int    `json:"numSets"`
	TwoAhead   bool   `json:"twoAhead"`
	Draw       bool   `json:"draw"`
	PointsWin  int    `json:"pointsWin"`
	PointsDraw int    `json:"pointsDraw"`
}

// IsSingle .
func (t Tournament) IsSingle() bool {
	return t.NameType == "single"
}

// IsDouble .
func (t Tournament) IsDouble() bool {
	return !t.IsSingle()
}

// IsBYP .
func (t Tournament) IsBYP() bool {
	return t.NameType == "byp"
}

// IsDYP .
func (t Tournament) IsDYP() bool {
	return t.NameType == "dyp"
}
