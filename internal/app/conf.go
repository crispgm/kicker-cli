package app

import (
	"errors"
	"os"

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
	data, err := os.ReadFile(app.FilePath)
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
	err = os.WriteFile(app.FilePath, b, 0o644)
	return err
}

// AddEvent .
func (app *App) AddEvent(events ...entity.Event) {
	app.Conf.Events = append(app.Conf.Events, events...)
}

// GetEvent returns event with the given id. Otherwise, return nil.
func (app App) GetEvent(id string) *entity.Event {
	for _, e := range app.Conf.Events {
		if id == e.ID {
			return &e
		}
	}

	return nil
}

// DeleteEvent delete an event
func (app *App) DeleteEvent(id string) error {
	s := -1
	for i, e := range app.Conf.Events {
		if id == e.ID {
			s = i
		}
	}
	if s < 0 {
		return errors.New("Event not found")
	}

	app.Conf.Events = append(app.Conf.Events[:s], app.Conf.Events[s+1:]...)
	return nil
}

// AddPlayer .
func (app *App) AddPlayer(players ...entity.Player) {
	app.Conf.Players = append(app.Conf.Players, players...)
}

// GetPlayer returns player with the given id. Otherwise, return nil.
func (app App) GetPlayer(id string) *entity.Player {
	for _, p := range app.Conf.Players {
		if id == p.ID {
			return &p
		}
	}

	return nil
}

// DeletePlayer delete a player
func (app *App) DeletePlayer(id string) error {
	s := -1
	for i, p := range app.Conf.Players {
		if id == p.ID {
			s = i
		}
	}
	if s < 0 {
		return errors.New("Player not found")
	}

	app.Conf.Players = append(app.Conf.Players[:s], app.Conf.Players[s+1:]...)
	return nil
}
