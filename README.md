# kicker-cli

A Foosball data aggregator, analyzers, and manager based on [Kickertool](https://app.kickertool.de/).

Kickertool is a powerful tournament software that enables everyone to run foosball event.
`kicker-cli` allows you to maintain data from Kickertool with the support of organization and event management,
by which organizers could manage, analyze and rank inside an organization across multiple events.

## Features

- [x] Players database
- [x] Multiple files data aggregation
- [x] Players' rank by win rate and ELO score
- [ ] Organization database
- [ ] Tournament database
- [ ] Data file management
- [ ] Result-based score

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

Init a workspace:
```shell
# init with current folder
kicker-cli init
# init with existing folder
kicker-cli init --workspace=/path/to/workspace
```

Import `.ktool` files:
```shell
# download from Kickertool and then
kicker-cli import /path/to/ktool
```

Get rank:
```shell
kicker-cli rank
```

Show players:
```shell
kicker-cli player list
```

## Known Issues

- Game duration is inaccurate because we actually don't input the result as soon as the game is finished.

## License

[MIT](/LICENSE)
