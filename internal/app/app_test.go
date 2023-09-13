package app

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/crispgm/kicker-cli/internal/util"
	"github.com/stretchr/testify/assert"
)

func TestAppLoadConf(t *testing.T) {
	appErr := NewApp("/tmp/kicker/not/here", "")
	assert.NotNil(t, appErr.LoadConf())

	path := util.GetCIPath("../..")

	app := NewApp(path, ".kicker.yaml")
	assert.NotEmpty(t, app.Version)
	assert.Nil(t, app.LoadConf())
	assert.Equal(t, "1", app.Conf.ManifestVersion)
	assert.Equal(t, "My Foosball Community", app.Conf.Organization.Name)
	assert.NotEmpty(t, app.Conf.Players)
	assert.NotEmpty(t, app.Conf.Events)

	assert.Equal(t, fmt.Sprintf("%s/%s", path, "data"), app.DataPath())
	assert.NotNil(t, app.GetEvent("191321a4-a709-4ff2-8ad2-ce8c20bb8265"))
	assert.Nil(t, app.GetEvent("d68fa-5f9e-49df-b576-12b29299c215"))
	numOfEvents := len(app.Conf.Events)
	assert.Error(t, app.DeleteEvent("d68fa-5f9e-49df-b576-12b29299c215"))
	assert.Nil(t, app.DeleteEvent("191321a4-a709-4ff2-8ad2-ce8c20bb8265"))
	assert.Len(t, app.Conf.Events, numOfEvents-1)

	assert.NotNil(t, app.GetPlayer("dbb556e1-ecb2-4aad-bbc7-ed240a1b9bfc"))
	assert.Nil(t, app.GetPlayer("dbb556e1-1111-4aad-bbc7-ed240a1b9bfc"))
	numOfPlayers := len(app.Conf.Players)
	assert.Error(t, app.DeletePlayer("d68fa-5f9e-49df-b576-12b29299c215"))
	assert.Nil(t, app.DeletePlayer("dbb556e1-ecb2-4aad-bbc7-ed240a1b9bfc"))
	assert.Len(t, app.Conf.Players, numOfPlayers-1)
}

func TestAppWriteConf(t *testing.T) {
	seed := rand.Int()
	fn := fmt.Sprintf(".kicker.yaml.%d", seed)
	app := NewApp("/tmp", fn)
	err := app.WriteConf()
	assert.NoError(t, err)
	assert.FileExists(t, app.FilePath)
}
