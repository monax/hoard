package storage

import "testing"

func TestDefaultS3Config(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultS3Config())
}
