package main

import (
	"fmt"
	"strings"

	"github.com/cep21/xdgbasedir"
	cli "github.com/jawher/mow.cli"
	"github.com/monax/hoard/v8/config"
	"github.com/monax/hoard/v8/encryption"
)

func Config(cmd *cli.Cmd) {
	conf := config.DefaultHoardConfig

	outputOpt := cmd.StringOpt("o output", "",
		"Instead of writing to STDOUT, write to the specified file")

	overwriteOpt := cmd.BoolOpt("f force", false,
		"Overwrite config file if it exists.")

	jsonOpt := cmd.BoolOpt("j json", false,
		"Print config to STDOUT as a single line of JSON (suitable for the --env config source)")

	yamlOpt := cmd.BoolOpt("y yaml", false,
		"Print config to STDOUT as a YAML specification")

	initOpt := cmd.BoolOpt("i init", false, "Write file to "+
		"XDG standard location")

	chunkSizeOpt := cmd.IntOpt("c chunk-size", config.DefaultChunkSize,
		"number of bytes on which to split plaintext/ciphertext across message boundaries when streaming")

	secretsOpt := cmd.StringsOpt("s secret", nil, "Pairs of PublicID and Passphrase to use as symmetric secrets in config")

	arg := cmd.StringArg("CONFIG", "", fmt.Sprintf("Storage type to generate, one of: %s",
		strings.Join(configTypes(), ", ")))

	cmd.Spec = "[--json | --yaml] | (([--output=<output file>] |  [--init]) [--force]) CONFIG " +
		"[--secret=<PublicID:Passphrase>...] [--chunk-size=<message chunk size in bytes>]"

	cmd.Action = func() {
		store, err := config.GetDefaultStorage(config.StorageType(*arg))
		if err != nil {
			fatalf("Error fetching default config for %v: %v", arg, err)
		}
		conf.Storage = store
		conf.ChunkSize = int64(*chunkSizeOpt)
		if len(*secretsOpt) > 0 {
			conf.Secrets = &config.Secrets{
				Symmetric: make([]*config.SymmetricSecret, len(*secretsOpt)),
			}
			for i, ss := range *secretsOpt {
				pair := strings.Split(ss, ":")
				if len(pair) != 2 {
					fatalf("got symmetric secret specification '%s' but must be specified as <PublicID:Passphrase>", ss)
				}
				salt, err := encryption.NewNonce(encryption.NonceSize)
				if err != nil {
					fatalf("failed to generate salt: %v", err)
				}

				data, err := encryption.DeriveSecretKey([]byte(pair[1]), salt)
				if err != nil {
					fatalf("could not derive secret key for config: %v", err)
				}

				conf.Secrets.Symmetric[i] = new(config.SymmetricSecret)
				conf.Secrets.Symmetric[i].PublicID = pair[0]
				conf.Secrets.Symmetric[i].SecretKey = data
			}
		}
	}

	cmd.After = func() {
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
}
