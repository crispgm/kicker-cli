package app

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppLoadConf(t *testing.T) {
	appErr := NewApp("/tmp/kicker/not/here", "")
	assert.NotNil(t, appErr.LoadConf())

	ciMode := os.Getenv("KICKER_CLI_CI_MODE")
	path := "../.."
	if ciMode == "1" {
		path = "."
	}

	app := NewApp(path, ".kicker.yaml")
	assert.NotEmpty(t, app.Version)
	assert.Nil(t, app.LoadConf())
	assert.Equal(t, "1", app.Conf.ManifestVersion)
	assert.Equal(t, "My Foosball Community", app.Conf.Organization.Name)
	assert.NotEmpty(t, app.Conf.Players)
	assert.NotEmpty(t, app.Conf.Events)

	assert.Equal(t, fmt.Sprintf("%s/%s", path, "data"), app.DataPath())
	assert.NotNil(t, app.GetEvent("978d68fa-5f9e-49df-b576-12b29299c215"))
	assert.Nil(t, app.GetEvent("d68fa-5f9e-49df-b576-12b29299c215"))
}

func TestAppWriteConf(t *testing.T) {
	seed := rand.Int()
	fn := fmt.Sprintf(".kicker.yaml.%d", seed)
	app := NewApp("/tmp", fn)
	err := app.WriteConf()
	assert.NoError(t, err)
	assert.FileExists(t, app.FilePath)
}
