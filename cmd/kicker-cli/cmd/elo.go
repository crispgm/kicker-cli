package cmd

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/pkg/rating"
	"github.com/crispgm/kicker-cli/pkg/rating/elo"
)

var (
	evalAlgo       string
	evalEloKFactor int
	teamMode       bool
)

func init() {
	evaluateCmd.Flags().StringVarP(&evalAlgo, "algorithm", "a", "elo", "rating algorithm")
	evaluateCmd.Flags().IntVarP(&evalEloKFactor, "elo-k", "k", elo.K, "K factor")
	rootCmd.AddCommand(evaluateCmd)
}

var evaluateCmd = &cobra.Command{
	Use:     "evaluate",
	Aliases: []string{"eval"},
	Short:   "Simple tool to evaluate estimated changes between two teams/players.",
	Long: `Simple tool to evaluate estimated changes between two teams/players.
$ kicker-cli evaluate -a elo 1100 1200
$ kicker-cli evaluate -a elo 1103 1203 1289 1013
$ kicker-cli evaluate -a elo -k 20 1103 1203 1289 1013`,
	Run: func(cmd *cobra.Command, args []string) {
		if evalAlgo != "elo" {
			errorMessageAndExit("invalid algorithm")
		}
		numOfPlayers := len(args)
		if numOfPlayers < 2 {
			pterm.Error.Println("Invalid params")
			cmd.Usage()
			return
		}
		if numOfPlayers >= 4 {
			teamMode = true
		}
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
		pterm.Printfln("Estimated %s rating (k=%d):", evalAlgo, evalEloKFactor)
		fmt.Println()
		eloMain(t1p1Score, t1p2Score, t2p1Score, t2p2Score)
	},
}

func eloMain(t1p1Score, t1p2Score, t2p1Score, t2p2Score float64) {
	var (
		t1AvgScore float64
		t2AvgScore float64
	)
	t1AvgScore = (t1p1Score + t1p2Score) / 2
	t2AvgScore = (t2p1Score + t2p2Score) / 2
	er := elo.Elo{K: float64(evalEloKFactor)}
	er.InitialScore(t1p1Score, t2AvgScore)
	t1p1Win := er.Calculate(rating.Win)
	er.InitialScore(t1p2Score, t2AvgScore)
	t1p1Loss := er.Calculate(rating.Loss)
	er.InitialScore(t1p2Score, t2AvgScore)
	t1p2Win := er.Calculate(rating.Win)
	er.InitialScore(t1p2Score, t2AvgScore)
	t1p2Loss := er.Calculate(rating.Loss)
	er.InitialScore(t2p1Score, t1AvgScore)
	t2p1Win := er.Calculate(rating.Win)
	er.InitialScore(t2p1Score, t1AvgScore)
	t2p1Loss := er.Calculate(rating.Loss)
	er.InitialScore(t2p2Score, t1AvgScore)
	t2p2Win := er.Calculate(rating.Win)
	er.InitialScore(t2p2Score, t1AvgScore)
	t2p2Loss := er.Calculate(rating.Loss)

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
