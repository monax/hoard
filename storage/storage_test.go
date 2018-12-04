package storage

import (
	"testing"

	"crypto/sha256"
	"sync"

	"strconv"

	"github.com/stretchr/testify/assert"
)

const concurrentAccessGoRoutines = 200

// Generic test suite for all Stores
func testStore(t *testing.T, store Store) {
	address := bs("address")
	data := bs("data")
	getPutGet(t, store, address, data)

	stat, err := store.Stat(address)
	if assert.NoError(t, err) {
		assert.True(t, stat.Exists)
		assert.Equal(t, uint64(len(data)), stat.Size)
	}

	stat, err = store.Stat(bs("bar"))
	if assert.NoError(t, err) {
		assert.False(t, stat.Exists)
	}

	retrieved, err := store.Get(bs("foo"))
	assert.Nil(t, retrieved)
	assert.Error(t, err)

	// Has a '/' under standard encoding
	getPutGet(t, store, []byte{0, 0, 63, 0, 0}, bs("bar-data"))

	testConcurrentContentAddressedStore(t, store)
}

func testConcurrentContentAddressedStore(t *testing.T, store Store) {
	cas := NewContentAddressedStore(MakeAddresser(sha256.New), store)
	wg := new(sync.WaitGroup)
	wg.Add(concurrentAccessGoRoutines)
	for i := 0; i < concurrentAccessGoRoutines; i++ {
		n := int64(i)
		go func() {
			data := bs("data", strconv.FormatInt(n, 2))
			putGetCAS(t, cas, data)
			wg.Done()
		}()
	}
	wg.Wait()
}

func getPutGet(t *testing.T, store Store, address, data []byte) {
	retrieved, err := store.Get(address)
	assert.Nil(t, retrieved, "Should be nothing at address")
	assert.Error(t, err, "Getting an address with no data should be "+
		"an error")

	// Put data at address
	address, err = store.Put(address, data)
	assert.NoError(t, err, "Should be able to Put data at address")

	retrieved, err = store.Get(address)
	assert.NoError(t, err, "Should be able to Get data from address")
	assert.Equal(t, data, retrieved)
}

func putGetCAS(t *testing.T, store ContentAddressedStore, data []byte) {
	// Put data at address
	address, err := store.Put(data)
	assert.NoError(t, err, "Should be able to Put data at address")

	retrieved, err := store.Get(address)
	assert.NoError(t, err, "Should be able to Get data from address")
	assert.Equal(t, data, retrieved)
}

func bs(strs ...string) []byte {
	var b []byte
	for _, s := range strs {
		b = append(b, []byte(s)...)
	}
	return b
}
