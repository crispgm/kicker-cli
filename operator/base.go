package operator

import "github.com/crispgm/kickertool-analyzer/model"

// BaseOperator .
type BaseOperator interface {
	ValidMode(string) bool
	Output() [][]string
}

// SupportedOperator .
var SupportedOperator = map[string]bool{
	model.ModeMonsterDYPPlayerStats: true,
	model.ModeMonsterDYPTeamStats:   true,
}

// Option .
type Option struct {
	RankMinThreshold int
	WithTime         bool
	WithHomeAway     bool
}
