package rating

import (
	"math"
)

const (
	// InitialScore of Elo rating
	initialScore = 1500.0
	// K is the constant for Elo rating
	defaultK = 40.0

	winScore  = 1.0
	drawScore = 0.5
	lossScore = 0.0
)

var _ Rating = (*Elo)(nil)

// Elo calculates Elo rating
type Elo struct {
	K int
}

// InitialScore .
func (e Elo) InitialScore() float64 {
	return initialScore
}

func (e Elo) chooseK(played int, score float64) float64 {
	k := defaultK
	if played >= 30 {
		if score >= 2400 {
			k = 10.0
		} else {
			k = 20.0
		}
	} else {
		k = 40.0
	}
	return k
}

// Calculate elo rating based on scores
// https://math.stackexchange.com/questions/838809/rating-system-for-2-vs-2-2-vs-1-and-1-vs-1-game
func (e Elo) Calculate(factors Factor) float64 {
	k := e.chooseK(factors.Played, factors.PlayerScore)
	if e.K > 0 {
		k = float64(e.K)
	}
	ra := factors.PlayerScore
	rb := factors.OpponentScore
	ea := 1.0 / (1.0 + math.Pow(10, (rb-ra)/400.0))
	sa := winScore
	if factors.Result == Draw {
		sa = drawScore
	} else if factors.Result == Loss {
		sa = lossScore
	}
	return ra + k*(float64(sa)-ea)
}
