package util

import "os"

// GetCIPath .
func GetCIPath(cwd string) string {
	ciMode := os.Getenv("KICKER_CLI_CI_MODE")
	path := cwd + "/test"
	if ciMode == "1" {
		path = "./test"
	}

	return path
}
