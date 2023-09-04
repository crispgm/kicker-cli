// Package model of kickertool
package model

// Modes
const (
	ModeSwissSystem       = "swiss"
	ModeRounds            = "rounds"
	ModeRoundRobin        = "round_robin"
	ModeMonsterDYP        = "monster_dyp"
	ModeDoubleElimination = "double_elimination"
	ModeElimination       = "elimination"
)

// Model an entity that contains ID and type
type Model struct {
	ID   string `json:"_id"`
	Type string `json:"type"`
}
