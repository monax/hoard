package storage

import (
	"io/ioutil"
	"os"
	"testing"

	"encoding/base64"

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
	testStore(t, fss)
}
