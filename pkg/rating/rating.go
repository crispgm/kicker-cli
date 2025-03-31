// Package rating multiple algorithms for rating
package rating

import "strings"

// ranking system
const (
	RSysWinRate = "WR"
	RSysELO     = "ELO"
	RSysITSF    = "ITSF"
	RSysATSA    = "ATSA"
	RSysKicker  = "KRP"
)

// NotSanctioned represents event that not sanctioned by this organization
const NotSanctioned = "NS"

// Kicker Points
const (
	KWorld       = "KWorld"
	KContinental = "KContinental"
	KDomestic    = "KDomestic"
	KLocal       = "KLocal"
	KCasual      = "KCasual"
)

// ITSF Points
const (
	ITSFWorldSeries   = "ITSFWorldSeries"
	ITSFInternational = "ITSFInternational"
	ITSFMasterSeries  = "ITSFMasterSeries"
	ITSFProTour       = "ITSFProTour"
)

// ATSA Points
const (
	ATSA2000 = "ATSA2000"
	ATSA1000 = "ATSA1000"
	ATSA750  = "ATSA750"
	ATSA500  = "ATSA500"
	ATSA250  = "ATSA250"
	ATSA50   = "ATSA50"
)

// literally win/draw/loss
const (
	Loss = iota + 1
	Draw
	Win
)

// Factor calculation variables
type Factor struct {
	PlayerScore   float64 // player score
	OpponentScore float64 // opponent player/team score
	Result        int     // game result
	Level         string  // tournament/game level
	Place         int     // place in tournament
	Played        int     // game played
}

// IsATSA .
func (f Factor) IsATSA() bool {
	return strings.HasPrefix(f.Level, "ATSA")
}

// IsITSF .
func (f Factor) IsITSF() bool {
	return strings.HasPrefix(f.Level, "ITSF")
}

// GetRankSystem .
func (f Factor) GetRankSystem() string {
	if f.IsITSF() {
		return RSysITSF
	} else if f.IsATSA() {
		return RSysATSA
	}

	return RSysKicker
}

// Rating interface of rating system
type Rating interface {
	InitialScore() float64
	Calculate(factors Factor) float64
}
