package operator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupported(t *testing.T) {
	assert.True(t, IsSupported("mdp"))
	assert.True(t, IsSupported("mdt"))
	assert.False(t, IsSupported("mds"))
}
