package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/crispgm/kicker-cli/pkg/elo"
)

func init() {
	eloCmd.Flags().IntVarP(&eloKFactor, "elo-k", "k", elo.K, "K factor")
	rootCmd.AddCommand(eloCmd)
}

var (
	eloKFactor int
	teamMode   bool
)

var eloCmd = &cobra.Command{
	Use:   "elo",
	Short: "Simple tool to show estimated ELO changes between two teams/players.",
	Long: `Simple tool to show estimated ELO changes between two teams/players.
$ pelo 1100 1200
$ pelo 1103 1203 1289 1013
$ pelo -k 20 1103 1203 1289 1013`,
	Run: eloMain,
}

func eloMain(cmd *cobra.Command, args []string) {
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
		t1p1Score = convertToFloat(os.Args[1])
		t1p2Score = convertToFloat(os.Args[2])
		t2p1Score = convertToFloat(os.Args[3])
		t2p2Score = convertToFloat(os.Args[4])
	} else {
		t1p1Score = convertToFloat(os.Args[1])
		t1p2Score = convertToFloat(os.Args[1])
		t2p1Score = convertToFloat(os.Args[2])
		t2p2Score = convertToFloat(os.Args[2])
	}
	rhw := elo.Rate{
		T1P1Score: t1p1Score,
		T1P2Score: t1p2Score,
		T2P1Score: t2p1Score,
		T2P2Score: t2p2Score,
		HostWin:   true,
		K:         float64(eloKFactor),
	}
	rhw.CalcEloRating()
	raw := elo.Rate{
		T1P1Score: t1p1Score,
		T1P2Score: t1p2Score,
		T2P1Score: t2p1Score,
		T2P2Score: t2p2Score,
		HostWin:   false,
		K:         float64(eloKFactor),
	}
	raw.CalcEloRating()

	pterm.Info.Printf("Estimated Elo rating (k=%d):\n", eloKFactor)

	if teamMode {
		fmt.Println("- If host won:")
		hostWonData := [][]string{
			{"Team", "Player", "Cur Score", "Expectation", "New Score", "Rating Change"},
			{"A", "1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.2f%%", rhw.T1P1Exp*100), fmt.Sprintf("%.f", rhw.T1P1Score), pterm.Green("+", rhw.T1P1Score-t1p1Score)},
			{"A", "2", fmt.Sprintf("%.f", t1p2Score), fmt.Sprintf("%.2f%%", rhw.T1P2Exp*100), fmt.Sprintf("%.f", rhw.T1P2Score), pterm.Green("+", rhw.T1P2Score-t1p2Score)},
			{"B", "1", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.2f%%", rhw.T2P1Exp*100), fmt.Sprintf("%.f", rhw.T2P1Score), pterm.Red(rhw.T2P1Score - t2p1Score)},
			{"B", "2", fmt.Sprintf("%.f", t2p2Score), fmt.Sprintf("%.2f%%", rhw.T2P2Exp*100), fmt.Sprintf("%.f", rhw.T2P2Score), pterm.Red(rhw.T2P2Score - t2p2Score)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(hostWonData).WithBoxed(true).Render()

		fmt.Println("- If away won:")
		awayWonData := [][]string{
			{"Team", "Player", "Cur Score", "Expect", "New Score", "Rating Change"},
			{"A", "1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.2f%%", raw.T1P1Exp*100), fmt.Sprintf("%.f", raw.T1P1Score), pterm.Red(raw.T1P1Score - t1p1Score)},
			{"A", "2", fmt.Sprintf("%.f", t1p2Score), fmt.Sprintf("%.2f%%", raw.T1P2Exp*100), fmt.Sprintf("%.f", raw.T1P2Score), pterm.Red(raw.T1P2Score - t1p2Score)},
			{"B", "1", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.2f%%", raw.T2P1Exp*100), fmt.Sprintf("%.f", raw.T2P1Score), pterm.Green("+", raw.T2P1Score-t2p1Score)},
			{"B", "2", fmt.Sprintf("%.f", t2p2Score), fmt.Sprintf("%.2f%%", raw.T2P2Exp*100), fmt.Sprintf("%.f", raw.T2P2Score), pterm.Green("+", raw.T2P2Score-t2p2Score)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(awayWonData).WithBoxed(true).Render()
	} else {
		fmt.Println("- If host won:")
		hostWonData := [][]string{
			{"Player", "Cur Score", "Expectation", "New Score", "Rating Change"},
			{"1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.2f%%", rhw.T1P1Exp*100), fmt.Sprintf("%.f", rhw.T1P1Score), pterm.Green("+", rhw.T1P1Score-t1p1Score)},
			{"2", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.2f%%", rhw.T2P1Exp*100), fmt.Sprintf("%.f", rhw.T2P1Score), pterm.Red(rhw.T2P1Score - t2p1Score)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(hostWonData).WithBoxed(true).Render()

		fmt.Println("- If away won:")
		awayWonData := [][]string{
			{"Player", "Cur Score", "Expect", "New Score", "Rating Change"},
			{"1", fmt.Sprintf("%.f", t1p1Score), fmt.Sprintf("%.2f%%", raw.T1P1Exp*100), fmt.Sprintf("%.f", raw.T1P1Score), pterm.Red(raw.T1P1Score - t1p1Score)},
			{"2", fmt.Sprintf("%.f", t2p1Score), fmt.Sprintf("%.2f%%", raw.T2P1Exp*100), fmt.Sprintf("%.f", raw.T2P1Score), pterm.Green("+", raw.T2P1Score-t2p1Score)},
		}
		pterm.DefaultTable.WithHasHeader().WithData(awayWonData).WithBoxed(true).Render()
	}
}

func convertToFloat(in string) float64 {
	out, _ := strconv.Atoi(in)
	return float64(out)
}
