package stat

// BaseStat .
type BaseStat interface {
	ValidMode() bool
	Output() interface{}
}

// SupportedStat .
var SupportedStat = map[string]bool{
	"mts": true,
	"mtt": true,
}

// Option .
type Option struct {
	RankMinThreshold int
}
