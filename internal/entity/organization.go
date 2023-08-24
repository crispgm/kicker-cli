package entity

import (
	"github.com/crispgm/kicker-cli/internal/util"
)

// Organization .
type Organization struct {
	ID   string `yaml:"id"`
	Name string `yaml:"name"`
}

// NewOrganization creates an organization with name and UUID
func NewOrganization(name string) *Organization {
	return &Organization{
		ID:   util.UUID(),
		Name: name,
	}
}
