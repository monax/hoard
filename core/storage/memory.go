package storage

import "fmt"

type memoryStore struct {
	memory map[string][]byte
}

func NewMemoryStore() Store {
	return &memoryStore{
		memory: make(map[string][]byte),
	}
}

func (ms *memoryStore) Put(address, data []byte) error {
	ms.memory[string(address)] = data
	return nil
}

func (ms *memoryStore) Get(address []byte) ([]byte, error) {
	return ms.memory[string(address)], nil
}

func (ms *memoryStore) Stat(address []byte) (*StatInfo, error) {
	data, exists := ms.memory[string(address)]
	return &StatInfo{
		Exists: exists,
		Size:   uint64(len(data)),
	}, nil
}

func (ms *memoryStore) Location(address []byte) string {
	return fmt.Sprintf("memfs://%x", address)
}
