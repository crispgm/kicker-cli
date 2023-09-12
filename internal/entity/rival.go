package entity

// Rival .
type Rival struct {
	Team1 Team
	Team2 Team

	Played      int
	Win         int
	Loss        int
	Draw        int
	TimePlayed  int
	TimePerGame int
}
