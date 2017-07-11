package config

import (
	"fmt"

	"path"

	"os/user"

	"code.monax.io/platform/hoard/config/logging"
	"code.monax.io/platform/hoard/config/storage"
	"code.monax.io/platform/hoard/core/logging/structure"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/cep21/xdgbasedir"
)

const (
	DefaultListenAddress       = "tcp://localhost:53431"
	DefaultAddressEncodingName = "base64"
	DefaultHoardConfigFileName = "hoard.toml"
)

var (
	DefaultHoardConfig = NewHoardConfig(DefaultListenAddress,
		storage.NewMemoryConfig(DefaultAddressEncodingName),
		logging.NewLoggingConfig(logging.Logfmt, structure.InfoChannel,
			structure.TraceChannel))
)

func DefaultS3Config() *storage.StorageConfig {
	usr, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("Could not get home directory: %s", err))
	}
	s3c, err := storage.NewS3Config(DefaultAddressEncodingName,
		"monax-hoard-test",
		"store",
		"eu-central-1",
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{
			Filename: path.Join(usr.HomeDir, ".aws", "credentials"),
			Profile:  "default",
		},
		&credentials.StaticProvider{
			Value: credentials.Value{
				AccessKeyID:     "",
				SecretAccessKey: "",
				SessionToken:    "",
			},
		},
	)
	if err != nil {
		panic(fmt.Errorf("Could not generate example config: %s", err))
	}
	return s3c
}

func DefaultFileSystemConfig() *storage.StorageConfig {

	dataDir, err := xdgbasedir.DataHomeDirectory()
	if err != nil {
		panic(fmt.Errorf("Could not get XDG data dir: %s", err))
	}
	return storage.NewFileSystemConfig(DefaultAddressEncodingName,
		path.Join(dataDir, "hoard"))
}

func DefaultMemoryConfig() *storage.StorageConfig {
	return storage.NewMemoryConfig(DefaultAddressEncodingName)
}
