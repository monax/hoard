package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v8/cmd"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/server"
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

	environmentOpt := hoardApp.BoolOpt("env-config", false,
		fmt.Sprintf("Parse the contents of the environment variable %s as a complete JSON config",
			config.DefaultJSONConfigEnvironmentVariable))

	secretsFromEnv := hoardApp.BoolOpt("env-secrets", false,
		fmt.Sprintf("Decode the environment variables pointed to by the symmetric public IDs."))

	// This string spec is parsed by mow.cli and has actual semantic significance
	// around optionality and ordering of options and arguments
	hoardApp.Spec = "[--config=<path to config file> | --env-config] [--address=<address to listen on>] [--logging] [--env-secrets]"

	cmd.AddVersionCommand(hoardApp)

	hoardApp.Action = func() {
		conf, err := hoardConfigCascade(*environmentOpt, *configFileOpt).Get(nil)
		if err != nil {
			fatalf("Could not get Hoard config: %s", err)
		}

		// I can't think of a good reason to allow this...
		if conf.ChunkSize == 0 {
			conf.ChunkSize = config.DefaultChunkSize
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

		symmetricProvider, err := config.NewSymmetricProvider(conf.Secrets, *secretsFromEnv)
		if err != nil {
			fatalf("Could not load symmetric keys: %s", err)
		}
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
		"examples demonstrating some features and need to be edited.", Config)

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
