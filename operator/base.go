package operator

import "github.com/crispgm/kickertool-analyzer/model"

// BaseOperator .
type BaseOperator interface {
	ValidMode(string) bool
	Output() [][]string
	Details() []model.EntityPlayer
}

// supportedOperator .
var supportedOperator = map[string]bool{
	model.ModeMonsterDYPPlayerStats: true,
	model.ModeMonsterDYPTeamStats:   true,
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
	WithTime         bool
	WithHomeAway     bool
	WithPoint        bool
	Incremental      bool
}
