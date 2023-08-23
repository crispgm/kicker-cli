// Package parser .
package parser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/crispgm/kicker-cli/internal/entity"
)

// ParsePlayer .
func ParsePlayer(fn string) ([]entity.Player, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var players []entity.Player
	err = json.Unmarshal(data, &players)
	if err != nil {
		return nil, err
	}
	return players, err
}

// WritePlayer .
func WritePlayer(fn string, players []entity.Player) error {
	b, err := json.Marshal(players)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(fn, b, 0o644)
	return err
}
