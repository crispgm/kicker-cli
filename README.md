# kicker-cli

Foosball organization, event, and data analyzers based on [Kickertool]().

Kickertool is a powerful tournament software that enables everyone to run foosball event.
`kicker-cli` allows you to maintain data from Kickertool with the support of organization and event management,
by which organizers could manage, analyze and rank inside an organization across multiple events.

## Features

- [x] Players database
- [x] Multiple files data aggregation
- [x] Players' rank by win rate and ELO score
- [ ] Organization
- [ ] Data file management
- [ ] Result-based score
- [ ] Tournament support with levels

### Game Modes

- [x] MonsterDYP

### Operators

- Player Statistics
- Team Statistics

## Installation

```bash
go install github.com/crispgm/kicker-cli/cmd/kicker-cli@latest
```

## Usage

1. Download `.ktool` files
2. Setup players' database
3. Run command

```text
Usage of kicker-cli:

  -dry-run
        Dry Run
  -elo-k int
        Elo K factor (default 10)
  -incremental
        Update player's data incrementally
  -mode string
        Stat mode. Supported: mdp, mdt
  -order-by wr
        Order by wr (win rate) or `elo` (ELO ranking) (default "wr")
  -player string
        Players' data file
  -rmt int
        Rank minimum threshold
  -with-home-away
        With home/away analysis
  -with-point
        With point analysis
  -with-time
        With time analysis
```

## Known Issues

- Game duration is inaccurate because we actually don't input the result as soon as the game is finished.

## License

[MIT](/LICENSE)
