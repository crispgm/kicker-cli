package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

var showScore bool

func init() {
	eventInfoCmd.Flags().BoolVarP(&showScore, "show-score", "s", false, "show scores of each set")
	eventCmd.AddCommand(eventInfoCmd)
}

var eventInfoCmd = &cobra.Command{
	Use:     "info",
	Aliases: []string{"show"},
	Short:   "Show event details",
	Long:    "Show event details",
	Run: func(cmd *cobra.Command, args []string) {
		var arg string
		if isInputFromPipe() {
			scanner := bufio.NewScanner(bufio.NewReader(bufio.NewReader(os.Stdin)))
			for scanner.Scan() {
				arg = strings.Trim(scanner.Text(), " \t\n")
				break
			}
			if len(arg) == 0 {
				errorMessageAndExit("Please present an event ID")
			}
		} else {
			if len(args) == 0 {
				errorMessageAndExit("Please present an event ID")
			}
			arg = args[0]
		}
		instance := initInstanceAndLoadConf()
		e := instance.GetEvent(arg)
		if e == nil {
			errorMessageAndExit("No event(s) found")
		}
		table := initEventInfoHeader()
		t, r := loadAndShowEventInfo(&table, instance.DataPath(), instance.Conf.Players, e, instance.Conf.Organization.Timezone)
		_ = pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
		table = showGames(r.PreliminaryRounds, t.Options)
		if len(table) > 0 {
			pterm.Println("Rounds:")
			_ = pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
		sort.SliceStable(r.LoserBracket, func(i, j int) bool {
			return true
		})
		table = showGames(r.LoserBracket, t.Options)
		if len(table) > 0 {
			pterm.Println("Loser Bracket:")
			_ = pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
		table = showGames(r.WinnerBracket, t.Options)
		if len(table) > 0 {
			pterm.Println("Winner Bracket:")
			_ = pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
		// show result
		table = showResults(r.Ranks, t.IsSingle())
		if len(table) > 0 {
			pterm.Println("Result:")
			_ = pterm.DefaultTable.WithHasHeader(false).WithData(table).WithBoxed(!globalNoBoxes).Render()
		}
	},
}

func showGames(games []entity.Game, options model.Options) [][]string {
	var table [][]string
	if len(games) == 0 {
		return table
	}
	numOfSets := len(games[0].Sets)
	if numOfSets == 1 && options.FastInput {
		showScore = false
	}

	if showScore {
		for _, g := range games {
			numOfSets := len(g.Sets)
			if numOfSets > 1 {
				table = append(table, []string{
					g.Name,
					g.Team1[0],
					fmt.Sprintf("%d:%d", g.Point1, g.Point2),
					g.Team2[0],
				})
				for _, s := range g.Sets {
					if len(g.Team1) == 1 {
						table = append(table, []string{
							"",
							"",
							fmt.Sprintf("%d:%d", s.Point1, s.Point2),
							"",
						})
					} else {
						table = append(table, []string{
							"",
							"",
							fmt.Sprintf("%d:%d", s.Point1, s.Point2),
							"",
						})
					}
				}
			} else {
				for _, s := range g.Sets {
					if len(g.Team1) == 1 {
						table = append(table, []string{
							g.Name,
							g.Team1[0],
							fmt.Sprintf("%d:%d", s.Point1, s.Point2),
							g.Team2[0],
						})
					} else {
						table = append(table, []string{
							g.Name,
							fmt.Sprintf("%s/%s", g.Team1[0], g.Team1[1]),
							fmt.Sprintf("%d:%d", s.Point1, s.Point2),
							fmt.Sprintf("%s/%s", g.Team2[0], g.Team2[1]),
						})
					}
				}
			}
		}
	} else {
		for _, g := range games {
			if len(g.Team1) == 1 {
				table = append(table, []string{
					g.Name,
					g.Team1[0],
					fmt.Sprintf("%d:%d", g.Point1, g.Point2),
					g.Team2[0],
				})
			} else {
				table = append(table, []string{
					g.Name,
					fmt.Sprintf("%s/%s", g.Team1[0], g.Team1[1]),
					fmt.Sprintf("%d:%d", g.Point1, g.Point2),
					fmt.Sprintf("%s/%s", g.Team2[0], g.Team2[1]),
				})
			}
		}
	}
	return table
}

func showResults(ranks [][]entity.Player, single bool) [][]string {
	var table [][]string
	var rank = 1
	for i := len(ranks) - 1; i >= 0; i-- {
		r := ranks[i]
		var level []string
		if rank == 1 {
			level = append(level, "Champion")
		} else if rank == 2 {
			level = append(level, "Runner-Up")
		} else {
			level = append(level, fmt.Sprintf("%d", rank))
		}
		if single {
			rank += len(r)
			for _, p := range r {
				level = append(level, p.Name)
			}
		} else {
			rank += len(r) / 2
			for i := 0; i < len(r); i += 2 {
				p1 := r[i]
				p2 := r[i+1]
				level = append(level, fmt.Sprintf("%s/%s", p1.Name, p2.Name))
			}
		}
		table = append(table, level)
	}
	return table
}
