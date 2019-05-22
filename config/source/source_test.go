package source

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/monax/hoard/v4/config"
	"github.com/monax/hoard/v4/config/logging"
	"github.com/monax/hoard/v4/config/storage"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	jsonString := config.DefaultHoardConfig.JSONString()
	os.Setenv(DefaultJSONConfigEnvironmentVariable, jsonString)
	conf, err := Environment(DefaultJSONConfigEnvironmentVariable).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, jsonString, conf.JSONString())
}

func TestTOMLFile(t *testing.T) {
	tomlString := config.DefaultHoardConfig.TOMLString()
	file := writeConfigFile(t, tomlString)
	defer os.Remove(file)
	conf, err := File(file).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, tomlString, conf.TOMLString())
}

func TestYAMLFile(t *testing.T) {
	yamlString := config.DefaultHoardConfig.YAMLString()
	file := writeConfigFile(t, yamlString)
	defer os.Remove(file)
	conf, err := File(file).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, yamlString, conf.YAMLString())
}

func TestCascade(t *testing.T) {
	// Both fall through so baseConfig returned
	conf, err := Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File("")).Get(config.DefaultHoardConfig)
	assert.NoError(t, err)
	assert.Equal(t, *config.DefaultHoardConfig, *conf)

	// Env not set so falls through to file
	fileConfig := config.DefaultHoardConfig.TOMLString()
	file := writeConfigFile(t, fileConfig)
	defer os.Remove(file)
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, fileConfig, conf.TOMLString())

	// Env set so caught by environment source
	envConfig := config.NewHoardConfig("unix:///tmp/hoard.sock'", storage.DefaultCloudConfig("aws"),
		logging.DefaultConfig)
	os.Setenv(DefaultJSONConfigEnvironmentVariable, envConfig.JSONString())
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(config.DefaultHoardConfig)
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
