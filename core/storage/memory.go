package storage

import (
	"fmt"
	"sync"
)

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

func (ms *memoryStore) Put(address []byte, data []byte) ([]byte, error) {
	ms.mtx.Lock()
	ms.memory[string(address)] = data
	ms.mtx.Unlock()
	return address, nil
}

func (ms *memoryStore) Get(address []byte) ([]byte, error) {
	data, exists := ms.get(address)
	if !exists {
		return nil, ErrorAddressNotFound(address)
	}
	return data, nil
}

func (ms *memoryStore) Stat(address []byte) (*StatInfo, error) {
	data, exists := ms.get(address)
	return &StatInfo{
		Exists: exists,
		Size:   uint64(len(data)),
	}, nil
}

func (ms *memoryStore) Location(address []byte) string {
	return fmt.Sprintf("memfs://%x", address)
}

func (ms *memoryStore) Name() string {
	return "memoryStore"
}

func (ms *memoryStore) get(address []byte) ([]byte, bool) {
	ms.mtx.RLock()
	data, exists := ms.memory[string(address)]
	ms.mtx.RUnlock()
	return data, exists
}
