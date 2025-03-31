package rating

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRank(t *testing.T) {
	r := Rank{}
	f := Factor{
		Level:       ATSA50,
		Place:       1,
		PlayerScore: 0.0,
	}

	assert.Zero(t, r.InitialScore())
	assert.Equal(t, 50.0, r.Calculate(f))

	f.Place = 2
	assert.Equal(t, 35.0, r.Calculate(f))
	f.Place = 3
	assert.Equal(t, 25.0, r.Calculate(f))
	f.Place = 4
	assert.Equal(t, 20.0, r.Calculate(f))
	f.Place = 5
	assert.Equal(t, 15.0, r.Calculate(f))
	f.Place = 6
	assert.Equal(t, 15.0, r.Calculate(f))
	f.Place = 7
	assert.Equal(t, 9.0, r.Calculate(f))
	f.Place = 9
	assert.Equal(t, 8.0, r.Calculate(f))
}
