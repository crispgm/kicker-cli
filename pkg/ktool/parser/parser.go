// Package parser .
package parser

import (
	"encoding/json"
	"os"

	"github.com/crispgm/kicker-cli/pkg/ktool/model"
)

// ParseFile .
func ParseFile(fn string) (*model.Tournament, error) {
	data, err := os.ReadFile(fn)
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
