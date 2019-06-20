package config

import "testing"

func TestDefaultMemoryConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultMemory())
}
