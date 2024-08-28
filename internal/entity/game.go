package entity

// Game is stat for single game
type Game struct {
	Team1 []string
	Team2 []string

	GameType   int // 1: qualification, 2: elimination
	TimeStart  int
	TimeEnd    int
	TimePlayed int
	Winner     int
	Point1     int
	Point2     int

	Name string
	Sets []Set
}

// Game types
const (
	GameTypeQualification = iota + 1
	GameTypeElimination
)

// Set is stat for single set
type Set struct {
	Point1 int
	Point2 int
}
