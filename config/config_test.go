package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultHoardConfig(t *testing.T) {
	assertHoardConfigSerialisation(t,
		func(conf *HoardConfig) string {
			return conf.TOMLString()
		},
		HoardConfigFromTOMLString,
		DefaultHoardConfig)
	assertHoardConfigSerialisation(t,
		func(conf *HoardConfig) string {
			return conf.JSONString()
		},
		HoardConfigFromJSONString,
		DefaultHoardConfig)
}

func assertHoardConfigSerialisation(t *testing.T,
	serialise func(*HoardConfig) string,
	deserialise func(string) (*HoardConfig, error),
	hoardConfig *HoardConfig) {

	hoardConfigRoundTrip, err := deserialise(serialise(hoardConfig))
	assert.NoError(t, err)
	hoardConfigString := serialise(hoardConfig)
	assert.NotEmpty(t, hoardConfigString)
	assert.Equal(t, hoardConfigString, serialise(hoardConfigRoundTrip))
}
