package storage

import (
	"fmt"

	"bytes"

	"errors"

	"code.monax.io/platform/hoard/core/storage"
	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-kit/kit/log"
)

const DefaultAddressEncodingName = "base64"

var DefaultConfig = NewMemoryConfig(DefaultAddressEncodingName)

type StorageType string

const (
	Unspecified StorageType = ""
	Memory      StorageType = "memory"
	Filesystem  StorageType = "filesystem"
	S3          StorageType = "s3"
	IPFS        StorageType = "ipfs"
)

type StorageConfig struct {
	// Acts a string enum
	StorageType StorageType
	// Address encoding name
	AddressEncoding string
	// Embedding a pointer to each type of config struct allows us to access the
	// relevant one, while at the same time those that are left as nil will be
	// omitted from being serialised.
	*FileSystemConfig
	*S3Config
	*IPFSConfig
}

func NewStorageConfig(storageType StorageType, addressEncoding string) *StorageConfig {
	return &StorageConfig{
		StorageType:     storageType,
		AddressEncoding: addressEncoding,
	}
}

func StoreFromStorageConfig(storageConfig *StorageConfig,
	logger log.Logger) (storage.Store, error) {

	addressEncoding, err := storage.GetAddressEncoding(storageConfig.AddressEncoding)
	if err != nil {
		return nil, err
	}

	switch storageConfig.StorageType {
	case Memory, Unspecified:
		return storage.NewMemoryStore(), nil
	case Filesystem:
		fsc := storageConfig.FileSystemConfig
		if fsc == nil {
			return nil, errors.New("Filesystem storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if fsc.RootDirectory == "" {
			return nil, errors.New("RootDirectory key must be non-empty in " +
				"filesystem storage config.")
		}
		return storage.NewFileSystemStore(fsc.RootDirectory, addressEncoding), nil
	case S3:
		s3c := storageConfig.S3Config
		if s3c == nil {
			return nil, errors.New("S3 configuration must be supplied to use " +
				"the S3 storage backend")
		}

		creds, err := AWSCredentialsFromChain(s3c.CredentialsProviderChain)
		if err != nil {
			return nil, fmt.Errorf("Could not create credentials: %s", err)
		}

		var region *string
		if s3c.Region != "" {
			region = aws.String(s3c.Region)
		}

		awsConfig := &aws.Config{
			Credentials: creds,
			Region:      region,
			Logger: aws.LoggerFunc(func(keyvals ...interface{}) {
				logger.Log(keyvals...)
			}),
		}

		return storage.NewS3Store(s3c.Bucket, s3c.Prefix, addressEncoding,
			awsConfig, logger)
	default:
		return nil, fmt.Errorf("Did not recognise storage type '%s'",
			storageConfig.StorageType)
	}
}

func ConfigFromString(tomlString string) (*StorageConfig, error) {
	storageConfig := new(StorageConfig)
	_, err := toml.Decode(tomlString, storageConfig)
	if err != nil {
		return nil, err
	}
	return storageConfig, nil
}

func (storageConfig *StorageConfig) TOMLString() string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(storageConfig)
	if err != nil {
		return fmt.Sprintf("<Could not serialise StorageConfig>")
	}
	return buf.String()
}
