package config

func NewMemory(addressEncoding string) *Storage {
	return NewStorage(Memory, addressEncoding)
}

func DefaultMemory() *Storage {
	return NewMemory(DefaultAddressEncodingName)
}
