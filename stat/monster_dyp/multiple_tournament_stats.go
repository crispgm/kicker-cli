package monsterdyp

import (
	"fmt"
	"sort"

	"github.com/crispgm/kickertool-analyzer/model"
)

// MultipleTournamentStats generate statistics data of multiple monster DYP tournaments
type MultipleTournamentStats struct {
	tournaments []model.Tournament
	players     []model.EntityPlayer
}

// NewMultipleTournamentStats .
func NewMultipleTournamentStats(tournaments []model.Tournament, players []model.EntityPlayer) *MultipleTournamentStats {
	return &MultipleTournamentStats{
		tournaments: tournaments,
		players:     players,
	}
}

// ValidMode .
func (m MultipleTournamentStats) ValidMode() bool {
	for _, t := range m.tournaments {
		if t.Mode != model.ModeMonsterDYP {
			return false
		}
	}

	return true
}

// Output .
func (m *MultipleTournamentStats) Output() []model.EntityPlayer {
	data := make(map[string]model.EntityPlayer)
	for _, t := range m.tournaments {
		teams := make(map[string]model.Team)
		players := make(map[string]model.Player)
		for _, p := range t.Players {
			if !p.Removed {
				var found bool
				for _, ep := range m.players {
					if ep.IsPlayer(p.Name) {
						found = true
						if _, ok := data[ep.Name]; !ok {
							data[ep.Name] = model.EntityPlayer{Name: ep.Name}
						}
						p.Name = ep.Name
						players[p.ID] = p
						break
					}
				}
				if !found {
					fmt.Println(p.Name, "not found")
				}
			}
		}
		for _, t := range t.Teams {
			teams[t.ID] = t
		}

		for _, r := range t.Rounds {
			if !r.Deactivated && !r.Skipped {
				for _, p := range r.Plays {
					if !p.Valid {
						continue
					}
					team1 := teams[p.Team1.ID]
					team2 := teams[p.Team2.ID]
					t1p1 := players[team1.Players[0].ID]
					t1p2 := players[team1.Players[1].ID]
					t2p1 := players[team2.Players[0].ID]
					t2p2 := players[team2.Players[1].ID]
					t1p1Data := data[t1p1.Name]
					t1p2Data := data[t1p2.Name]
					t2p1Data := data[t2p1.Name]
					t2p2Data := data[t2p2.Name]
					timePlayed := p.TimeEnd - p.TimeStart
					t1p1Data.TimePlayed += timePlayed
					t1p2Data.TimePlayed += timePlayed
					t2p1Data.TimePlayed += timePlayed
					t2p2Data.TimePlayed += timePlayed
					if t1p1Data.LongestGameTime < timePlayed || t1p1Data.LongestGameTime == 0 {
						t1p1Data.LongestGameTime = timePlayed
					}
					if t1p1Data.ShortestGameTime > timePlayed || t1p1Data.ShortestGameTime == 0 {
						t1p1Data.ShortestGameTime = timePlayed
					}
					if t1p2Data.LongestGameTime < timePlayed || t1p2Data.LongestGameTime == 0 {
						t1p2Data.LongestGameTime = timePlayed
					}
					if t1p2Data.ShortestGameTime > timePlayed || t1p2Data.ShortestGameTime == 0 {
						t1p2Data.ShortestGameTime = timePlayed
					}
					if t2p1Data.LongestGameTime < timePlayed || t2p1Data.LongestGameTime == 0 {
						t2p1Data.LongestGameTime = timePlayed
					}
					if t2p1Data.ShortestGameTime > timePlayed || t2p1Data.ShortestGameTime == 0 {
						t2p1Data.ShortestGameTime = timePlayed
					}
					if t2p2Data.LongestGameTime < timePlayed || t2p2Data.LongestGameTime == 0 {
						t2p2Data.LongestGameTime = timePlayed
					}
					if t2p2Data.ShortestGameTime > timePlayed || t2p2Data.ShortestGameTime == 0 {
						t2p2Data.ShortestGameTime = timePlayed
					}
					for _, d := range p.Disciplines {
						for _, s := range d.Sets {
							t1p1Data.Played++
							t1p2Data.Played++
							t2p1Data.Played++
							t2p2Data.Played++
							if s.Team1 > s.Team2 {
								t1p1Data.Won++
								t1p2Data.Won++
								t2p1Data.Lost++
								t2p2Data.Lost++
								t1p1Data.GoalsWon += (s.Team1 - s.Team2)
								t1p2Data.GoalsWon += (s.Team1 - s.Team2)
								t2p1Data.GoalsInLost += (s.Team1 - s.Team2)
								t2p2Data.GoalsInLost += (s.Team1 - s.Team2)
							} else if s.Team2 > s.Team1 {
								t1p1Data.Lost++
								t1p2Data.Lost++
								t2p1Data.Won++
								t2p2Data.Won++
								t2p1Data.GoalsWon += (s.Team2 - s.Team1)
								t2p2Data.GoalsWon += (s.Team2 - s.Team1)
								t1p1Data.GoalsInLost += (s.Team2 - s.Team1)
								t1p2Data.GoalsInLost += (s.Team2 - s.Team1)
							} else {
								t1p1Data.Draws++
								t1p2Data.Draws++
								t2p1Data.Draws++
								t2p2Data.Draws++
							}
							t1p1Data.Goals += s.Team1
							t1p2Data.Goals += s.Team1
							t2p1Data.Goals += s.Team2
							t2p2Data.Goals += s.Team2
							t1p1Data.GoalsIn += s.Team2
							t1p2Data.GoalsIn += s.Team2
							t2p1Data.GoalsIn += s.Team1
							t2p2Data.GoalsIn += s.Team1
						}
					}
					data[t1p1.Name] = t1p1Data
					data[t1p2.Name] = t1p2Data
					data[t2p1.Name] = t2p1Data
					data[t2p2.Name] = t2p2Data
				}
			}
		}
	}

	var sliceData []model.EntityPlayer
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Won) / float32(d.Played) * 100.0
			d.PointsPerGame = float32(d.Goals) / float32(d.Played)
			d.TimePerGame = d.TimePlayed / d.Played / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			d.DiffPerWon = float32(d.GoalsWon) / float32(d.Won)
			d.DiffPerLost = float32(d.GoalsInLost) / float32(d.Lost)
		}
		sliceData = append(sliceData, d)
	}
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].WinRate > sliceData[j].WinRate {
			return true
		} else if sliceData[i].WinRate == sliceData[j].WinRate {
			if sliceData[i].GoalDiff > sliceData[j].GoalDiff {
				return true
			} else if sliceData[i].GoalDiff == sliceData[j].GoalDiff {
				return sliceData[i].Goals > sliceData[j].Goals
			}
		}
		return false
	})
	return sliceData
}
