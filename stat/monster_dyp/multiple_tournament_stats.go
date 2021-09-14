package monsterdyp

import (
	"fmt"
	"sort"

	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/stat"
)

// MultipleTournamentStats generate statistics data of multiple monster DYP tournaments
type MultipleTournamentStats struct {
	option      stat.Option
	tournaments []model.Tournament
	players     []model.EntityPlayer
}

// NewMultipleTournamentStats .
func NewMultipleTournamentStats(tournaments []model.Tournament, players []model.EntityPlayer, option stat.Option) *MultipleTournamentStats {
	return &MultipleTournamentStats{
		option:      option,
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
func (m *MultipleTournamentStats) Output() [][]string {
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
			for _, p := range r.Plays {
				if !p.Valid || p.Deactivated || p.Skipped {
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
				m.playedTimeStats(&t1p1Data, timePlayed)
				m.playedTimeStats(&t1p2Data, timePlayed)
				m.playedTimeStats(&t2p1Data, timePlayed)
				m.playedTimeStats(&t2p2Data, timePlayed)
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
							t1p1Data.HomeWon++
							t1p2Data.HomeWon++
							t2p1Data.AwayLost++
							t2p2Data.AwayLost++
							t1p1Data.GoalsWon += (s.Team1 - s.Team2)
							t1p2Data.GoalsWon += (s.Team1 - s.Team2)
							t2p1Data.GoalsInLost += (s.Team1 - s.Team2)
							t2p2Data.GoalsInLost += (s.Team1 - s.Team2)
						} else if s.Team2 > s.Team1 {
							t1p1Data.Lost++
							t1p2Data.Lost++
							t2p1Data.Won++
							t2p2Data.Won++
							t1p1Data.HomeLost++
							t1p2Data.HomeLost++
							t2p1Data.AwayWon++
							t2p2Data.AwayWon++
							t2p1Data.GoalsWon += (s.Team2 - s.Team1)
							t2p2Data.GoalsWon += (s.Team2 - s.Team1)
							t1p1Data.GoalsInLost += (s.Team2 - s.Team1)
							t1p2Data.GoalsInLost += (s.Team2 - s.Team1)
						} else {
							// basically not approachable
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

	var sliceData []model.EntityPlayer
	for _, d := range data {
		d.GoalDiff = d.Goals - d.GoalsIn
		if d.Played != 0 {
			d.WinRate = float32(d.Won) / float32(d.Played) * 100.0
			d.HomeWonRate = float32(d.HomeWon) / float32(d.HomeWon+d.HomeLost) * 100.0
			d.AwayWonRate = float32(d.AwayWon) / float32(d.AwayWon+d.AwayLost) * 100.0
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
		if sliceData[i].Played >= m.option.RankMinThreshold && sliceData[j].Played < m.option.RankMinThreshold {
			return true
		}
		if sliceData[i].Played < m.option.RankMinThreshold && sliceData[j].Played >= m.option.RankMinThreshold {
			return false
		}

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

	header := []string{"#", "Name", "Num", "Won", "Lost", "G+", "G-", "G±", "WR%", "PPG", "LPG", "DPW", "DPL"}
	haHeader := []string{"HW", "HL", "HW%", "AW", "AL", "AW%"}
	timeHeader := []string{"TPG", "LGP", "SGP"}
	if m.option.WithHostAway {
		header = append(header, haHeader...)
	}
	if m.option.WithTime {
		header = append(header, timeHeader...)
	}
	table := [][]string{header}
	for i, d := range sliceData {
		item := []string{
			fmt.Sprintf("%d", i+1),
			d.Name,
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Won),
			fmt.Sprintf("%d", d.Lost),
			fmt.Sprintf("%d", d.Goals),
			fmt.Sprintf("%d", d.GoalsIn),
			fmt.Sprintf("%d", d.GoalDiff),
			fmt.Sprintf("%.0f%%", d.WinRate),
			fmt.Sprintf("%.2f", d.PointsPerGame),
			fmt.Sprintf("%.2f", d.PointsInPerGame),
			fmt.Sprintf("%.2f", d.DiffPerWon),
			fmt.Sprintf("%.2f", d.DiffPerLost),
		}
		if m.option.WithHostAway {
			item = append(item, []string{
				fmt.Sprintf("%d", d.HomeWon),
				fmt.Sprintf("%d", d.HomeLost),
				fmt.Sprintf("%.0f%%", d.HomeWonRate),
				fmt.Sprintf("%d", d.AwayWon),
				fmt.Sprintf("%d", d.HomeLost),
				fmt.Sprintf("%.0f%%", d.AwayWonRate),
			}...)
		}
		if m.option.WithTime {
			item = append(item, []string{
				fmt.Sprintf("%02d:%02d", d.TimePerGame/60, d.TimePerGame%60),
				fmt.Sprintf("%02d:%02d", d.LongestGameTime/60, d.LongestGameTime%60),
				fmt.Sprintf("%02d:%02d", d.ShortestGameTime/60, d.ShortestGameTime%60),
			}...)
		}
		table = append(table, item)
	}
	return table
}

func (MultipleTournamentStats) playedTimeStats(data *model.EntityPlayer, timePlayed int) {
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
