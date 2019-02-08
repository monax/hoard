package storage

import "testing"

func TestDefaultCloudConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultCloudConfig("aws"))
	assertStorageConfigSerialisation(t, DefaultCloudConfig("azure"))
	assertStorageConfigSerialisation(t, DefaultCloudConfig("gcp"))
}
