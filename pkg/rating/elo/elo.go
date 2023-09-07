// Package elo calculate ELO score
package elo

import (
	"math"

	"github.com/crispgm/kicker-cli/pkg/rating"
)

const (
	// InitialScore of Elo rating
	InitialScore = 1500.0
	// K is the constant for Elo rating
	K = 10.0

	// WonScore is the score for winning
	WonScore = 1.0
	// DrewScore is the score for drawing
	DrewScore = 0.5
	// LostScore is the score for losing
	LostScore = 0.0
)

var _ rating.Rating = (*EloRating)(nil)

// EloRating calculates Elo rating
type EloRating struct {
	ra, rb float64

	K float64
}

// InitialScore .
func (er *EloRating) InitialScore(ra, rb float64) {
	er.ra = ra
	er.rb = rb
}

// Calculate elo rating based on scores
// https://math.stackexchange.com/questions/838809/rating-system-for-2-vs-2-2-vs-1-and-1-vs-1-game
func (er *EloRating) Calculate(wonDrewLost int) float64 {
	ea := 1.0 / (1.0 + math.Pow(10, (er.rb-er.ra)/400.0))
	sa := WonScore
	if wonDrewLost == rating.Drew {
		sa = DrewScore
	} else if wonDrewLost == rating.Lost {
		sa = LostScore
	}
	ra := er.ra + er.K*(float64(sa)-ea)
	return ra
}
