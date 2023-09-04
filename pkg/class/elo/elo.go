// Package elo calculate ELO score
package elo

import "math"

const (
	// InitialScore of Elo rating
	InitialScore = 1000.0
	// K is the constant for Elo rating
	K = 10.0

	// WonScore is the score for winning
	WonScore = 1.0
	// DrawScore is the score for losing
	DrawScore = 0.5
	// LostScore is the score for losing
	LostScore = 0.0
)

// Rate holds Elo rating
type Rate struct {
	T1P1Score float64
	T1P1Exp   float64
	T1P2Score float64
	T1P2Exp   float64
	T2P1Score float64
	T2P1Exp   float64
	T2P2Score float64
	T2P2Exp   float64
	HostWin   bool
	K         float64
}

// CalcEloRating calc elo rating based on scores
// https://math.stackexchange.com/questions/838809/rating-system-for-2-vs-2-2-vs-1-and-1-vs-1-game
func (r *Rate) CalcEloRating() {
	// team average scores
	team1Score := (float64(r.T1P1Score) + float64(r.T1P2Score)) / 2
	team2Score := (float64(r.T2P1Score) + float64(r.T2P2Score)) / 2

	// expectations
	r.T1P1Exp = 1 / (1 + math.Pow(10, float64(team2Score-float64(r.T1P1Score))/400))
	r.T1P2Exp = 1 / (1 + math.Pow(10, float64(team2Score-float64(r.T1P2Score))/400))
	r.T2P1Exp = 1 / (1 + math.Pow(10, float64(team1Score-float64(r.T2P1Score))/400))
	r.T2P2Exp = 1 / (1 + math.Pow(10, float64(team1Score-float64(r.T2P2Score))/400))

	// update scores
	delta1Score := 0.0
	delta2Score := 0.0
	if r.HostWin {
		delta1Score = WonScore
	} else {
		delta2Score = WonScore
	}
	r.T1P1Score = math.Round(float64(r.T1P1Score) + r.K*(delta1Score-float64(r.T1P1Exp)))
	r.T1P2Score = math.Round(float64(r.T1P2Score) + r.K*(delta1Score-float64(r.T1P2Exp)))
	r.T2P1Score = math.Round(float64(r.T2P1Score) + r.K*(delta2Score-float64(r.T2P1Exp)))
	r.T2P2Score = math.Round(float64(r.T2P2Score) + r.K*(delta2Score-float64(r.T2P2Exp)))

	if r.T1P1Score < 0 {
		r.T1P1Score = 0
	}
	if r.T1P2Score < 0 {
		r.T1P2Score = 0
	}
	if r.T2P1Score < 0 {
		r.T2P1Score = 0
	}
	if r.T2P2Score < 0 {
		r.T2P2Score = 0
	}
}
