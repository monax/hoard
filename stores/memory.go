package stores

import (
	"fmt"
	"sync"
)

var _ Store = (*memoryStore)(nil)

type memoryStore struct {
	memory map[string][]byte
	mtx    *sync.RWMutex
}

func NewMemoryStore() *memoryStore {
	return &memoryStore{
		memory: make(map[string][]byte),
		mtx:    new(sync.RWMutex),
	}
}

func (inv *memoryStore) Put(address []byte, data []byte) ([]byte, error) {
	inv.mtx.Lock()
	inv.memory[string(address)] = data
	inv.mtx.Unlock()
	return address, nil
}

func (inv *memoryStore) Delete(address []byte) error {
	inv.mtx.Lock()
	inv.memory[string(address)] = nil
	inv.mtx.Unlock()
	return nil
}

func (inv *memoryStore) Get(address []byte) ([]byte, error) {
	data, exists := inv.get(address)
	if !exists {
		return nil, ErrorAddressNotFound(address)
	}
	return data, nil
}

func (inv *memoryStore) Stat(address []byte) (*StatInfo, error) {
	data, exists := inv.get(address)
	return &StatInfo{
		Exists: exists,
		Size_:  uint64(len(data)),
	}, nil
}

func (inv *memoryStore) Location(address []byte) string {
	return fmt.Sprintf("memfs://%x", address)
}

func (inv *memoryStore) Name() string {
	return "memoryStore"
}

func (inv *memoryStore) get(address []byte) ([]byte, bool) {
	inv.mtx.RLock()
	data, exists := inv.memory[string(address)]
	inv.mtx.RUnlock()
	return data, exists
}
