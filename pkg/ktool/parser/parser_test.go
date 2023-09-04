package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/stretchr/testify/assert"
)

func TestParseGame(t *testing.T) {
	ciMode := os.Getenv("KICKER_CLI_CI_MODE")
	path := "../../.."
	if ciMode == "1" {
		path = "."
	}

	testCases := []struct {
		mode      string
		rounds    bool
		ko        bool
		leftLevel bool
	}{
		{model.ModeSwissSystem, true, true, true},
		{model.ModeRounds, true, true, false},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s", tc.mode), func(t *testing.T) {
			fn := fmt.Sprintf("%s/test/data/test_%s.ktool", path, tc.mode)
			trn, err := ParseFile(fn)
			if assert.NoError(t, err) {
				assert.Equal(t, "table_soccer", trn.Sport.Name)
				assert.Equal(t, tc.mode, trn.Mode)
				assert.NotZero(t, trn.NumRounds)
				assert.Equal(t, "Tournament", trn.Type)

				if tc.rounds {
					assert.NotZero(t, trn.Rounds)
					assert.Equal(t, "Round", trn.Rounds[0].Type)
					assert.Equal(t, "1", trn.Rounds[0].Name)
				}

				if tc.ko {
					assert.NotEmpty(t, trn.KnockOffs)
					assert.Equal(t, "KO", trn.KnockOffs[0].Type)
					assert.NotEmpty(t, trn.KnockOffs[0].Levels)
					assert.Equal(t, "Level", trn.KnockOffs[0].Third.Type)
					if tc.leftLevel {
						assert.NotEmpty(t, trn.KnockOffs[0].LeftLevels)
					}
				}
			}
		})
	}
}
