package config

import "testing"

func TestDefaultFileSystemConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, NewDefaultFileSystemConfig())
}
