package main

import (
	"flag"
	"os"
	"time"

	"github.com/crispgm/kickertool-analyzer/model"
	"github.com/crispgm/kickertool-analyzer/operator"
	monsterdyp "github.com/crispgm/kickertool-analyzer/operator/monster_dyp"
	"github.com/crispgm/kickertool-analyzer/parser"
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
	withPoint        bool
)

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry Run")
	flag.StringVar(&mode, "mode", "", "Stat mode. Supported: mdp, mdt")
	flag.StringVar(&player, "player", "", "Players' data file")
	flag.IntVar(&rankMinThreshold, "rmt", 0, "Rank Minimum Threshold")
	flag.BoolVar(&withTime, "with-time", false, "With Time Analysis")
	flag.BoolVar(&withHomeAway, "with-home-away", false, "With Home/Away Analysis")
	flag.BoolVar(&withPoint, "with-point", false, "With Point Analysis")
	flag.Parse()

	// check mode
	if supported, ok := operator.SupportedOperator[mode]; !ok || !supported {
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
		pterm.Info.Println("Parsing", fn)
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
	var statOperator operator.BaseOperator
	option := operator.Option{
		RankMinThreshold: rankMinThreshold,
		WithTime:         withTime,
		WithHomeAway:     withHomeAway,
		WithPoint:        withPoint,
	}
	if mode == model.ModeMonsterDYPPlayerStats {
		statOperator = monsterdyp.NewPlayerStats(games, option)
	} else if mode == model.ModeMonsterDYPTeamStats {
		statOperator = monsterdyp.NewTeamStats(games, option)
	}
	valid := true
	for _, t := range tournaments {
		if t.Mode != model.ModeMonsterDYP {
			valid = false
			break
		}
	}
	if valid {
		pterm.Info.Println("Briefing:", c.Briefing())
		table := statOperator.Output()
		if !dryRun {
			pterm.DefaultTable.WithHasHeader().WithData(table).WithBoxed(true).Render()
		}
	} else {
		pterm.Error.Println("Unsupported tournament mode for this operator")
	}
}
