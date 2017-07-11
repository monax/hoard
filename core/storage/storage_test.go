package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func testStore(t *testing.T, store Store) {
	address := bs("address")
	data := bs("data")

	retrieved, err := store.Get(address)
	assert.Nil(t, retrieved, "Should be nothing at address")
	assert.Error(t, err, "Getting an address with no data should be "+
		"an error")

	// Put data at address
	err = store.Put(address, data)
	assert.NoError(t, err, "Should be able to Put data at address")

	retrieved, err = store.Get(address)
	assert.NoError(t, err, "Should be able to Get data from address")
	assert.Equal(t, data, retrieved)

	stat, err := store.Stat(address)
	if assert.NoError(t, err) {
		assert.True(t, stat.Exists)
		assert.Equal(t, uint64(len(data)), stat.Size)
	}

	stat, err = store.Stat(bs("bar"))
	if assert.NoError(t, err) {
		assert.False(t, stat.Exists)
	}

	retrieved, err = store.Get(bs("foo"))
	assert.Nil(t, retrieved)
	assert.Error(t, err)

}

func bs(s string) []byte {
	return ([]byte)(s)
}
