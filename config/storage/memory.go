package storage

func NewMemoryConfig(addressEncoding string) *StorageConfig {
	return NewStorageConfig(Memory, addressEncoding)
}
