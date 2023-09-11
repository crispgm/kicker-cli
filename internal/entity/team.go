package entity

// Team .
type Team struct {
	Player1          string
	Player2          string
	Played           int
	TimePlayed       int
	Win              int
	Draw             int
	Loss             int
	Goals            int
	GoalsIn          int
	GoalDiff         int
	WinRate          float32
	TimePerGame      int
	PointsPerGame    float32
	PointsInPerGame  float32
	GoalsWin         int
	DiffPerWin       float32
	GoalsInLoss      int
	DiffPerLoss      float32
	LongestGameTime  int
	ShortestGameTime int
}
