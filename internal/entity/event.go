package entity

import (
	"github.com/crispgm/kicker-cli/internal/util"
)

// DefaultPoints for a event
const DefaultPoints = 50

// Event .
type Event struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	MD5         string `yaml:"md5"`
	KickerLevel string `yaml:"kicker_level"`
	ITSFLevel   string `yaml:"itsf_level"`
	ATSALevel   string `yaml:"atsa_level"`
	URL         string `yaml:"url"`
}

// NewEvent creates an event
func NewEvent(path, name string, level string) *Event {
	return &Event{
		ID:          util.UUID(),
		Name:        name,
		Path:        path,
		KickerLevel: level,
	}
}
