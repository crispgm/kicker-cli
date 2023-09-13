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
	K = 40.0

	// WinScore is the score for winning
	WinScore = 1.0
	// DrawScore is the score for drawing
	DrawScore = 0.5
	// LossScore is the score for losing
	LossScore = 0.0
)

var _ rating.Rating = (*Elo)(nil)

// Elo calculates Elo rating
type Elo struct {
	ra, rb float64

	K float64
}

// InitialScore .
func (er *Elo) InitialScore(ra, rb float64) {
	er.ra = ra
	er.rb = rb
}

// Calculate elo rating based on scores
// https://math.stackexchange.com/questions/838809/rating-system-for-2-vs-2-2-vs-1-and-1-vs-1-game
func (er *Elo) Calculate(winDrawLoss int) float64 {
	ea := 1.0 / (1.0 + math.Pow(10, (er.rb-er.ra)/400.0))
	sa := WinScore
	if winDrawLoss == rating.Draw {
		sa = DrawScore
	} else if winDrawLoss == rating.Loss {
		sa = LossScore
	}
	ra := er.ra + er.K*(float64(sa)-ea)
	return ra
}
