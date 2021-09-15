package main

import (
	"flag"
	"os"
	"time"

	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/parser"
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
	withHomeAway     bool
)

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry Run")
	flag.StringVar(&mode, "mode", "", "Stat mode. Supported: mts, mtt")
	flag.StringVar(&player, "player", "", "Players' data file")
	flag.IntVar(&rankMinThreshold, "rmt", 0, "Rank Minimum Threshold")
	flag.BoolVar(&withTime, "with-time", false, "With Time Analysis")
	flag.BoolVar(&withHomeAway, "with-home-away", false, "With Home/Away Analysis")
	flag.Parse()

	// check mode
	if supported, ok := stat.SupportedStat[mode]; !ok || !supported {
		pterm.Error.Println("Invalid mode", mode)
		os.Exit(1)
	}
	pterm.Info.Println("Statistics mode:", mode)

	// load players
	if len(player) == 0 {
		pterm.Error.Println("No given player file")
		os.Exit(1)
	}
	pterm.Info.Println("Loading players ...")
	players, err := parser.ParsePlayer(player)
	if err != nil {
		pterm.Error.Println("Load players failed:", err)
		os.Exit(1)
	}

	// load tournaments
	files = flag.Args()
	if len(files) == 0 {
		pterm.Error.Println("No given files")
		os.Exit(1)
	}
	pterm.Info.Println("Loading tournaments ...")
	var tournaments []model.Tournament

	// parsing
	p, _ := pterm.DefaultProgressbar.WithTotal(len(files)).WithRemoveWhenDone().WithTitle("Processing tournaments data").Start()
	for _, fn := range files {
		pterm.Success.Println("Parsing", fn)
		t, err := parser.ParseTournament(fn)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		tournaments = append(tournaments, *t)
		p.Increment()
		time.Sleep(time.Millisecond * 100)
	}
	c := parser.NewConverter()
	games, err := c.Normalize(tournaments, players)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}

	// calculating
	var statInfo stat.BaseStat
	option := stat.Option{
		RankMinThreshold: rankMinThreshold,
		WithTime:         withTime,
		WithHomeAway:     withHomeAway,
	}
	if mode == "mts" {
		statInfo = monsterdyp.NewMultipleTournamentStats(games, option)
	} else if mode == "mtt" {
		statInfo = monsterdyp.NewMultipleTournamentTeamStats(games, option)
	}
	valid := true
	for _, t := range tournaments {
		if t.Mode != model.ModeMonsterDYP {
			valid = false
			break
		}
	}
	if valid {
		table := statInfo.Output()
		if !dryRun {
			pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
		}
	}
}
