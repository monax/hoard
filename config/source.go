package config

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/cep21/xdgbasedir"
)

const DefaultHoardConfigFileName = "hoard.conf"
const DefaultJSONConfigEnvironmentVariable = "HOARD_JSON_CONFIG"
const STDINFileIdentifier = "-"

type Provider interface {
	// Description of where this provider sources its config from
	From() string
	// Get the config possibly overriding values passed in from baseConfig
	Get(baseConfig *HoardConfig) (*HoardConfig, error)
	// Return a copy of the provider that does nothing if skip is true
	SetSkip(skip bool) Provider
	// Whether to skip this provider
	Skip() bool
}

var _ Provider = &configSource{}

type configSource struct {
	from     string
	skip     bool
	provider func(baseConfig *HoardConfig) (*HoardConfig, error)
}

func (cs *configSource) From() string {
	return cs.from
}

func (cs *configSource) Get(baseConfig *HoardConfig) (*HoardConfig, error) {
	return cs.provider(baseConfig)
}

func (cs *configSource) Skip() bool {
	return cs.skip
}

// Returns a copy of the configSource with skip set as passed in
func (cs *configSource) SetSkip(skip bool) Provider {
	return &configSource{
		skip:     skip,
		from:     cs.from,
		provider: cs.provider,
	}
}

func Cascade(logWriter io.Writer, shortCircuit bool, sources ...Provider) *configSource {
	var fromStrings []string
	for _, source := range sources {
		if !source.Skip() {
			fromStrings = append(fromStrings, source.From())
		}
	}
	return &configSource{
		from: strings.Join(fromStrings, " then "),
		provider: func(baseConfig *HoardConfig) (*HoardConfig, error) {
			for _, source := range sources {
				if !source.Skip() {
					writeLog(logWriter, fmt.Sprintf("Trying to source config from %s", source.From()))
					conf, err := source.Get(baseConfig)
					if err != nil {
						return nil, err
					}
					if conf != nil {
						if shortCircuit {
							writeLog(logWriter, fmt.Sprintf("Using config from %s", source.From()))
							return conf, nil
						}
						writeLog(logWriter,
							fmt.Sprintf("Using config from %s and checking next source for overrides",
								source.From()))
						baseConfig = conf
					}
				}
			}
			if baseConfig == nil {
				return nil, errors.New("config cascade could not establish a config")
			}
			return baseConfig, nil
		},
	}
}

// Source from file
func File(configFile string) *configSource {
	return &configSource{
		skip: configFile == "",
		from: fmt.Sprintf("Storage file at '%s'", configFile),
		provider: func(baseConfig *HoardConfig) (*HoardConfig, error) {
			return fromFile(configFile)
		},
	}
}

// Try to find config by using XDG base dir spec
func XDGBaseDir() *configSource {
	skip := false
	// Look for config in standard XDG specified locations
	configFile, err := xdgbasedir.GetConfigFileLocation(DefaultHoardConfigFileName)
	if err == nil {
		_, err := os.Stat(configFile)
		// Skip if config  file does not exist at default location
		skip = os.IsNotExist(err)
	}
	return &configSource{
		skip: skip,
		from: fmt.Sprintf("XDG base dir"),
		provider: func(baseConfig *HoardConfig) (*HoardConfig, error) {
			if err != nil {
				return nil, err
			}
			return fromFile(configFile)
		},
	}
}

// Source from a single environment variable with config embedded in JSON
func Environment(key string) *configSource {
	jsonString := os.Getenv(key)
	return &configSource{
		skip: jsonString == "",
		from: fmt.Sprintf("'%s' environment variable (as JSON)", key),
		provider: func(baseConfig *HoardConfig) (*HoardConfig, error) {
			conf, err := HoardConfigFromJSONString(jsonString)
			if err != nil {
				return nil, err
			}
			return conf, nil
		},
	}
}

func Default() *configSource {
	return &configSource{
		from: "defaults",
		provider: func(baseConfig *HoardConfig) (*HoardConfig, error) {
			return DefaultHoardConfig, nil
		},
	}
}

func fromFile(configFile string) (*HoardConfig, error) {
	bs, err := readFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("could not read config file '%s': %s",
			configFile, err)
	}
	if len(bs) == 0 {
		return nil, fmt.Errorf("empty config")
	}

	configSpec := string(bs)
	if hoardConf, err := HoardConfigFromTOMLString(configSpec); err == nil {
		return hoardConf, err
	}
	return HoardConfigFromYAMLString(configSpec)
}

func readFile(configFile string) ([]byte, error) {
	if configFile == STDINFileIdentifier {
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(configFile)
}

func writeLog(writer io.Writer, msg string) {
	if writer != nil {
		writer.Write(([]byte)(msg))
		writer.Write(([]byte)("\n"))
	}
}
