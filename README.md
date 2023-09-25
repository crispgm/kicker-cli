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
- Ranks for players and teams with performance (win rate) and [ranking systems (points, & ELO score)](/docs/ranking_system.md)

### Analyzers

- Double Player Ranks
- Double Team Ranks
- Single Player Ranks
- Double Team Rivals
- Single Player Rivals

### Supported Game Modes

#### Round Games

- MonsterDYP
- Rounds
- Round Robin
- Swiss System

#### Knockoff Games

- Double Elimination
- Elimination

## Installation

```bash
go install github.com/crispgm/kicker-cli/cmd/kicker-cli@latest
```

## Usage

Init a workspace:
```shell
# init with current folder
kicker-cli init --name "MyAwesomeFoos"

# init with existing folder
kicker-cli init --name "MyAwesomeFoos" --path /path/to/workspace

# show organization
kicker-cli org
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
kicker-cli event ls --before 2023-08-23
kicker-cli event ls --name-type dyp

# show info
kicker-cli event info 351e00bf-025c-4243-b381-2f5a135c3070

# get rank
kicker-cli event rank -m double_player_rank -t byp
kicker-cli event rank 351e00bf-025c-4243-b381-2f5a135c3070 -m double_player_rank

# open event
kicker-cli event open 351e00bf-025c-4243-b381-2f5a135c3070
#  => xdg-open event-url

# delete event
kicker-cli event delete 351e00bf-025c-4243-b381-2f5a135c3070
```

Player management:
```shell
# list players
kicker-cli player ls

# delete player
kicker-cli player delete 13d4ea60-f6ff-48da-be1e-413d38328cb0
```

Evaluation:
```shell
# elo scores
kicker-cli eval -a elo 1100 1200
kicker-cli eval -a elo 1103 1203 1289 1013
kicker-cli eval -a elo -k 20 1103 1203 1289 1013

# points gained
kicker-cli eval -a rank -s ATSA50 -s KLocal 1
```

## Known Issues & Limitations

- Game duration is inaccurate because we actually don't input the result as soon as the game is finished. So we abandoned all time factors.

## License

[MIT](/LICENSE)
