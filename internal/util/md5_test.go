package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5(t *testing.T) {
	ciMode := os.Getenv("KICKER_CLI_CI_MODE")
	path := "../.."
	if ciMode == "1" {
		path = "."
	}
	path += "/test/data/test_round_robin.ktool"
	md5, err := MD5CheckSum(path)
	assert.NoError(t, err)
	assert.Equal(t, "3ba7c806a52baf4efbfb4962f62a36d1", md5)
}
