package entity

// Record .
type Record struct {
	PreliminaryRounds []Game
	WinnerBracket     []Game
	LoserBracket      []Game
	ThirdPlace        *Game

	AllGames []Game
}
