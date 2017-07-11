package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHoardConfig(t *testing.T) {
	assertHoardConfigSerialisation(t, DefaultHoardConfig)
}

func assertHoardConfigSerialisation(t *testing.T,
	hoardConfig *HoardConfig) {

	storageConfigOut, err := HoardConfigFromString(hoardConfig.TOMLString())
	assert.NoError(t, err)
	tomlString := hoardConfig.TOMLString()
	assert.NotEmpty(t, tomlString)
	assert.Equal(t, tomlString, storageConfigOut.TOMLString())
}
