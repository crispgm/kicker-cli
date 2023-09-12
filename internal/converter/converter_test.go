package converter

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/crispgm/kicker-cli/internal/entity"
	"github.com/crispgm/kicker-cli/pkg/ktool/model"
	"github.com/crispgm/kicker-cli/pkg/ktool/parser"
)

func TestConvertDoubleGames(t *testing.T) {
	ciMode := os.Getenv("KICKER_CLI_CI_MODE")
	path := "../.."
	if ciMode == "1" {
		path = "."
	}

	fn := fmt.Sprintf("%s/test/data/test_%s.ktool", path, "swiss")
	trn, err := parser.ParseFile(fn)
	if assert.NoError(t, err) && assert.NotNil(t, trn) {
		var ePlayers []entity.Player
		for _, p := range trn.Players {
			ePlayers = append(ePlayers, *entity.NewPlayer(p.Name))
		}
		nc := NewConverter()
		rec, err := nc.Normalize([]model.Tournament{*trn}, ePlayers)
		if assert.NoError(t, err) && assert.NotNil(t, rec) {
			assert.Len(t, rec.PreliminaryRounds, 10)
			assert.Len(t, rec.LoserBracket, 3)
			assert.Len(t, rec.WinnerBracket, 5)
			assert.Nil(t, rec.ThirdPlace)
			assert.Len(t, rec.AllGames, 18)
			assert.NotEmpty(t, nc.Briefing())
		}
	}
}
