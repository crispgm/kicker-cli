package operator

import (
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

// calculateELO calculate ELO for player
func calculateELO(played int, p1Elo, p2Elo float64, result int) float64 {
	eloCalc := rating.Elo{}
	factors := rating.Factor{
		Played:        played,
		PlayerScore:   p1Elo,
		OpponentScore: p2Elo,
		Result:        result,
	}
	return eloCalc.Calculate(factors)
}

func openSingleTournament(trn *model.Tournament) bool {
	if trn.IsSingle() {
		if trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}

func openDoubleTournament(trn *model.Tournament) bool {
	if trn.IsDouble() {
		if trn.Mode == model.ModeMonsterDYP ||
			trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}
