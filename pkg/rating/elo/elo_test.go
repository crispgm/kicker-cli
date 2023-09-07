package elo

import (
	"testing"

	"github.com/crispgm/kicker-cli/pkg/rating"
)

func TestEloScore(t *testing.T) {
	type args struct {
		Ra          int
		Rb          int
		K           int
		WonDrewLost int
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{
			name: "same rating draw",
			args: args{
				Ra:          1000,
				Rb:          1000,
				K:           40,
				WonDrewLost: rating.Drew,
			},
			want: 1000,
		},
		{
			name: "same rating win",
			args: args{
				Ra:          1000,
				Rb:          1000,
				K:           40,
				WonDrewLost: rating.Won,
			},
			want: 1020,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			er := EloRating{K: float64(tt.args.K)}
			er.InitialScore(float64(tt.args.Ra), float64(tt.args.Rb))
			if got := er.Calculate(tt.args.WonDrewLost); got != tt.want {
				t.Errorf("Calculate() = %v, want %v", got, tt.want)
			}
		})
	}
}
