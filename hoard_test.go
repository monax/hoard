package hoard

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/reference"
	"github.com/monax/hoard/v8/stores"
	"github.com/stretchr/testify/assert"
)

func TestDeterministicEncryptedStore(t *testing.T) {
	hrd := NewHoard(stores.NewMemoryStore(), config.NoopSecretManager, log.NewNopLogger())
	bunsIn := bs("hot buns")

	ref, err := hrd.Put(bunsIn, make([]byte, 32))
	assert.NoError(t, err)

	bunsOut, err := hrd.Get(ref)
	assert.Equal(t, bunsIn, bunsOut)

	_, err = hrd.Get(reference.New(ref.Address, pad("wrong secret", 32), nil, 1024))
	assert.Error(t, err)

	statInfo, err := hrd.Store().Stat(ref.Address)
	assert.NoError(t, err)
	assert.True(t, statInfo.Exists)
	// Our GCM cipher should be running an overhead of 16 bytes
	// (no IV, but 16-byte authentication tag)
	assert.Equal(t, uint64(len(bunsIn))+16+32, statInfo.Size_)

	loc := hrd.Store().Location(ref.Address)
	assert.Equal(t, "memfs://75b382c29b0d8382a09b856f7a0f00300548c9f369574f68cfc9c62fcab2d1dc", loc)

	// flip LSB of first byte of address to get an non-existent address
	ref.Address[0] = ref.Address[0] ^ 1
	statInfo, err = hrd.Store().Stat(ref.Address)
	assert.NoError(t, err)
	assert.False(t, statInfo.Exists)
}

func bs(s string) []byte {
	return ([]byte)(s)
}

func pad(s string, n int) []byte {
	b := make([]byte, n)
	copy(b, bs(s))
	return b
}
