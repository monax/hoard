package config

import (
	"bytes"
	"fmt"

	"code.monax.io/platform/hoard/config/logging"
	"code.monax.io/platform/hoard/config/storage"
	"github.com/BurntSushi/toml"
)

type HoardConfig struct {
	ListenAddress string
	Storage       *storage.StorageConfig
	Logging       *logging.LoggingConfig
	// TODO: SecretsConfig - how to access bootstrapping secrets
}

func NewHoardConfig(listenAddress string, storageConfig *storage.StorageConfig,
	loggingConfig *logging.LoggingConfig) *HoardConfig {
	return &HoardConfig{
		ListenAddress: listenAddress,
		Storage:       storageConfig,
		Logging:       loggingConfig,
	}
}

func HoardConfigFromString(tomlString string) (*HoardConfig, error) {
	hoardConfig := new(HoardConfig)
	_, err := toml.Decode(tomlString, hoardConfig)
	if err != nil {
		return nil, err
	}
	if hoardConfig.ListenAddress == "" {
		hoardConfig.ListenAddress = DefaultListenAddress
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
