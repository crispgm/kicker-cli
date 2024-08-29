package entity

import (
	"time"

	"github.com/crispgm/kicker-cli/internal/util"
)

// Organization .
type Organization struct {
	ID           string `yaml:"id"`
	Name         string `yaml:"name"`
	Timezone     string `yaml:"timezone"`
	KickerToolID string `yaml:"kicker_tool_id"`
}

// NewOrganization creates an organization with name and UUID
func NewOrganization(name string) *Organization {
	curTime := time.Now()
	return &Organization{
		ID:       util.UUID(),
		Name:     name,
		Timezone: curTime.Location().String(),
	}
}
