package rating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRatingIsMethods(t *testing.T) {
	f := Factor{Level: ATSA1000}
	assert.Equal(t, RSysATSA, f.GetRankSystem())
	assert.True(t, f.IsATSA())
	assert.False(t, f.IsITSF())
}
