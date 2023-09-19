package rating

import (
	"math"
	"testing"
)

func TestEloK(t *testing.T) {
	type args struct {
		Played int
		Score  float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "played 20, score lt 1000",
			args: args{
				Played: 20,
				Score:  1000.0,
			},
			want: 40,
		},
		{
			name: "played 20, score gte 2400",
			args: args{
				Played: 20,
				Score:  2401.0,
			},
			want: 40,
		},
		{
			name: "played 30, score gte 2400",
			args: args{
				Played: 30,
				Score:  2401.0,
			},
			want: 10,
		},
		{
			name: "played 30, score gte 2400",
			args: args{
				Played: 30,
				Score:  2301.0,
			},
			want: 20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := Elo{}
			if got := er.chooseK(tt.args.Played, tt.args.Score); math.Round(got) != tt.want {
				t.Errorf("Calculate(%d, %.2f) = %v, want %v", tt.args.Played, tt.args.Score, got, tt.want)
			}
		})
	}
}

func TestEloScore(t *testing.T) {
	type args struct {
		Ra     int
		Rb     int
		Played int
		Result int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "same rating draw",
			args: args{
				Ra:     1000,
				Rb:     1000,
				Result: Draw,
			},
			want: 1000,
		},
		{
			name: "same rating win",
			args: args{
				Ra:     1000,
				Rb:     1000,
				Result: Win,
			},
			want: 1020,
		},
		{
			name: "win with k = 20",
			args: args{
				Ra:     1980,
				Rb:     2000,
				Played: 30,
				Result: Win,
			},
			want: 1991,
		},
		{
			name: "win with k = 10",
			args: args{
				Ra:     2452,
				Rb:     2530,
				Played: 40,
				Result: Win,
			},
			want: 2458,
		},
		{
			name: "diff rating loss",
			args: args{
				Ra:     1400,
				Rb:     1500,
				Result: Loss,
			},
			want: 1386,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := Elo{}
			factors := Factor{
				PlayerScore:   float64(tt.args.Ra),
				OpponentScore: float64(tt.args.Rb),
				Result:        tt.args.Result,
				Played:        tt.args.Played,
			}
			if got := er.Calculate(factors); math.Round(got) != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
