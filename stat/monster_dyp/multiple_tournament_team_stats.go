package monsterdyp

import (
	"fmt"
	"sort"

	"github.com/crispgm/kickertool-analyzer/model"
)

// MultipleTournamentTeamStats generate statistics data of multiple monster DYP tournaments by team
type MultipleTournamentTeamStats struct {
	tournaments []model.Tournament
	players     []model.EntityPlayer
}

// NewMultipleTournamentTeamStats .
func NewMultipleTournamentTeamStats(tournaments []model.Tournament, players []model.EntityPlayer) *MultipleTournamentTeamStats {
	return &MultipleTournamentTeamStats{
		tournaments: tournaments,
		players:     players,
	}
}

// ValidMode .
func (m MultipleTournamentTeamStats) ValidMode() bool {
	for _, t := range m.tournaments {
		if t.Mode != model.ModeMonsterDYP {
			return false
		}
	}

	return true
}

// Output .
func (m *MultipleTournamentTeamStats) Output() interface{} {
	data := make(map[string]model.EntityTeam)
	players := make(map[string]model.Player)
	for _, t := range m.tournaments {
		teams := make(map[string]model.Team)
		for _, p := range t.Players {
			if !p.Removed {
				var found bool
				for _, ep := range m.players {
					if ep.IsPlayer(p.Name) {
						found = true
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
			for _, p := range r.Plays {
				if !p.Valid || p.Deactivated || p.Skipped {
					continue
				}
				team1 := teams[p.Team1.ID]
				team2 := teams[p.Team2.ID]
				t1p1Name := players[team1.Players[0].ID].Name
				t1p2Name := players[team1.Players[1].ID].Name
				t2p1Name := players[team2.Players[0].ID].Name
				t2p2Name := players[team2.Players[1].ID].Name
				team1Name := fmt.Sprintf("%s_%s", t1p1Name, t1p2Name)
				if t1p1Name > t1p2Name {
					team1Name = fmt.Sprintf("%s_%s", t1p2Name, t1p1Name)
				}
				team2Name := fmt.Sprintf("%s_%s", t2p1Name, t2p2Name)
				if t2p1Name > t2p2Name {
					team2Name = fmt.Sprintf("%s_%s", t2p2Name, t2p1Name)
				}
				var et1, et2 model.EntityTeam
				if t, ok := data[team1Name]; ok {
					et1 = t
				} else {
					et1 = model.EntityTeam{
						Player1: t1p1Name,
						Player2: t1p2Name,
					}
				}
				if t, ok := data[team2Name]; ok {
					et2 = t
				} else {
					et2 = model.EntityTeam{
						Player1: t2p1Name,
						Player2: t2p2Name,
					}
				}
				timePlayed := p.TimeEnd - p.TimeStart
				et1.Played++
				et2.Played++
				et1.TimePlayed += timePlayed
				et2.TimePlayed += timePlayed
				m.playedTimeStats(&et1, timePlayed)
				m.playedTimeStats(&et2, timePlayed)
				for _, d := range p.Disciplines {
					for _, s := range d.Sets {
						if s.Team1 > s.Team2 {
							et1.Won++
							et2.Lost++
							et1.GoalsWon += (s.Team1 - s.Team2)
							et2.GoalsInLost += (s.Team1 - s.Team2)
						} else if s.Team2 > s.Team1 {
							et1.Lost++
							et2.Won++
							et2.GoalsWon += (s.Team2 - s.Team1)
							et1.GoalsInLost += (s.Team2 - s.Team1)
						} else {
							et1.Draws++
							et2.Draws++
						}
						et1.Goals += s.Team1
						et2.Goals += s.Team2
						et1.GoalsIn += s.Team2
						et2.GoalsIn += s.Team1
					}
				}
				data[team1Name] = et1
				data[team2Name] = et2
			}
		}
	}

	var sliceData []model.EntityTeam
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Won) / float32(d.Played) * 100.0
			d.PointsPerGame = float32(d.Goals) / float32(d.Played)
			d.PointsInPerGame = float32(d.GoalsIn) / float32(d.Played)
			d.TimePerGame = d.TimePlayed / d.Played / 1000
			d.LongestGameTime /= 1000
			d.ShortestGameTime /= 1000
			d.DiffPerWon = float32(d.GoalsWon) / float32(d.Won)
			d.DiffPerLost = float32(d.GoalsInLost) / float32(d.Lost)
		}
		sliceData = append(sliceData, d)
	}
	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= rankThreshold && sliceData[j].Played < rankThreshold {
			return true
		}
		if sliceData[i].Played < rankThreshold && sliceData[j].Played >= rankThreshold {
			return false
		}

		if sliceData[i].WinRate > sliceData[j].WinRate {
			return true
		} else if sliceData[i].WinRate == sliceData[j].WinRate {
			iWinLost := sliceData[i].Won - sliceData[i].Lost
			jWinLost := sliceData[j].Won - sliceData[j].Lost
			if iWinLost > jWinLost {
				return true
			} else if iWinLost == jWinLost {
				return sliceData[i].Won > sliceData[j].Won
			}
		}
		return false
	})
	return sliceData
}

func (MultipleTournamentTeamStats) playedTimeStats(data *model.EntityTeam, timePlayed int) {
	if timePlayed < 0 || timePlayed > 1000*60*15 {
		// consider illegal
		return
	}
	if data.LongestGameTime < timePlayed || data.LongestGameTime == 0 {
		data.LongestGameTime = timePlayed
	}
	if data.ShortestGameTime > timePlayed || data.ShortestGameTime == 0 {
		data.ShortestGameTime = timePlayed
	}
}
