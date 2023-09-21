package cmd

import (
	"fmt"
	"strconv"

	"github.com/crispgm/kicker-cli/pkg/rating"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
)

var evalRankSystems []string

func init() {
	evaluateRankCmd.Flags().StringArrayVarP(&evalRankSystems, "rank-system", "s", []string{rating.KLocal, rating.ATSA50, rating.ITSFProTour}, "rank systems")
	evaluateCmd.AddCommand(evaluateRankCmd)
}

var evaluateRankCmd = &cobra.Command{
	Use:   "rank",
	Short: "Evaluate estimated ranking points gained from event",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			pterm.Error.Println("Invalid params")
			cmd.Usage()
			return
		}
		rank := rating.Rank{}
		var table [][]string
		if !globalNoHeaders {
			table = append(table, []string{"Place", "Level", "Points Gained"})
		}
		for _, arg := range args {
			place, err := strconv.Atoi(arg)
			if err != nil {
				errorMessageAndExit(arg, "is not a valid place")
			}
			for _, level := range evalRankSystems {
				factors := rating.Factor{
					Place: place,
					Level: level,
				}
				gained := rank.Calculate(factors)
				table = append(table, []string{
					fmt.Sprintf("%d", place),
					level,
					fmt.Sprintf("+%d", int(gained)),
				})
			}
		}
		pterm.DefaultTable.WithHasHeader(!globalNoHeaders).WithData(table).WithBoxed(!globalNoBoxes).Render()
	},
}
