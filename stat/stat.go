package stat

// BaseStat .
type BaseStat interface {
	ValidMode() bool
	Output() [][]string
}

// SupportedStat .
var SupportedStat = map[string]bool{
	"mts": true,
	"mtt": true,
}

// Option .
type Option struct {
	RankMinThreshold int
	WithTime         bool
	WithHomeAway     bool
}
