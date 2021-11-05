package elo

import "testing"

func TestEloScore(t *testing.T) {
	r := Rate{
		T1P1Score: 1205,
		T1P2Score: 1073,
		T2P1Score: 904,
		T2P2Score: 895,
		HostWin:   false,
	}
	r.CalcEloRating()
	t.Log(r)
}
