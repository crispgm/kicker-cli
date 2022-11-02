# kickertool-analyzer

Data analyzer tools of Kickertool.

## kicker-cli

Statistics data of multiple `.ktool` files, by which organizers could analyzer and rank inside an organization across multiple events.

_Notice_: Only Monster DYP mode is supported.

```bash
go install github.com/crispgm/kickertool-analyzer/cmd/kicker-cli@latest
```

### Usage

1. Download `.ktool` files
2. Setup `players.json`
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

## pelo

Simple tool to show estimated ELO changes between two teams/players.

```shell
$ pelo 1100 1200
$ pelo 1103 1203 1289 1013
$ pelo -k 20 1103 1203 1289 1013
```
