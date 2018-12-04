package storage

import (
	"fmt"

	"github.com/monax/hoard/sync"
)

// The number of mutexes that will be available for all possible addresses.
// Raising this number will decrease mutex contention between addresses (that
// would otherwise be sharing a mutex) at the expense of using a little more
// space.
const addressMutexCount = 127

type syncStore struct {
	store NamedStore
	mtx   *sync.AddressRWMutex
}

// Wrap a Store to synchronise it with respect to address access. For each
// address exactly one writer can enter the Put method of the underlying store
// or multiple readers can enter the Get and Stat methods, but no simultaneous
// readers (Getters, Statters) and writers (Putters) are allowed. Concurrent
// reads and writes to different addresses are permitted so the underlying store
// must be goroutine-safe across addresses.
func NewSyncStore(store NamedStore) *syncStore {
	return &syncStore{
		store: store,
		mtx:   sync.NewAddressRWMutex(addressMutexCount),
	}
}

var _ Store = (*syncStore)(nil)

func (ss *syncStore) Get(address []byte) (data []byte, err error) {
	ss.mtx.RLock(address)
	defer ss.mtx.RUnlock(address)
	return ss.store.Get(address)
}

func (ss *syncStore) Stat(address []byte) (*StatInfo, error) {
	ss.mtx.RLock(address)
	defer ss.mtx.RUnlock(address)
	return ss.store.Stat(address)

}

func (ss *syncStore) Put(address []byte, data []byte) ([]byte, error) {
	ss.mtx.Lock(address)
	defer ss.mtx.Unlock(address)
	return ss.store.Put(address, data)
}

func (ss *syncStore) Location(address []byte) string {
	return ss.store.Location(address)
}

func (ss *syncStore) Name() string {
	return fmt.Sprintf("syncStore[mutexCount=%v](%s)", ss.mtx.Size(),
		ss.store.Name())
}
