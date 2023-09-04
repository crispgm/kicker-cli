package model

// KnockOff .
type KnockOff struct {
	Model

	Levels     []Level `json:"levels"`
	LeftLevels []Level `json:"leftLevels"`
	Third      Level   `json:"third"`
}

// Level .
type Level struct {
	Model

	Plays []Play `json:"plays"`
}
