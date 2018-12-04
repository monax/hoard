package hoard

import (
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/reference"
	"github.com/monax/hoard/storage"
	"github.com/stretchr/testify/assert"
)

func TestDeterministicEncryptedStore(t *testing.T) {
	hrd := NewHoard(storage.NewMemoryStore(), log.NewNopLogger())
	bunsIn := bs("hot buns")

	ref, err := hrd.Put(bunsIn, nil)
	assert.NoError(t, err)

	bunsOut, err := hrd.Get(ref)
	assert.Equal(t, bunsIn, bunsOut)

	_, err = hrd.Get(reference.New(ref.Address, pad("wrong secret", 32), nil))
	assert.Error(t, err)

	statInfo, err := hrd.Store().Stat(ref.Address)
	assert.NoError(t, err)
	assert.True(t, statInfo.Exists)
	// Our GCM cipher should be running an overhead of 16 bytes
	// (no IV, but 16-byte authentication tag)
	assert.Equal(t, uint64(len(bunsIn))+16, statInfo.Size)

	loc := hrd.Store().Location(ref.Address)
	assert.Equal(t, "memfs://2768b87bc22a7f0f20eb268937a00d682bc5dc21a62badf1d24c78633436a1b7", loc)

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
