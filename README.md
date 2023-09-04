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
- Rank for players and teams with win rate and ELO score
- [ ] Result-based score

### Game Modes

#### Round Games

- MonsterDYP
- Rounds
- [ ] Round Robin
- Swiss System

#### Knockoff Games

- Double Elimination
- Elimination

### Operators

- Player Statistics
- [ ] Player Rivals
- Team Statistics
- [ ] Team Rivals

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
kicker-cli init --path=/path/to/workspace
```

Import `.ktool` files:
```shell
# download from Kickertool and then
kicker-cli import --path=/path/to/ktool
```

Get rank:
```shell
kicker-cli rank --name "my event name"
kicker-cli rank --all
```

Show players:
```shell
kicker-cli player
```

ELO emulation:
```shell
kicker-cli elo 1100 1200
kicker-cli elo 1103 1203 1289 1013
kicker-cli elo 1103 1203 1289 1013 -k 20
```

## Known Issues

- Game duration is inaccurate because we actually don't input the result as soon as the game is finished.

## License

[MIT](/LICENSE)
