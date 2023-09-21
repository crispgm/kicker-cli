package rating

import (
	"sort"
)

var _ Rating = (*Rank)(nil)

var kickerRankTable = map[string]map[int]int{
	KWorld:       {1: 200, 2: 180, 3: 160, 4: 140, 5: 120, 9: 100, 17: 80, 33: 60, 65: 40, 129: 20, 257: 12, 513: 4},
	KContinental: {1: 150, 2: 135, 3: 120, 4: 105, 5: 90, 9: 75, 17: 60, 33: 45, 65: 30, 129: 15, 257: 9, 513: 3},
	KDomestic:    {1: 100, 2: 90, 3: 80, 4: 70, 5: 60, 9: 50, 17: 40, 33: 30, 65: 20, 129: 10, 257: 6, 513: 2},
	KLocal:       {1: 50, 2: 45, 3: 40, 4: 35, 5: 30, 9: 25, 17: 20, 33: 15, 65: 10, 129: 5, 257: 3, 513: 1},
	KCasual:      {1: 10, 2: 8, 3: 6, 4: 4, 5: 2, 9: 1},
}

var atsaRankTable = map[string]map[int]int{
	ATSA2000: {1: 2000, 2: 1200, 3: 720, 5: 360, 9: 180, 17: 90, 33: 45},
	ATSA1000: {1: 1000, 2: 600, 3: 360, 4: 240, 5: 180, 7: 120, 9: 90, 17: 45, 33: 25},
	ATSA500:  {1: 500, 2: 300, 3: 180, 4: 120, 5: 90, 7: 60, 9: 45, 13: 30, 17: 20, 33: 10},
	ATSA50:   {1: 50, 2: 30, 3: 20, 4: 17, 5: 15, 7: 12, 9: 9, 17: 4, 33: 1},
}

var itsfRankTable = map[string]map[int]int{
	ITSFWorldSeries:   {1: 200, 2: 180, 3: 160, 4: 140, 5: 120, 9: 100, 17: 80, 33: 60, 65: 40, 129: 20, 257: 12, 513: 4},
	ITSFInternational: {1: 150, 2: 135, 3: 120, 4: 105, 5: 90, 9: 75, 17: 60, 33: 45, 65: 30, 129: 15, 257: 9, 513: 3},
	ITSFMasterSeries:  {1: 100, 2: 90, 3: 80, 4: 70, 5: 60, 9: 50, 17: 40, 33: 30, 65: 20, 129: 10, 257: 6, 513: 2},
	ITSFProTour:       {1: 25, 2: 22, 3: 19, 4: 16, 5: 13, 9: 10, 17: 7, 33: 4, 65: 2, 129: 1},
}

// Rank calculates rank based rating
type Rank struct {
}

// InitialScore .
func (r Rank) InitialScore() float64 {
	return 0
}

// Calculate .
func (r Rank) Calculate(factors Factor) float64 {
	var rankTable map[string]map[int]int
	if factors.IsITSF() {
		rankTable = itsfRankTable
	} else if factors.IsATSA() {
		rankTable = atsaRankTable
	} else {
		rankTable = kickerRankTable
	}
	curScore := factors.PlayerScore
	incr := 0.0
	if table, ok := rankTable[factors.Level]; ok {
		var sortPts [][]int
		for pos, pts := range table {
			sortPts = append(sortPts, []int{pos, pts})
		}
		sort.SliceStable(sortPts, func(i, j int) bool {
			return sortPts[i][0] < sortPts[j][0]
		})
		incr = float64(sortPts[0][1])
		for _, pts := range sortPts {
			if factors.Place >= pts[0] {
				incr = float64(pts[1])
			} else {
				return curScore + incr
			}
		}
	}
	return 0
}
