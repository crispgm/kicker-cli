package model

import (
	"fmt"
	"time"
)

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
	FastInput  bool   `json:"fastInput"`
	Draw       bool   `json:"draw"`
	PointsWin  int    `json:"pointsWin"`
	PointsDraw int    `json:"pointsDraw"`
}

// IsSingle .
func (t Tournament) IsSingle() bool {
	return t.NameType == NameTypeSingle
}

// IsDouble .
func (t Tournament) IsDouble() bool {
	return !t.IsSingle()
}

// IsBYP .
func (t Tournament) IsBYP() bool {
	return t.NameType == NameTypeBYP
}

// IsDYP .
func (t Tournament) IsDYP() bool {
	return t.NameType == NameTypeDYP
}

// IsMonsterDYP .
func (t Tournament) IsMonsterDYP() bool {
	return t.NameType == NameTypeMonsterDYP
}

// PreliminaryMode .
func (t Tournament) PreliminaryMode() string {
	if len(t.Rounds) == 0 {
		return ""
	}

	return t.Mode
}

// EliminationMode .
func (t Tournament) EliminationMode() string {
	if len(t.KnockOffs) == 0 {
		return ""
	}

	ko := t.KnockOffs[0]
	if len(ko.LeftLevels) > 0 {
		return ModeDoubleElimination
	}

	return ModeElimination
}

// TournamentMode .
func (t Tournament) TournamentMode() string {
	eMode := t.EliminationMode()
	pMode := t.PreliminaryMode()
	if eMode == "" && pMode == "" {
		return t.Mode
	}

	if eMode == "" {
		return pMode
	} else if pMode == "" {
		return eMode
	}

	return fmt.Sprintf("%s, %s", pMode, eMode)
}
