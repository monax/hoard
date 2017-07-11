package storage

import "testing"

func TestDefaultMemoryConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultMemoryConfig())
}
