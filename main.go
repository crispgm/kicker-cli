package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"

	"github.com/crispgm/kickertool-analyzer/model"
)

func main() {
	argc := len(os.Args)
	if argc <= 1 {
		os.Exit(1)
	}
	var tournaments []model.Tournament
	for _, fn := range os.Args[1:] {
		fmt.Println("Parsing tournaments data:", fn)
		t, err := parseTournament(fn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tournaments = append(tournaments, *t)
	}

	fmt.Println()
	playerStats(tournaments)
}

func parseTournament(fn string) (*model.Tournament, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var tournament model.Tournament
	err = json.Unmarshal(data, &tournament)
	if err != nil {
		return nil, err
	}
	return &tournament, err
}

func playerStats(tournaments []model.Tournament) {
	var data = make(map[string]model.EntityPlayer)
	for _, t := range tournaments {
		var teams = make(map[string]model.Team)
		var players = make(map[string]model.Player)
		for _, p := range t.Players {
			if !p.Removed {
				var found bool
				for _, ep := range model.AllPlayers {
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
							} else if s.Team2 > s.Team1 {
								t1p1Data.Lost++
								t1p2Data.Lost++
								t2p1Data.Won++
								t2p2Data.Won++
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
	fmt.Println("Name\t\tNum\tWon\tLost\tG+\tG-\tG±\tPPG\tWR")
	for _, d := range sliceData {
		if len(d.Name) < 8 {
			d.Name = d.Name + "        "
		}
		fmt.Printf("%s\t%d\t%d\t%d\t%d\t%d\t%d\t%.2f\t%.0f%%\n",
			d.Name,
			d.Played, d.Won, d.Lost,
			d.Goals, d.GoalsIn, d.GoalDiff,
			d.PointsPerGame, d.WinRate)
	}
}
