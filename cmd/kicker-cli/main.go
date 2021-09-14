package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/stat"
	monsterdyp "github.com/crispgm/kickertool-analyzer/stat/monster_dyp"
	"github.com/pterm/pterm"
)

// flags
var (
	mode   string
	player string
	files  []string

	dryRun bool

	// Options
	rankMinThreshold int
	withTime         bool
	withHostAway     bool
)

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry Run")
	flag.StringVar(&mode, "mode", "", "Stat mode. Supported: mts, mtt")
	flag.StringVar(&player, "player", "", "Players' data file")
	flag.IntVar(&rankMinThreshold, "rmt", 0, "Rank Minimum Threshold")
	flag.BoolVar(&withTime, "with-time", false, "With Time Analysis")
	flag.BoolVar(&withHostAway, "with-host-away", false, "With Host/Away Analysis")
	flag.Parse()

	// check mode
	if supported, ok := stat.SupportedStat[mode]; !ok || !supported {
		fmt.Println("Invalid mode", mode)
		os.Exit(1)
	}
	pterm.Info.Println("Statistics mode:", mode)

	// load players
	if len(player) == 0 {
		fmt.Println("No given player file")
		os.Exit(1)
	}
	pterm.Info.Println("Loading players ...")
	players, err := parsePlayer(player)
	if err != nil {
		pterm.Error.Println("Load players failed:", err)
	}

	// load tournaments
	files = flag.Args()
	if len(files) == 0 {
		fmt.Println("No given files")
		os.Exit(1)
	}
	pterm.Info.Println("Loading tournaments ...")
	var tournaments []model.Tournament

	// parsing
	p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithTitle("Processing tournaments data").Start()
	for _, fn := range files {
		pterm.Success.Println("Parsing", fn)
		t, err := parseTournament(fn)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		tournaments = append(tournaments, *t)
		p.Increment()
		time.Sleep(time.Millisecond * 100)
	}

	// calculating
	var statInfo stat.BaseStat
	option := stat.Option{
		RankMinThreshold: rankMinThreshold,
		WithTime:         withTime,
		WithHostAway:     withHostAway,
	}
	if mode == "mts" {
		statInfo = monsterdyp.NewMultipleTournamentStats(tournaments, players, option)
	} else if mode == "mtt" {
		statInfo = monsterdyp.NewMultipleTournamentTeamStats(tournaments, players, option)
	}
	if statInfo.ValidMode() {
		table := statInfo.Output()
		if !dryRun {
			pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
		}
	}
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

func parsePlayer(fn string) ([]model.EntityPlayer, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var players []model.EntityPlayer
	err = json.Unmarshal(data, &players)
	if err != nil {
		return nil, err
	}
	return players, err
}
