// Package operator .
package operator

import "github.com/crispgm/kicker-cli/internal/entity"

// Operator .
type Operator interface {
	SupportedFormats(nameType string, mode string) bool
	Input(games []entity.Game, players []entity.Player, options Option)
	Output() [][]string
}

// Option .
type Option struct {
	OrderBy          string
	RankMinThreshold int
	EloKFactor       int
	WithHeader       bool
	WithTime         bool
	WithHomeAway     bool
	WithGoals        bool
}
