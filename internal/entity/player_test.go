package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPlayer(t *testing.T) {
	p := NewPlayer("Aaa Bbbb")
	p.AddAlias("A", "bb")

	assert.True(t, p.IsPlayer("Aaa bbbb"))
	assert.False(t, p.IsPlayer("Aaa bbb"))
	assert.True(t, p.IsPlayer("A"))
	assert.False(t, p.IsPlayer("B"))
	assert.True(t, p.IsPlayer("BB"))
}
