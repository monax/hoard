package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/cep21/xdgbasedir"
	"github.com/go-kit/kit/log"
	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v7/cmd"
	"github.com/monax/hoard/v7/config"
	"github.com/monax/hoard/v7/server"
)

func main() {
	hoardApp := cli.App("hoard",
		"A content-addressed deterministically encrypted blob storage system")

	listenAddressOpt := hoardApp.StringOpt("a address", "",
		"local address for hoard to listen on encoded as a URL with the "+
			"network protocol as the scheme, for example 'tcp://localhost:54192' "+
			"or 'unix:///tmp/hoard.sock'")

	loggingOpt := hoardApp.BoolOpt("l logging", false,
		"Whether to emit any operational logging")

	configFileOpt := hoardApp.StringOpt("c config", "", "Path to "+
		"config file. If omitted default config is used. Use '-' to read config from STDIN.")

	environmentOpt := hoardApp.BoolOpt("e env", false,
		fmt.Sprintf("Parse the contents of the environment variable %s as a complete JSON config",
			config.DefaultJSONConfigEnvironmentVariable))

	// This string spec is parsed by mow.cli and has actual semantic significance
	// around optionality and ordering of options and arguments
	hoardApp.Spec = "[--config=<path to config file> | --env] [--address=<address to listen on>] [--logging]"

	cmd.AddVersionCommand(hoardApp)

	hoardApp.Action = func() {
		conf, err := hoardConfigCascade(*environmentOpt, *configFileOpt).Get(nil)
		if err != nil {
			fatalf("Could not get Hoard config: %s", err)
		}

		var logger log.Logger

		if *loggingOpt {
			logger, err = config.Logger(conf.Logging, os.Stderr)
			if err != nil {
				fatalf("Could not create logging form logging config: %s", err)
			}
		}

		store, err := StoreFromStorageConfig(conf.Storage, logger)
		if err != nil {
			fatalf("Could not configure store from storage config: %s", err)
		}
		if *listenAddressOpt != "" {
			conf.ListenAddress = *listenAddressOpt
		}
		symmetricProvider := config.NewSymmetricProvider(conf.Secrets)
		openPGPConf := config.NewOpenPGPSecret(conf.Secrets)
		secretsManager := config.SecretsManager{Provider: symmetricProvider, OpenPGP: openPGPConf}
		serv := server.New(conf.ListenAddress, store, secretsManager, conf.ChunkSize, logger)
		// Catch interrupt etc
		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt, os.Kill, syscall.SIGTERM)
		go func(c chan os.Signal) {
			sig := <-c
			printf("\nCaught %s signal: shutting down...", sig)
			// Make sure we clean up
			serv.Stop()
			os.Exit(0)
		}(signalCh)

		printf("Starting hoard daemon on %s with chunk size %d on %s...", conf.ListenAddress, conf.ChunkSize,
			store.Name())
		err = serv.Serve()
		if err != nil {
			fatalf("Could not start hoard server: %s", err)
		}
	}

	hoardApp.Command("config", "Initialise Hoard configuration by "+
		"printing an example configuration file to STDOUT. Most config files emitted are "+
		"examples demonstrating some features and need to be edited.",
		func(configCmd *cli.Cmd) {
			conf := config.DefaultHoardConfig

			outputOpt := configCmd.StringOpt("o output", "",
				"Instead of writing to STDOUT, write to the specified file")

			overwriteOpt := configCmd.BoolOpt("f force", false,
				"Overwrite config file if it exists.")

			jsonOpt := configCmd.BoolOpt("j json", false,
				"Print config to STDOUT as a single line of JSON (suitable for the --env config source)")

			yamlOpt := configCmd.BoolOpt("y yaml", false,
				"Print config to STDOUT as a YAML specification")

			initOpt := configCmd.BoolOpt("i init", false, "Write file to "+
				"XDG standard location")

			chunkSizeOpt := configCmd.IntOpt("c chunk-size", config.DefaultChunkSize,
				"number of bytes on which to split plaintext/ciphertext across message boundaries when streaming")

			secretsOpt := configCmd.StringsOpt("s secret", nil, "Pairs of PublicID and Passphrase to use as symmetric secrets in config")

			arg := configCmd.StringArg("CONFIG", "", fmt.Sprintf("Storage type to generate, one of: %s",
				strings.Join(configTypes(), ", ")))

			configCmd.Spec = "[--json | --yaml] | (([--output=<output file>] |  [--init]) [--force]) CONFIG " +
				"[--secret=<PublicID:Passphrase>...] [--chunk-size=<message chunk size in bytes>]"

			configCmd.Action = func() {
				store, err := config.GetDefaultStorage(config.StorageType(*arg))
				if err != nil {
					fatalf("Error fetching default config for %v: %v", arg, err)
				}
				conf.Storage = store
				conf.ChunkSize = *chunkSizeOpt
				if len(*secretsOpt) > 0 {
					conf.Secrets = &config.Secrets{
						Symmetric: make([]config.SymmetricSecret, len(*secretsOpt)),
					}
					for i, ss := range *secretsOpt {
						pair := strings.Split(ss, ":")
						if len(pair) != 2 {
							fatalf("got symmetric secret specification '%s' but must be specified as <PublicID:Passphrase>", ss)
						}
						conf.Secrets.Symmetric[i].PublicID = pair[0]
						conf.Secrets.Symmetric[i].Passphrase = pair[1]
					}
				}
			}

			configCmd.After = func() {
				configString := conf.TOMLString()
				if *jsonOpt {
					configString = conf.JSONString()
				} else if *yamlOpt {
					configString = conf.YAMLString()
				}
				if *initOpt {
					configFileName, err := xdgbasedir.GetConfigFileLocation(
						config.DefaultHoardConfigFileName)
					if err != nil {
						fatalf("Error getting config file location: %s", err)
					}
					outputOpt = &configFileName
				}
				if *outputOpt != "" {
					printf("Writing to config file '%s'", *outputOpt)
					err := writeFile(*outputOpt, ([]byte)(configString), *overwriteOpt)
					if err != nil {
						fatalf("Error writing config file: %s", err)
					}
				} else {
					fmt.Print(configString)
				}
			}
		})

	hoardApp.Run(os.Args)
}

// Print informational output to Stderr
func printf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func hoardConfigCascade(env bool, configFile string) config.Provider {
	return config.Cascade(os.Stderr, true,
		config.Environment(config.DefaultJSONConfigEnvironmentVariable).SetSkip(!env),
		config.File(configFile).SetSkip(configFile == ""),
		config.XDGBaseDir(),
		config.Default())
}

func writeFile(filename string, data []byte, overwrite bool) error {
	if _, err := os.Stat(filename); overwrite || os.IsNotExist(err) {
		return ioutil.WriteFile(filename, data, 0666)
	}
	return fmt.Errorf("file '%s' already exists", filename)
}

func configTypes() []string {
	storageTypes := config.GetStorageTypes()
	configTypes := make([]string, len(storageTypes))
	for i, st := range storageTypes {
		configTypes[i] = string(st)
	}
	return configTypes
}
