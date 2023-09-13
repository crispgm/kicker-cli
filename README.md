# kicker-cli

[![build](https://github.com/crispgm/kicker-cli/actions/workflows/build.yml/badge.svg)](https://github.com/crispgm/kicker-cli/actions/workflows/build.yml)
[![codecov](https://codecov.io/gh/crispgm/kicker-cli/graph/badge.svg?token=MMGE8I5YNI)](https://codecov.io/gh/crispgm/kicker-cli)

A Foosball data aggregator, analyzers, and manager based on [Kickertool](https://app.kickertool.de/).

Kickertool is a powerful tournament software that enables everyone to run foosball event.
`kicker-cli` allows you to maintain data from Kickertool with the support of organization and event management,
by which organizers could manage, analyze and rank inside an organization across multiple events.

## Features

- Organization, Events, and Players database
- Multiple files data aggregation
- Rank for players and teams with [points](/docs/ranking_system.md), win rate and ELO score

### Game Modes

#### Round Games

- MonsterDYP
- Rounds
- Round Robin
- Swiss System

#### Knockoff Games

- Double Elimination
- Elimination

### Operators

- Double Player Ranks: BYP, DYP, and MonsterDYP
- Double Team Ranks: BYP, DYP, and MonsterDYP
- Single Player Ranks: Single
- Double Team Rivals: BYP, DYP, and MonsterDYP
- Single Player Rivals: Single

## Installation

```bash
go install github.com/crispgm/kicker-cli/cmd/kicker-cli@latest
```

## Usage

Init a workspace:
```shell
# init with current folder
kicker-cli init
# init with existing folder
kicker-cli init /path/to/workspace
```

Import `.ktool` files:
```shell
# download from Kickertool and then
kicker-cli import /path/to/ktool
```

Event management:
```shell
# list events
kicker-cli event ls
# show info
kicker-cli event info 351e00bf-025c-4243-b381-2f5a135c3070
# get rank
kicker-cli event rank -m double_player_rank -t byp
kicker-cli event rank 351e00bf-025c-4243-b381-2f5a135c3070 -m double_player_rank
# open event
kicekr-cli event open 351e00bf-025c-4243-b381-2f5a135c3070
#  => xdg-open event-url
```

Show players:
```shell
kicker-cli player
```

Evaluation:
```shell
kicker-cli eval -a elo 1100 1200
kicker-cli eval -a elo 1103 1203 1289 1013
kicker-cli eval -a elo -k 20 1103 1203 1289 1013
```

## Known Issues

- Game duration is inaccurate because we actually don't input the result as soon as the game is finished.

## License

[MIT](/LICENSE)
