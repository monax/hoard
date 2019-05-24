package config

import "testing"

func TestDefaultIPFSConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultIPFSConfig())
}
