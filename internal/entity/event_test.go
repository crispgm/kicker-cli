package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	p := NewEvent(".", "test", 50)

	assert.NotEmpty(t, p.ID)
	assert.Equal(t, p.Name, "test")
	assert.Equal(t, p.Points, 50)
}
