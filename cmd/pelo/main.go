package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/crispgm/kickertool-analyzer/elo"
	"github.com/pterm/pterm"
)

func main() {
	if len(os.Args) != 5 {
		pterm.Error.Println("Invalid params")
	}

	t1p1Score := convertToFloat(os.Args[1])
	t1p2Score := convertToFloat(os.Args[2])
	t2p1Score := convertToFloat(os.Args[3])
	t2p2Score := convertToFloat(os.Args[4])
	rhw := elo.Rate{
		T1P1Score: t1p1Score,
		T1P2Score: t1p2Score,
		T2P1Score: t2p1Score,
		T2P2Score: t2p2Score,
		HostWin:   true,
	}
	rhw.CalcEloRating()
	raw := elo.Rate{
		T1P1Score: t1p1Score,
		T1P2Score: t1p2Score,
		T2P1Score: t2p1Score,
		T2P2Score: t2p2Score,
		HostWin:   false,
	}
	raw.CalcEloRating()

	fmt.Println("Estimated Elo rating:")
	fmt.Println("- If host won:")
	fmt.Println("  > [Team 1] Player 1:", t1p1Score, "->", rhw.T1P1Score, pterm.Green("+", rhw.T1P1Score-t1p1Score))
	fmt.Println("  > [Team 1] Player 2:", t1p2Score, "->", rhw.T1P2Score, pterm.Green("+", rhw.T1P2Score-t1p2Score))
	fmt.Println("  > [Team 2] Player 3:", t2p1Score, "->", rhw.T2P1Score, pterm.Red(rhw.T2P1Score-t2p1Score))
	fmt.Println("  > [Team 2] Player 4:", t2p2Score, "->", rhw.T2P2Score, pterm.Red(rhw.T2P2Score-t2p2Score))
	fmt.Println("- If away won:")
	fmt.Println("  > [Team 1] Player 1:", t1p1Score, "->", raw.T1P1Score, pterm.Red(raw.T1P1Score-t1p1Score))
	fmt.Println("  > [Team 1] Player 2:", t1p2Score, "->", raw.T1P2Score, pterm.Red(raw.T1P2Score-t1p2Score))
	fmt.Println("  > [Team 2] Player 3:", t2p1Score, "->", raw.T2P1Score, pterm.Green("+", raw.T2P1Score-t2p1Score))
	fmt.Println("  > [Team 2] Player 4:", t2p2Score, "->", raw.T2P2Score, pterm.Green("+", raw.T2P2Score-t2p2Score))
}

func convertToFloat(in string) float64 {
	out, _ := strconv.Atoi(in)
	return float64(out)
}
