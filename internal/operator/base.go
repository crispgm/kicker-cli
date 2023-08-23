// Package operator .
package operator

import (
	"github.com/crispgm/kickertool-analyzer/internal/entity"
)

// BaseOperator .
type BaseOperator interface {
	ValidMode(string) bool
	Output() [][]string
	Details() []entity.Player
}

// supportedOperator .
var supportedOperator = map[string]bool{
	entity.ModeMonsterDYPPlayerStats: true,
	entity.ModeMonsterDYPTeamStats:   true,
}

// IsSupported .
func IsSupported(mode string) bool {
	if supported, ok := supportedOperator[mode]; ok && supported {
		return true
	}
	return false
}

// Option .
type Option struct {
	OrderBy          string
	RankMinThreshold int
	EloKFactor       int
	WithTime         bool
	WithHomeAway     bool
	WithPoint        bool
	Incremental      bool
}
