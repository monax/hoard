package config

import "testing"

func TestDefaultCloudConfig(t *testing.T) {
	assertStorageConfigSerialisation(t, DefaultCloud("aws"))
	assertStorageConfigSerialisation(t, DefaultCloud("azure"))
	assertStorageConfigSerialisation(t, DefaultCloud("gcp"))
}
