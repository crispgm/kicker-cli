package util

import (
	"os/user"
	"path/filepath"
	"strings"
)

// ExpandUserHomeDir check whether path starts with tilde("~")
// If it is, expand home dir
// Otherwise, just return
func ExpandUserHomeDir(path string) string {
	usr, err := user.Current()
	if err != nil {
		return path
	}
	dir := usr.HomeDir

	if path == "~" {
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		path = filepath.Join(dir, path[2:])
	}
	return path
}
