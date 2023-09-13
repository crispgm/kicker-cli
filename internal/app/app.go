// Package app .
package app

import (
	"errors"
	"path/filepath"

	"github.com/crispgm/kicker-cli/internal/entity"
)

// Version of app
const Version = "1.0.0"

// App .
type App struct {
	Version  string
	Path     string
	FilePath string
	Name     string
	Conf     Conf
}

// NewApp creates an app instance
func NewApp(path, name string) *App {
	return &App{
		Version:  Version,
		Path:     path,
		Name:     name,
		FilePath: filepath.Join(path, name),
	}
}

// DataPath returns path to data files
func (app App) DataPath() string {
	return filepath.Join(app.Path, "/data")
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
