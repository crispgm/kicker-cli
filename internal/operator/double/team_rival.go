package double

import (
	"fmt"
	"sort"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/internal/operator"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var _ operator.Operator = (*TeamRival)(nil)

// TeamRival generate statistics data of multiple monster DYP tournaments by team
type TeamRival struct {
	options operator.Option
	games   []entity.Game
}

// SupportedFormats .
func (t TeamRival) SupportedFormats(trn *model.Tournament) bool {
	if trn.IsDouble() {
		if trn.Mode == model.ModeMonsterDYP ||
			trn.Mode == model.ModeSwissSystem || trn.Mode == model.ModeRounds || trn.Mode == model.ModeRoundRobin ||
			trn.Mode == model.ModeDoubleElimination || trn.Mode == model.ModeElimination {
			return true
		}
	}

	return false
}

// Input .
func (t *TeamRival) Input(games []entity.Game, players []entity.Player, options operator.Option) {
	t.games = games
	t.options = options
}

// Output .
func (t *TeamRival) Output() [][]string {
	data := make(map[string]entity.Rival)
	for _, g := range t.games {
		t1p1Name := g.Team1[0]
		t1p2Name := g.Team1[1]
		t2p1Name := g.Team2[0]
		t2p2Name := g.Team2[1]
		team1Name := fmt.Sprintf("%s/%s", t1p1Name, t1p2Name)
		if t1p1Name > t1p2Name {
			team1Name = fmt.Sprintf("%s/%s", t1p2Name, t1p1Name)
		}
		team2Name := fmt.Sprintf("%s/%s", t2p1Name, t2p2Name)
		if t2p1Name > t2p2Name {
			team2Name = fmt.Sprintf("%s/%s", t2p2Name, t2p1Name)
		}

		rivalName := fmt.Sprintf("%s_vs_%s", team1Name, team2Name)
		rivalNameAlt := fmt.Sprintf("%s_vs_%s", team2Name, team1Name)

		var rival entity.Rival
		if _, ok := data[rivalName]; ok {
			rival = data[rivalName]
		} else if _, ok := data[rivalNameAlt]; ok {
			rivalName = rivalNameAlt
			rival = data[rivalNameAlt]
		} else {
			rival = entity.Rival{
				Team1: entity.Team{
					Player1: t1p1Name,
					Player2: t1p2Name,
				},
				Team2: entity.Team{
					Player1: t2p1Name,
					Player2: t2p2Name,
				},
			}
		}

		rival.Played++
		rival.Team1.Played++
		rival.Team2.Played++
		timePlayed := g.TimePlayed
		rival.TimePlayed += timePlayed
		rival.Team1.TimePlayed += timePlayed
		rival.Team2.TimePlayed += timePlayed

		if g.Point1 > g.Point2 {
			rival.Win++
			rival.Team1.Win++
			rival.Team2.Loss++
			rival.Team1.GoalsWin += (g.Point1 - g.Point2)
			rival.Team2.GoalsInLoss += (g.Point1 - g.Point2)
		} else if g.Point1 < g.Point2 {
			rival.Loss++
			rival.Team2.Win++
			rival.Team1.Loss++
			rival.Team2.GoalsWin += (g.Point2 - g.Point1)
			rival.Team1.GoalsInLoss += (g.Point2 - g.Point1)
		} else {
			rival.Draw++
			rival.Team1.Draw++
			rival.Team2.Draw++
		}
		rival.Team1.Goals += g.Point1
		rival.Team2.Goals += g.Point2
		rival.Team1.GoalsIn += g.Point2
		rival.Team2.GoalsIn += g.Point1

		data[rivalName] = rival
	}

	var sliceData []entity.Rival
	for _, d := range data {
		sliceData = append(sliceData, d)
	}

	sort.SliceStable(sliceData, func(i, j int) bool {
		if sliceData[i].Played >= t.options.MinimumPlayed && sliceData[j].Played < t.options.MinimumPlayed {
			return true
		}
		if sliceData[i].Played < t.options.MinimumPlayed && sliceData[j].Played >= t.options.MinimumPlayed {
			return false
		}

		if sliceData[i].Played > sliceData[j].Played {
			return true
		}
		return false
	})

	if t.options.Head > 0 && len(sliceData) > t.options.Head {
		sliceData = sliceData[:t.options.Head]
	} else if t.options.Tail > 0 && len(sliceData) > t.options.Tail {
		sliceData = sliceData[len(sliceData)-t.options.Tail:]
	}

	header := []string{"#", "Team1", "Team2", "Num", "Win", "Loss", "Draw"}
	table := [][]string{}
	if t.options.WithHeader {
		table = append(table, header)
	}
	for i, d := range sliceData {
		if d.Played == 0 {
			continue
		}
		item := []string{
			fmt.Sprintf("%d", i+1),
			fmt.Sprintf("%s/%s", d.Team1.Player1, d.Team1.Player2),
			fmt.Sprintf("%s/%s", d.Team2.Player1, d.Team2.Player2),
			fmt.Sprintf("%d", d.Played),
			fmt.Sprintf("%d", d.Win),
			fmt.Sprintf("%d", d.Loss),
			fmt.Sprintf("%d", d.Draw),
		}
		table = append(table, item)
	}
	return table
}
