# Ranking System

## Kicker Ranking System

We provides and implemented both merit-based ranking (KRP) and ELO-based (KES) ranking.
Otherwise, ITSF and ATSA ranking systems are already built-in.
All the points and scores can be used simultenously.

### Kicker Ranking Points (KRP)

KRP follows ITSF points but with our own interpretation of event class:

| Place | World | Continental | Domestic | Local | Casual |
| ----- | ----- | ----------- | -------- | ----- | ------ |
| 1     | 200   | 150         | 100      | 50    | 10     |
| 2     | 180   | 135         | 90       | 45    | 8      |
| 3     | 160   | 120         | 80       | 40    | 6      |
| 4     | 140   | 105         | 70       | 35    | 4      |
| 5     | 120   | 90          | 60       | 30    | 2      |
| 9     | 100   | 75          | 50       | 25    | 1      |
| 17    | 80    | 60          | 40       | 20    | 0      |
| 33    | 60    | 45          | 30       | 15    | 0      |
| 65    | 40    | 30          | 20       | 10    | 0      |
| 129   | 20    | 15          | 10       | 5     | 0      |
| 257   | 12    | 9           | 6        | 3     | 0      |
| 513   | 4     | 3           | 2        | 1     | 0      |

### Kicker ELO Scores (KES)

KES follows FIDE K-factor's choice:

| K-factor | Used for players with ratings ...   |
| -------- | ----------------------------------- |
| K = 40   | below 30 games                      |
| K = 20   | under 2400                          |
| K = 10   | at least 2400 and at least 30 games |

## Reference

- [ITSF PLAYERS WORLD RANKING SYSTEM](https://www.tablesoccer.org/rules/documents/ITSF_Player_Ranking_System.pdf)
- [ATSA Ranking Points](https://asiatablesoccer.glide.page/dl/Players/s/255ce0)
- [Elo rating system - Wikipedia](https://en.wikipedia.org/wiki/Elo_rating_system)
