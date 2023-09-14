package entity

// Game is stat for single game
type Game struct {
	Team1 []string
	Team2 []string

	TimeStart  int
	TimeEnd    int
	TimePlayed int
	Point1     int
	Point2     int

	Name string
	Sets []Set
}

// Set is stat for single set
type Set struct {
	Point1 int
	Point2 int
}
