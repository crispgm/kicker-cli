package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/crispgm/kickertool-analyzer/elo"
	"github.com/pterm/pterm"
)

const initialElo = 1500

// Team .
type Team struct {
	Player1 string
	Player2 string
}

// Match .
type Match struct {
	Team1   Team
	Team2   Team
	HomeWin bool
	Type    string
}

// Player .
type Player struct {
	Name    string
	Won     int
	Lost    int
	WinRate float64
	EloRank float64
}

func main() {
	var matches []Match
	for _, fn := range os.Args[1:] {
		pterm.Info.Println("Parsing", fn)
		m, err := parseFile(fn)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		matches = append(matches, m...)
	}

	mp := make(map[string]*Player)
	for _, m := range matches {
		t1p1 := getOrCreate(mp, m.Team1.Player1)
		t1p2 := getOrCreate(mp, m.Team1.Player2)
		t2p1 := getOrCreate(mp, m.Team2.Player1)
		t2p2 := getOrCreate(mp, m.Team2.Player2)
		if m.HomeWin {
			t1p1.Won++
			t1p2.Won++
			t2p1.Lost++
			t2p2.Lost++
		} else {
			t1p1.Lost++
			t1p2.Lost++
			t2p1.Won++
			t2p2.Won++
		}
		k := 20.0
		if m.Type == "DYP" {
			k = 30
		}
		rhw := elo.Rate{
			T1P1Score: t1p1.EloRank,
			T1P2Score: t1p2.EloRank,
			T2P1Score: t2p1.EloRank,
			T2P2Score: t2p2.EloRank,
			HostWin:   m.HomeWin,
			K:         k,
		}
		rhw.CalcEloRating()
		t1p1.EloRank = rhw.T1P1Score
		t1p2.EloRank = rhw.T1P2Score
		t2p1.EloRank = rhw.T2P1Score
		t2p2.EloRank = rhw.T2P2Score
	}
	var playerRank []*Player
	for _, p := range mp {
		p.WinRate = float64(p.Won) / float64(p.Won+p.Lost)
		playerRank = append(playerRank, p)
	}
	sort.SliceStable(playerRank, func(i, j int) bool {
		return playerRank[i].EloRank > playerRank[j].EloRank
	})
	output := [][]string{
		{"Player", "Won", "Lost", "Win%", "Elo"},
	}
	for _, p := range playerRank {
		output = append(output, []string{p.Name, fmt.Sprintf("%d", p.Won), fmt.Sprintf("%d", p.Lost), fmt.Sprintf("%.2f%%", p.WinRate*100), fmt.Sprintf("%.0f", p.EloRank)})
	}
	pterm.DefaultTable.WithHasHeader().WithData(output).WithBoxed(true).Render()
}

func getOrCreate(mp map[string]*Player, pn string) *Player {
	if p, ok := mp[pn]; ok {
		return p
	}
	mp[pn] = &Player{
		Name:    pn,
		EloRank: initialElo,
	}

	return mp[pn]
}

func parseFile(fn string) ([]Match, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	sData := string(data)
	lines := strings.Split(sData, "\n")
	matchType := "open double"
	if strings.Contains(fn, "DYP") {
		matchType = "dyp"
	}
	var matches []Match
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		cols := strings.Split(line, ",")
		t1 := cols[0]
		t2 := cols[1]
		res := cols[2]
		homeWin := true
		if res == "L" {
			homeWin = false
		}
		t1ps := strings.Split(t1, "/")
		t2ps := strings.Split(t2, "/")
		m := Match{
			Team1: Team{
				Player1: t1ps[0],
				Player2: t1ps[1],
			},
			Team2: Team{
				Player1: t2ps[0],
				Player2: t2ps[1],
			},
			HomeWin: homeWin,
			Type:    matchType,
		}
		matches = append(matches, m)
	}

	return matches, nil
}
