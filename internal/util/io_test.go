package util

import (
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCopyFile(t *testing.T) {
	path := GetCIPath("../..")
	seed := rand.Int()
	src := path + "/../README.md"
	dst := fmt.Sprintf("/tmp/README.md.%d", seed)
	CopyFile(src, dst)
	if assert.FileExists(t, dst) {
		os.Remove(dst)
	}

	assert.Panics(t, func() {
		CopyFile("REAMDE.md", "/tmp/xxxx/xxxxx")
	})
	assert.Error(t, CopyFile(src, "/tmp/xxx/xxx"))
}
