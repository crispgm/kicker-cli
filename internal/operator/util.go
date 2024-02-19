package operator

import (
	"encoding/csv"
	"encoding/json"
	"os"

	"github.com/pterm/pterm"

	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/rating"
)

// calculateELO calculate ELO for player
func calculateELO(played int, p1Elo, p2Elo float64, result int) float64 {
	eloCalc := rating.Elo{}
	factors := rating.Factor{
		Played:        played,
		PlayerScore:   p1Elo,
		OpponentScore: p2Elo,
		Result:        result,
	}
	return eloCalc.Calculate(factors)
}

func openSingleTournament(trn *model.Tournament) bool {
	if trn.IsSingle() {
		if trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}

func openDoubleTournament(trn *model.Tournament) bool {
	if trn.IsDouble() {
		if trn.Mode == model.ModeMonsterDYP ||
			trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}

func output(opt Option, header []string, body [][]string) {
	var table [][]string
	if opt.WithHeader {
		table = append(table, header)
		table = append(table, body...)
	}

	if opt.OutputFormat == "default" {
		_ = pterm.DefaultTable.WithHasHeader(opt.WithHeader).WithData(table).WithBoxed(opt.WithBoxes).Render()
	} else if opt.OutputFormat == "csv" {
		csvwriter := csv.NewWriter(os.Stdout)
		for _, row := range table {
			_ = csvwriter.Write(row)
		}
		csvwriter.Flush()
	} else if opt.OutputFormat == "tsv" {
		csvwriter := csv.NewWriter(os.Stdout)
		csvwriter.Comma = '\t'
		for _, row := range table {
			_ = csvwriter.Write(row)
		}
		csvwriter.Flush()
	} else if opt.OutputFormat == "json" {
		var jsonData []map[string]string
		for _, row := range body {
			jsonItem := make(map[string]string)
			for i, item := range row {
				key := header[i]
				if key == "#" {
					key = "Index"
				} else if key == "WR%" {
					key = "WinRate"
				}
				jsonItem[key] = item
			}
			jsonData = append(jsonData, jsonItem)
		}
		encoder := json.NewEncoder(os.Stdout)
		_ = encoder.Encode(jsonData)
	}
}
