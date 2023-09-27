package entity

import (
	"strings"

	"github.com/crispgm/kicker-cli/internal/util"
)

// Player is a real world player
type Player struct {
	ID      string   `yaml:"id"`
	Name    string   `yaml:"name"`
	Aliases []string `yaml:"aliases"`
	ATSAID  string   `yaml:"atsa_id,omitempty"`
	ITSFID  string   `yaml:"itsf_id,omitempty"`
	// major statistical data, not write
	EventsPlayed int     `yaml:"-"`
	GamesPlayed  int     `yaml:"-"`
	Win          int     `yaml:"-"`
	Loss         int     `yaml:"-"`
	Draw         int     `yaml:"-"`
	WinRate      float64 `yaml:"-"`
	EloRating    float64 `yaml:"-"` // kicker ELO scores
	KickerPoints int     `yaml:"-"` // kicker ranking points
	ATSAPoints   int     `yaml:"-"` // ATSA points
	ITSFPoints   int     `yaml:"-"` // ITSF points
	// minor statistical data
	HomeWin     int     `yaml:"-"`
	HomeLoss    int     `yaml:"-"`
	HomeWinRate float64 `yaml:"-"`
	AwayWin     int     `yaml:"-"`
	AwayLoss    int     `yaml:"-"`
	AwayWinRate float64 `yaml:"-"`
	// data not implemented right now
	TimePlayed       int     `yaml:"-"`
	LongestGameTime  int     `yaml:"-"`
	ShortestGameTime int     `yaml:"-"`
	TimePerGame      int     `yaml:"-"`
	Goals            int     `yaml:"-"`
	GoalsIn          int     `yaml:"-"`
	GoalDiff         int     `yaml:"-"`
	PointsPerGame    float64 `yaml:"-"`
	PointsInPerGame  float64 `yaml:"-"`
	GoalsWin         int     `yaml:"-"`
	DiffPerWin       float64 `yaml:"-"`
	GoalsInLoss      int     `yaml:"-"`
	DiffPerLoss      float64 `yaml:"-"`
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
