package stores

import (
	"fmt"

	"github.com/monax/hoard/v4/sync"
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

func (inv *syncStore) Get(address []byte) (data []byte, err error) {
	inv.mtx.RLock(address)
	defer inv.mtx.RUnlock(address)
	return inv.store.Get(address)
}

func (inv *syncStore) Stat(address []byte) (*StatInfo, error) {
	inv.mtx.RLock(address)
	defer inv.mtx.RUnlock(address)
	return inv.store.Stat(address)

}

func (inv *syncStore) Put(address []byte, data []byte) ([]byte, error) {
	inv.mtx.Lock(address)
	defer inv.mtx.Unlock(address)
	return inv.store.Put(address, data)
}

func (inv *syncStore) Location(address []byte) string {
	return inv.store.Location(address)
}

func (inv *syncStore) Name() string {
	return fmt.Sprintf("syncStore[mutexCount=%v](%s)", inv.mtx.Size(),
		inv.store.Name())
}
