package entity

import (
	"github.com/crispgm/kicker-cli/internal/util"
)

// DefaultPoints for a event
const DefaultPoints = 50

// Event .
type Event struct {
	ID         string `yaml:"id"`
	Name       string `yaml:"name"`
	Path       string `yaml:"path"`
	MD5        string `yaml:"md5"`
	Points     int    `yaml:"points"`
	ITSFPoints int    `yaml:"itsf_points"`
	ATSAPoints int    `yaml:"atsa_points"`
	URL        string `yaml:"url"`
}

// NewEvent creates an event
func NewEvent(path, name string, points int) *Event {
	return &Event{
		ID:     util.UUID(),
		Name:   name,
		Path:   path,
		Points: points,
	}
}
