// Package operator .
package operator

import (
	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

// Operator .
type Operator interface {
	SupportedFormats(trn *model.Tournament) bool
	Input(tournaments []entity.Tournament, players []entity.Player, options Option)
	Output() [][]string
}

// Option .
type Option struct {
	OrderBy       string
	MinimumPlayed int
	Head          int
	Tail          int
	WithHeader    bool
}
