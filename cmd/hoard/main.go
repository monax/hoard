package main

import (
	"fmt"
	"os"

	"os/signal"
	"syscall"

	"io/ioutil"

	"github.com/cep21/xdgbasedir"
	"github.com/go-kit/kit/log"
	"github.com/jawher/mow.cli"
	"github.com/monax/hoard/cmd"
	"github.com/monax/hoard/config"
	"github.com/monax/hoard/config/logging"
	"github.com/monax/hoard/config/source"
	"github.com/monax/hoard/config/storage"
	"github.com/monax/hoard/server"
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
		"config file. If omitted default config is used.")

	environmentOpt := hoardApp.BoolOpt("e env", false,
		fmt.Sprintf("Parse the contents of the enironment variable %s as a complete JSON config",
			source.DefaultJSONConfigEnvironmentVariable))

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
			logger, err = logging.LoggerFromLoggingConfig(conf.Logging, os.Stderr)
			if err != nil {
				fatalf("Could not create logging form logging config: %s", err)
			}
		}

		store, err := storage.StoreFromStorageConfig(conf.Storage, logger)
		if err != nil {
			fatalf("Could not configure store from storage config: %s", err)
		}
		if *listenAddressOpt != "" {
			conf.ListenAddress = *listenAddressOpt
		}
		serv := server.New(conf.ListenAddress, store, logger)
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

		printf("Starting hoard daemon on %s with %s...", conf.ListenAddress,
			store.Name())
		err = serv.Serve()
		if err != nil {
			fatalf("Could not start hoard server: %s", err)
		}
	}

	hoardApp.Command("init", "Initialise Hoard configuration by "+
		"writing an example configuration file. Most config files emitted are "+
		"examples demonstrating some features and need to be edited.",
		func(initCmd *cli.Cmd) {
			conf := config.DefaultHoardConfig

			outputOpt := initCmd.StringOpt("o output", "",
				"Instead of writing to standard config location, write to the "+
					"specified file, pass '-' to write the config to STDOUT")

			overwriteOpt := initCmd.BoolOpt("f force", false,
				"Overwrite config file if it exists.")

			jsonOpt := initCmd.BoolOpt("j json", false,
				"Output as JSON (suitable for environment variable parsing)")

			initCmd.BoolOpt("t toml", false,
				"Output as TOML (the default)")

			initCmd.Spec = "[--output=<output file or '-' for STDOUT>] [--json | --toml] [--force]"

			initCmd.Command("mem", "Emit initial config with memory "+
				"storage backend.",
				func(s3Cmd *cli.Cmd) {
					s3Cmd.Action = func() {
						conf.Storage = storage.DefaultMemoryConfig()
					}
				})

			initCmd.Command("fs", "Emit initial config with "+
				"filesystem storage backend.",
				func(s3Cmd *cli.Cmd) {
					s3Cmd.Action = func() {
						conf.Storage = storage.DefaultFileSystemConfig()
					}
				})

			initCmd.Command("s3", "Emit initial config with S3 storage "+
				"backend.",
				func(s3Cmd *cli.Cmd) {
					s3Cmd.Action = func() {
						conf.Storage = storage.DefaultS3Config()
					}
				})

			initCmd.After = func() {
				configString := conf.TOMLString()
				if *jsonOpt {
					configString = conf.JSONString()
				}
				if *outputOpt == "-" {
					fmt.Print(configString)
				} else {
					if *outputOpt == "" {
						configFileName, err := xdgbasedir.GetConfigFileLocation(
							source.DefaultHoardConfigFileName)
						if err != nil {
							fatalf("Error getting config file location: %s", err)
						}
						outputOpt = &configFileName
					}
					printf("Writing to config file '%s'", *outputOpt)
					err := writeFile(*outputOpt, ([]byte)(configString), *overwriteOpt)
					if err != nil {
						fatalf("Error writing config file: %s", err)
					}
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

func hoardConfigCascade(env bool, configFile string) source.ConfigProvider {
	return source.Cascade(os.Stderr, true,
		source.Environment(source.DefaultJSONConfigEnvironmentVariable).SetSkip(!env),
		source.File(configFile).SetSkip(configFile == ""),
		source.XDGBaseDir(),
		source.Default())
}

func writeFile(filename string, data []byte, overwrite bool) error {
	if _, err := os.Stat(filename); overwrite || os.IsNotExist(err) {
		return ioutil.WriteFile(filename, data, 0666)
	}
	return fmt.Errorf("file '%s' already exists", filename)
}
