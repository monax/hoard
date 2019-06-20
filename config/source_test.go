package config

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	jsonString := DefaultHoardConfig.JSONString()
	os.Setenv(DefaultJSONConfigEnvironmentVariable, jsonString)
	conf, err := Environment(DefaultJSONConfigEnvironmentVariable).Get(&HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, jsonString, conf.JSONString())
}

func TestTOMLFile(t *testing.T) {
	tomlString := DefaultHoardConfig.TOMLString()
	file := writeConfigFile(t, tomlString)
	defer os.Remove(file)
	conf, err := File(file).Get(&HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, tomlString, conf.TOMLString())
}

func TestYAMLFile(t *testing.T) {
	yamlString := DefaultHoardConfig.YAMLString()
	file := writeConfigFile(t, yamlString)
	defer os.Remove(file)
	conf, err := File(file).Get(&HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, yamlString, conf.YAMLString())
}

func TestCascade(t *testing.T) {
	// Both fall through so baseConfig returned
	conf, err := Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File("")).Get(DefaultHoardConfig)
	assert.NoError(t, err)
	assert.Equal(t, *DefaultHoardConfig, *conf)

	// Env not set so falls through to file
	fileConfig := DefaultHoardConfig.TOMLString()
	file := writeConfigFile(t, fileConfig)
	defer os.Remove(file)
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(&HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, fileConfig, conf.TOMLString())

	// Env set so caught by environment source
	envConfig := NewHoardConfig("unix:///tmp/hoard.sock'", DefaultCloud("aws"), DefaultLogging)
	os.Setenv(DefaultJSONConfigEnvironmentVariable, envConfig.JSONString())
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(DefaultHoardConfig)
	assert.NoError(t, err)
	assert.Equal(t, envConfig.TOMLString(), conf.TOMLString())
}

func writeConfigFile(t *testing.T, hoardConfig string) string {
	f, err := ioutil.TempFile("", DefaultHoardConfigFileName)
	assert.NoError(t, err)
	f.Write(([]byte)(hoardConfig))
	f.Close()
	return f.Name()
}
