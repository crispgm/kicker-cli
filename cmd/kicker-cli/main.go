package main

import (
	"flag"
	"os"
	"strings"

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
	orderBy          string
	rankMinThreshold int
	withTime         bool
	withHomeAway     bool
	withPoint        bool
	incremental      bool
)

func main() {
	flag.BoolVar(&dryRun, "dry-run", false, "Dry Run")
	flag.StringVar(&mode, "mode", "", "Stat mode. Supported: mdp, mdt")
	flag.StringVar(&player, "player", "", "Players' data file")
	flag.StringVar(&orderBy, "order-by", "wr", "Order by `wr` (win rate) or `elo` (ELO ranking)")
	flag.IntVar(&rankMinThreshold, "rmt", 0, "Rank minimum threshold")
	flag.BoolVar(&withTime, "with-time", false, "With time analysis")
	flag.BoolVar(&withHomeAway, "with-home-away", false, "With home/away analysis")
	flag.BoolVar(&withPoint, "with-point", false, "With point analysis")
	flag.BoolVar(&incremental, "incremental", false, "Update player's data incrementally")
	flag.Parse()

	// check mode
	if !operator.IsSupported(mode) {
		pterm.Error.Println("Invalid mode", mode)
		os.Exit(1)
	}
	pterm.Info.Println("Statistics mode:", mode)

	// check orderBy
	if mode == model.ModeMonsterDYPTeamStats && orderBy != "wr" && orderBy != "elo" {
		pterm.Error.Println("Invalid order", orderBy)
		os.Exit(1)
	}
	pterm.Info.Println("Order by:", orderBy)

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
	p, _ := pterm.DefaultProgressbar.
		WithTotal(len(files)).
		WithRemoveWhenDone().
		WithTitle("Processing tournaments data").
		Start()
	for _, fn := range files {
		if !strings.HasSuffix(fn, ".ktool") {
			pterm.Warning.Println("Not .ktool file, ignored:", fn)
			continue
		}
		pterm.Info.Println("Parsing", fn)
		t, err := parser.ParseTournament(fn)
		if err != nil {
			pterm.Error.Println(err)
			os.Exit(1)
		}
		tournaments = append(tournaments, *t)
		p.Increment()
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
		OrderBy:          orderBy,
		RankMinThreshold: rankMinThreshold,
		WithTime:         withTime,
		WithHomeAway:     withHomeAway,
		WithPoint:        withPoint,
		Incremental:      incremental,
	}
	if mode == model.ModeMonsterDYPPlayerStats {
		statOperator = monsterdyp.NewPlayerStats(games, players, option)
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
			if incremental {
				players = statOperator.Details()
				if len(players) > 0 {
					err = parser.WritePlayer(player, players)
					if err != nil {
						panic(err)
					}
				}
			}
		}
	} else {
		pterm.Error.Println("Unsupported tournament mode for this operator")
	}
}
