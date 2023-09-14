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

	Name  string `json:"name"`
	Plays []Play `json:"plays"`
}
