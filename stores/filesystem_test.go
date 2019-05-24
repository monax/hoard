package stores

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileSystemStore(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "filesystem_test")
	defer func() {
		err := os.RemoveAll(tempDir)
		if err != nil {
			panic(err)
		}
	}()
	assert.NoError(t, err)

	fss, err := NewFileSystemStore(tempDir, base64.URLEncoding)

	assert.NoError(t, err)
	RunTests(t, fss)
}
