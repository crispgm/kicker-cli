package util

import "testing"

func TestHomeDir(t *testing.T) {
	t.Log(ExpandUserHomeDir("~"))
	t.Log(ExpandUserHomeDir("~/dev"))
}
