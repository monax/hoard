package storage

import "testing"

func TestDefaultFileSystemConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultFileSystemConfig())
}
