package entity

// Team .
type Team struct {
	Player1          string
	Player2          string
	Played           int
	TimePlayed       int
	Won              int
	Lost             int
	Draws            int
	Goals            int
	GoalsIn          int
	GoalDiff         int
	WinRate          float32
	TimePerGame      int
	PointsPerGame    float32
	PointsInPerGame  float32
	GoalsWon         int
	DiffPerWon       float32
	GoalsInLost      int
	DiffPerLost      float32
	LongestGameTime  int
	ShortestGameTime int
}
