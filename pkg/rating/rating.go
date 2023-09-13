// Package rating multiple algorithms for rating
package rating

// NotSanctioned represents event that not sanctioned by this organization
const NotSanctioned = "NS"

// ITSF Points
const (
	ITSFWorldSeries   = "ITSFWorldSeries"
	ITSFInternational = "ITSFInternational"
	ITSFMasterSeries  = "ITSFWorldSeries"
	ITSFProTour       = "ITSFProTour"
)

// ATSA Points
const (
	ATSA2000 = "ATSA2000"
	ATSA1000 = "ATSA1000"
	ATSA500  = "ATSA500"
	ATSA50   = "ATSA50"
	ATSA25   = "ATSA25"
)

// literally win/draw/loss
const (
	Loss = iota + 1
	Draw
	Win
)

// Rating interface of rating system
type Rating interface {
	InitialScore(float64, float64)
	Calculate(int) float64
}
