// Package parser .
package parser

import (
	"encoding/json"
	"io/ioutil"

	"github.com/crispgm/kickertool-analyzer/pkg/ktool/model"
)

// ParseTournament .
func ParseTournament(fn string) (*model.Tournament, error) {
	data, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}

	var tournament model.Tournament
	err = json.Unmarshal(data, &tournament)
	if err != nil {
		return nil, err
	}
	return &tournament, err
}
