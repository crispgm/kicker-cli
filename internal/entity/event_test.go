package entity

import (
	"testing"

	"github.com/crispgm/kicker-cli/pkg/rating"
	"github.com/stretchr/testify/assert"
)

func TestNewEvent(t *testing.T) {
	p := NewEvent(".", "test", rating.KWorld)

	assert.NotEmpty(t, p.ID)
	assert.Equal(t, p.Name, "test")
	assert.Equal(t, rating.KWorld, p.KickerLevel)
}
