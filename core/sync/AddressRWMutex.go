package sync

import (
	"sync"

	"github.com/OneOfOne/xxhash"
)

type AddressRWMutex struct {
	mtxs       []sync.RWMutex
	hasherPool sync.Pool
	mutexCount uint64
}

// Create a mutex that provides a pseudo-independent set of mutexes for addresses
// where the address space is mapped into possibly much smaller set of backing
// mutexes using the xxhash (non-cryptographic)
// hash function // modulo size. If some addresses collide modulo size they will be unnecessary
// contention between those addresses, but you can trade space against contention
// as desired.
func NewAddressRWMutex(mutexCount int) *AddressRWMutex {
	return &AddressRWMutex{
		// max slice length is bounded by max(int) thus the argument type
		mtxs: make([]sync.RWMutex, mutexCount, mutexCount),
		hasherPool: sync.Pool{
			New: func() interface{} {
				return xxhash.New64()
			},
		},
		mutexCount: uint64(mutexCount),
	}
}

func (mtx *AddressRWMutex) Lock(address []byte) {
	mtx.mutex(address).Lock()
}

func (mtx *AddressRWMutex) Unlock(address []byte) {
	mtx.mutex(address).Unlock()
}

func (mtx *AddressRWMutex) RLock(address []byte) {
	mtx.mutex(address).RLock()
}

func (mtx *AddressRWMutex) RUnlock(address []byte) {
	mtx.mutex(address).RUnlock()
}

// Return the size of the underlying array of mutexes
func (mtx *AddressRWMutex) Size() uint64 {
	return mtx.mutexCount
}

func (mtx *AddressRWMutex) mutex(address []byte) *sync.RWMutex {
	return &mtx.mtxs[mtx.index(address)]
}

func (mtx *AddressRWMutex) index(address []byte) uint64 {
	return mtx.hash(address) % mtx.mutexCount
}

func (mtx *AddressRWMutex) hash(address []byte) uint64 {
	h := mtx.hasherPool.Get().(*xxhash.XXHash64)
	defer func() {
		h.Reset()
		mtx.hasherPool.Put(h)
	}()
	h.Write(address)
	return h.Sum64()
}
