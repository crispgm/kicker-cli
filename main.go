package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/crispgm/kickertool-analyzer/model"
	monsterdyb "github.com/crispgm/kickertool-analyzer/stat/monster_dyb"
	"github.com/olekukonko/tablewriter"
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
	data := monsterdyb.MultipleTournamentStats(tournaments)
	outputTable(data)
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

func outputTable(data []model.EntityPlayer) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Num", "Won", "Lost", "G+", "G-", "G±", "WR%", "PPG", "TPG"})
	for _, d := range data {
		table.Append([]string{
			d.Name,
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Won),
			fmt.Sprintf("%d", d.Lost),
			fmt.Sprintf("%d", d.Goals),
			fmt.Sprintf("%d", d.GoalsIn),
			fmt.Sprintf("%d", d.GoalDiff),
			fmt.Sprintf("%.0f%%", d.WinRate),
			fmt.Sprintf("%.2f", d.PointsPerGame),
			fmt.Sprintf("%02d:%02d", d.TimePerGame/60, d.TimePerGame%60),
		})
	}

	table.Render()
}
