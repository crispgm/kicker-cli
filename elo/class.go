package elo

import "math"

// Classes definition
var Classes = []Class{
	{
		Low:   2400,
		High:  math.MaxInt32,
		Title: "Grand Master",
	},
	{
		Low:   2200,
		High:  2399,
		Title: "Master",
	},
	{
		Low:   2000,
		High:  2199,
		Title: "Expert",
	},
	{
		Low:   1800,
		High:  1999,
		Title: "Pro",
	},
	{
		Low:   1600,
		High:  1799,
		Title: "Semi-Pro A",
	},
	{
		Low:   1400,
		High:  1599,
		Title: "Semi-Pro B",
	},
	{
		Low:   1200,
		High:  1399,
		Title: "Amateur A",
	},
	{
		Low:   1000,
		High:  1199,
		Title: "Amateur B",
	},
	{
		Low:   800,
		High:  999,
		Title: "Rookie A",
	},
	{
		Low:   600,
		High:  799,
		Title: "Rookie B",
	},
	{
		Low:   400,
		High:  599,
		Title: "Rookie C",
	},
	{
		Low:   200,
		High:  399,
		Title: "Starter A",
	},
	{
		Low:   100,
		High:  199,
		Title: "Starter B",
	},
}
