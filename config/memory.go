package config

func NewDefaultMemory() *Storage {
	conf := NewDefaultStorage()
	conf.StorageType = Memory
	return conf
}
