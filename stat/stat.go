package stat

import "github.com/crispgm/kickertool-analyzer/model"

// BaseStat .
type BaseStat interface {
	ValidMode() bool
	Output() []model.EntityPlayer
}

// SupportedStat .
var SupportedStat = map[string]bool{
	"mts": true,
}
