package storage

import (
	"fmt"

	"path"

	"os/user"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/cep21/xdgbasedir"
)

const DefaultAddressEncodingName = "base64"

func DefaultS3Config() *StorageConfig {
	usr, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("Could not get home directory: %s", err))
	}
	s3c, err := NewS3Config(DefaultAddressEncodingName,
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

func DefaultFileSystemConfig() *StorageConfig {

	dataDir, err := xdgbasedir.DataHomeDirectory()
	if err != nil {
		panic(fmt.Errorf("Could not get XDG data dir: %s", err))
	}
	return NewFileSystemConfig(DefaultAddressEncodingName,
		path.Join(dataDir, "hoard"))
}

func DefaultMemoryConfig() *StorageConfig {
	return NewMemoryConfig(DefaultAddressEncodingName)
}
