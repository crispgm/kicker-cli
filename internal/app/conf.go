package app

import (
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/crispgm/kicker-cli/internal/entity"
)

// DefaultName of configuration file
const DefaultName = ".kicker.yaml"

// Conf .
type Conf struct {
	ManifestVersion string `yaml:"manifest_version"`

	Organization entity.Organization `yaml:"organization"`

	Events  []entity.Event  `yaml:"events"`
	Players []entity.Player `yaml:"players"`
}

// LoadConf .
func (app *App) LoadConf() error {
	data, err := ioutil.ReadFile(app.FilePath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &app.Conf)
	if err != nil {
		return err
	}
	return nil
}

// WriteConf .
func (app *App) WriteConf() error {
	b, err := yaml.Marshal(app.Conf)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(app.FilePath, b, 0o644)
	return err
}
