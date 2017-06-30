package storage

import (
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

	store := NewFileSystemStore(tempDir, NewFlatAddressSegmenter())
	address := ([]byte)("address")
	data := ([]byte)("data")

	store.Put(address, data)
	retrieved, err := store.Get(address)
	assert.NoError(t, err)
	assert.Equal(t, data, retrieved)

	stat, err := store.Stat(address)
	assert.NoError(t, err)
	assert.True(t, stat.Exists)
	assert.Equal(t, uint64(len(data)), stat.Size)

	stat, err = store.Stat(([]byte)("bar"))
	assert.NoError(t, err)
	assert.False(t, stat.Exists)
}
