// Package cmd .
package cmd

import (
	"github.com/spf13/cobra"
)

var projectBase string

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
}
