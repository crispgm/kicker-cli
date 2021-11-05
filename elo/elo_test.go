package elo

import "testing"

func TestEloScore(t *testing.T) {
	r := Rate{
		T1P1Score: 1500,
		T1P2Score: 1400,
		T2P1Score: 1200,
		T2P2Score: 900,
		HostWin:   false,
	}
	r.CalcEloRating()
	t.Log(r)
}
