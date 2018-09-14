package storage

import "testing"

func TestDefaultGCSConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultGCSConfig())
}
