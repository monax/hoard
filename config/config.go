package config

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/monax/hoard/config/logging"
	"github.com/monax/hoard/config/secrets"
	"github.com/monax/hoard/config/storage"
)

const DefaultListenAddress = "tcp://:53431"

var DefaultHoardConfig = NewHoardConfig(DefaultListenAddress,
	storage.DefaultConfig, logging.DefaultConfig)

type HoardConfig struct {
	ListenAddress string
	Storage       *storage.StorageConfig
	Logging       *logging.LoggingConfig
	Secrets       *secrets.SecretsConfig
}

func NewHoardConfig(listenAddress string, storageConfig *storage.StorageConfig,
	loggingConfig *logging.LoggingConfig) *HoardConfig {
	return &HoardConfig{
		ListenAddress: listenAddress,
		Storage:       storageConfig,
		Logging:       loggingConfig,
	}
}

func HoardConfigFromJSONString(jsonString string) (*HoardConfig, error) {
	hoardConfig := new(HoardConfig)
	buf := bytes.NewBufferString(jsonString)
	decoder := json.NewDecoder(buf)
	err := decoder.Decode(hoardConfig)
	if err != nil {
		return nil, err
	}
	return hoardConfig, nil
}

func HoardConfigFromTOMLString(tomlString string) (*HoardConfig, error) {
	hoardConfig := new(HoardConfig)
	_, err := toml.Decode(tomlString, hoardConfig)
	if err != nil {
		return nil, err
	}
	return hoardConfig, nil
}

func (hoardConfig *HoardConfig) TOMLString() string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(hoardConfig)
	if err != nil {
		return fmt.Sprintf("<Could not serialise HoardConfig>")
	}
	return buf.String()
}

func (hoardConfig *HoardConfig) JSONString() string {
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err := encoder.Encode(hoardConfig)
	if err != nil {
		return fmt.Sprintf("<Could not serialise HoardConfig>")
	}
	return buf.String()
}
