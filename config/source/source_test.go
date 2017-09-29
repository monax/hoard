package source

import (
	"os"
	"testing"

	"io/ioutil"

	"github.com/monax/hoard/config"
	"github.com/monax/hoard/config/logging"
	"github.com/monax/hoard/config/storage"
	"github.com/stretchr/testify/assert"
)

func TestEnvironment(t *testing.T) {
	jsonString := config.DefaultHoardConfig.JSONString()
	os.Setenv(DefaultJSONConfigEnvironmentVariable, jsonString)
	conf, err := Environment(DefaultJSONConfigEnvironmentVariable).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, jsonString, conf.JSONString())
}

func TestFile(t *testing.T) {
	tomlString := config.DefaultHoardConfig.TOMLString()
	file := writeConfigFile(t, config.DefaultHoardConfig)
	defer os.Remove(file)
	conf, err := File(file).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, tomlString, conf.TOMLString())
}

func TestCascade(t *testing.T) {
	// Both fall through so baseConfig returned
	conf, err := Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File("")).Get(config.DefaultHoardConfig)
	assert.NoError(t, err)
	assert.Equal(t, *config.DefaultHoardConfig, *conf)

	// Env not set so falls through to file
	fileConfig := config.DefaultHoardConfig
	file := writeConfigFile(t, fileConfig)
	defer os.Remove(file)
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(&config.HoardConfig{})
	assert.NoError(t, err)
	assert.Equal(t, fileConfig.TOMLString(), conf.TOMLString())

	// Env set so caught by environment source
	envConfig := config.NewHoardConfig("unix:///tmp/hoard.sock'", storage.DefaultS3Config(),
		logging.DefaultConfig)
	os.Setenv(DefaultJSONConfigEnvironmentVariable, envConfig.JSONString())
	conf, err = Cascade(os.Stderr, true,
		Environment(DefaultJSONConfigEnvironmentVariable),
		File(file)).Get(config.DefaultHoardConfig)
	assert.NoError(t, err)
	assert.Equal(t, envConfig.TOMLString(), conf.TOMLString())
}

func writeConfigFile(t *testing.T, hoardConfig *config.HoardConfig) string {
	tomlString := hoardConfig.TOMLString()
	f, err := ioutil.TempFile("", DefaultHoardConfigFileName)
	assert.NoError(t, err)
	f.Write(([]byte)(tomlString))
	f.Close()
	return f.Name()
}
