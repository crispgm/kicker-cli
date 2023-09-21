package entity

import "github.com/crispgm/kicker-cli/pkg/ktool/model"

// Tournament .
type Tournament struct {
	Raw       model.Tournament
	Converted Record
}

// Record .
type Record struct {
	PreliminaryRounds []Game
	WinnerBracket     []Game
	LoserBracket      []Game
	ThirdPlace        *Game

	AllGames []Game
	Ranks    [][]Player
	Players  []Player
}
