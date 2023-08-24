package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateOrg(t *testing.T) {
	org := NewOrganization("test")
	assert.Equal(t, "test", org.Name)
	assert.NotEmpty(t, org.ID)
}
