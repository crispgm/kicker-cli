package entity

import "github.com/crispgm/kicker-cli/internal/util"

// DefaultPoints for a event
const DefaultPoints = 50

// Event .
type Event struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Path   string `yaml:"path"`
	Points int    `yaml:"points"`
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
