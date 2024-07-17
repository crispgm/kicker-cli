package model

// KnockOff .
type KnockOff struct {
	Model

	Levels     []Level `json:"levels"`
	LeftLevels []Level `json:"leftLevels"`
	Third      Level   `json:"third"`

	Options    Options `json:"options"`
	ThirdPlace bool    `json:"thirdPlace"`
	Double     bool    `json:"double"`
	Size       int     `json:"size"`
	Finished   bool    `json:"bool"`
}

// Level .
type Level struct {
	Model

	Name  string `json:"name"`
	Plays []Play `json:"plays"`
}
