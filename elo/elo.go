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

// Class for player score
type Class struct {
	Low   int
	High  int
	Title string
}

// Rate holds Elo rating
type Rate struct {
	T1P1Score float64
	T1P2Score float64
	T2P1Score float64
	T2P2Score float64
	HostWin   bool
}

// CalcEloRating calc elo rating based on scores
// https://math.stackexchange.com/questions/838809/rating-system-for-2-vs-2-2-vs-1-and-1-vs-1-game
func (r *Rate) CalcEloRating() {
	// team average scores
	team1Score := (float64(r.T1P1Score) + float64(r.T1P2Score)) / 2
	team2Score := (float64(r.T2P1Score) + float64(r.T2P2Score)) / 2

	// expectations
	t1p1exp := 1 / (1 + math.Pow(10, float64(team2Score-float64(r.T1P1Score))/400))
	t1p2exp := 1 / (1 + math.Pow(10, float64(team2Score-float64(r.T1P2Score))/400))
	t2p1exp := 1 / (1 + math.Pow(10, float64(team1Score-float64(r.T2P1Score))/400))
	t2p2exp := 1 / (1 + math.Pow(10, float64(team1Score-float64(r.T2P2Score))/400))

	// update scores
	delta1Score := 0.0
	delta2Score := 0.0
	if r.HostWin {
		delta1Score = WonScore
	} else {
		delta2Score = WonScore
	}
	r.T1P1Score = math.Round(float64(r.T1P1Score) + K*(delta1Score-float64(t1p1exp)))
	r.T1P2Score = math.Round(float64(r.T1P2Score) + K*(delta1Score-float64(t1p2exp)))
	r.T2P1Score = math.Round(float64(r.T2P1Score) + K*(delta2Score-float64(t2p1exp)))
	r.T2P2Score = math.Round(float64(r.T2P2Score) + K*(delta2Score-float64(t2p2exp)))
}
