package config

import (
	"testing"

	"code.monax.io/platform/hoard/config/storage"
	"github.com/stretchr/testify/assert"
)

func TestDefaultStorageConfig(t *testing.T) {
	// We are panicking on serialisation errors so check here
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Could not create default storage configs: %s", r)
		}
	}()

	assertStorageConfigSerialisation(t, DefaultMemoryConfig())
	assertStorageConfigSerialisation(t, DefaultFileSystemConfig())
	assertStorageConfigSerialisation(t, DefaultS3Config())
}

func assertStorageConfigSerialisation(t *testing.T,
	storageConfig *storage.StorageConfig) {

	storageConfigOut, err := storage.ConfigFromString(storageConfig.TOMLString())
	assert.NoError(t, err)
	tomlString := storageConfig.TOMLString()
	assert.NotEmpty(t, tomlString)
	assert.Equal(t, tomlString, storageConfigOut.TOMLString())
}
