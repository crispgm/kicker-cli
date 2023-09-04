package model

// Player .
type Player struct {
	Model

	Name             string `json:"_name"`
	Removed          bool   `json:"removed"`
	Deactivated      bool   `json:"deactivated"`
	MarkedForRemoval bool   `json:"markedForRemoval"`
}
