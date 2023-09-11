// Package rating multiple algorithms for rating
package rating

// literally win/draw/loss
const (
	Win = iota + 1
	Draw
	Loss
)

// Rating interface of rating system
type Rating interface {
	InitialScore(float64, float64)
	Calculate(int) float64
}
