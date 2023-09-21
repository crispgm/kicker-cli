package cmd

import (
	"fmt"

	"github.com/crispgm/kicker-cli/pkg/rating"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var (
	evalEloKFactor int
	teamMode       bool
)

func init() {
	evaluateEloCmd.Flags().IntVarP(&evalEloKFactor, "elo-k", "k", 40, "K factor")
	evaluateCmd.AddCommand(evaluateEloCmd)
}

var evaluateEloCmd = &cobra.Command{
	Use:   "elo",
	Short: "Evaluate estimated ELO changes between two teams/players",
	Run: func(cmd *cobra.Command, args []string) {
		numOfPlayers := len(args)
		if numOfPlayers < 2 {
			pterm.Error.Println("Invalid params")
			cmd.Usage()
			return
		}
		if numOfPlayers >= 4 {
			teamMode = true
		}
		eloMain(args)
	},
}

func eloMain(args []string) {
	var (
		t1p1Score float64
		t1p2Score float64
		t2p1Score float64
		t2p2Score float64
	)
	if teamMode {
		t1p1Score = convertToFloat(args[0])
		t1p2Score = convertToFloat(args[1])
		t2p1Score = convertToFloat(args[2])
		t2p2Score = convertToFloat(args[3])
	} else {
		t1p1Score = convertToFloat(args[0])
		t1p2Score = convertToFloat(args[0])
		t2p1Score = convertToFloat(args[1])
		t2p2Score = convertToFloat(args[1])
	}
	pterm.Printfln("Estimated ELO rating (k=%d):", evalEloKFactor)
	fmt.Println()

	var (
		t1AvgScore float64
		t2AvgScore float64
	)
	t1AvgScore = (t1p1Score + t1p2Score) / 2
	t2AvgScore = (t2p1Score + t2p2Score) / 2
	er := rating.Elo{
		K: evalEloKFactor,
	}
	t1p1Win := er.Calculate(rating.Factor{
		PlayerScore:   t1p1Score,
		OpponentScore: t2AvgScore,
		Result:        rating.Win,
		Played:        0,
	})
	t1p1Loss := er.Calculate(rating.Factor{
		PlayerScore:   t1p1Score,
		OpponentScore: t2AvgScore,
		Result:        rating.Loss,
		Played:        0,
	})
	t1p2Win := er.Calculate(rating.Factor{
		PlayerScore:   t1p2Score,
		OpponentScore: t2AvgScore,
		Result:        rating.Win,
		Played:        0,
	})
	t1p2Loss := er.Calculate(rating.Factor{
		PlayerScore:   t1p2Score,
		OpponentScore: t2AvgScore,
		Result:        rating.Loss,
		Played:        0,
	})
	t2p1Win := er.Calculate(rating.Factor{
		PlayerScore:   t2p1Score,
		OpponentScore: t1AvgScore,
		Result:        rating.Win,
		Played:        0,
	})
	t2p1Loss := er.Calculate(rating.Factor{
		PlayerScore:   t2p1Score,
		OpponentScore: t1AvgScore,
		Result:        rating.Loss,
		Played:        0,
	})
	t2p2Win := er.Calculate(rating.Factor{
		PlayerScore:   t2p2Score,
		OpponentScore: t1AvgScore,
		Result:        rating.Win,
		Played:        0,
	})
	t2p2Loss := er.Calculate(rating.Factor{
		PlayerScore:   t2p2Score,
		OpponentScore: t1AvgScore,
		Result:        rating.Loss,
		Played:        0,
	})

	if teamMode {
		fmt.Println("- If team1 won:")
		hostWinData := [][]string{
			{"Team", "Player", "Cur Score", "New Score", "Rating Change"},
			{"A", "1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.f", t1p1Win), pterm.Green(fmt.Sprintf("+%.f", t1p1Win-t1p1Score))},
			{"A", "2", fmt.Sprintf("%.f", t1p2Score), fmt.Sprintf("%.f", t1p2Win), pterm.Green(fmt.Sprintf("+%.f", t1p2Win-t1p2Score))},
			{"B", "1", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.f", t2p1Loss), pterm.Red(fmt.Sprintf("%.f", t2p1Loss-t2p1Score))},
			{"B", "2", fmt.Sprintf("%.f", t2p2Score), fmt.Sprintf("%.f", t2p2Loss), pterm.Red(fmt.Sprintf("%.f", t2p2Loss-t2p2Score))},
		}
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(hostWinData).WithBoxed(!globalNoBoxes).Render()
		fmt.Println()

		fmt.Println("- If team2 won:")
		awayWinData := [][]string{
			{"Team", "Player", "Cur Score", "New Score", "Rating Change"},
			{"A", "1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.f", t1p1Loss), pterm.Red(fmt.Sprintf("%.f", t1p1Loss-t1p1Score))},
			{"A", "2", fmt.Sprintf("%.f", t1p2Score), fmt.Sprintf("%.f", t1p2Loss), pterm.Red(fmt.Sprintf("%.f", t1p2Loss-t1p2Score))},
			{"B", "1", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.f", t2p1Win), pterm.Green(fmt.Sprintf("+%.f", t2p1Win-t2p1Score))},
			{"B", "2", fmt.Sprintf("%.f", t2p2Score), fmt.Sprintf("%.f", t2p2Win), pterm.Green(fmt.Sprintf("+%.f", t2p2Win-t2p2Score))},
		}
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(awayWinData).WithBoxed(!globalNoBoxes).Render()
	} else {
		fmt.Println("- If player1 won:")
		hostWinData := [][]string{
			{"Player", "Cur Score", "New Score", "Rating Change"},
			{"1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.f", t1p1Win), pterm.Green(fmt.Sprintf("+%.f", t1p1Win-t1p1Score))},
			{"2", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.f", t2p1Loss), pterm.Red(fmt.Sprintf("%.f", t2p1Loss-t2p1Score))},
		}
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(hostWinData).WithBoxed(!globalNoBoxes).Render()
		fmt.Println()

		fmt.Println("- If player2 won:")
		awayWinData := [][]string{
			{"Player", "Cur Score", "New Score", "Rating Change"},
			{"1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.f", t1p1Loss), pterm.Red(fmt.Sprintf("%.f", t1p1Loss-t1p1Score))},
			{"2", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.f", t2p1Win), pterm.Green(fmt.Sprintf("+%.f", t2p1Win-t2p1Score))},
		}
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(awayWinData).WithBoxed(!globalNoBoxes).Render()
	}
}
