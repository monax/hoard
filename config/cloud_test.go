package config

import "testing"

func TestDefaultCloudConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, NewDefaultCloud("aws"))
	assertStorageConfigSerialisation(t, NewDefaultCloud("azure"))
	assertStorageConfigSerialisation(t, NewDefaultCloud("gcp"))
}
