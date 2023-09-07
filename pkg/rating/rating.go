// Package rating multiple algorithms for rating
package rating

// literally won/draw/lost
const (
	Won = iota + 1
	Drew
	Lost
)

// Rating interface of rating system
type Rating interface {
	InitialScore(float64, float64)
	Calculate(int) float64
}
