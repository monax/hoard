package storage

import (
	"fmt"

	"bytes"

	"errors"

	"github.com/BurntSushi/toml"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/storage"
)

const DefaultAddressEncodingName = storage.Base64EncodingName

var DefaultConfig = NewMemoryConfig(DefaultAddressEncodingName)

type StorageType string

const (
	Unspecified StorageType = ""
	Memory      StorageType = "memory"
	Filesystem  StorageType = "filesystem"
	S3          StorageType = "s3"
	GCS         StorageType = "gcs"
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
	*GCSConfig
	*IPFSConfig
}

func NewStorageConfig(storageType StorageType, addressEncoding string) *StorageConfig {
	return &StorageConfig{
		StorageType:     storageType,
		AddressEncoding: addressEncoding,
	}
}

func StoreFromStorageConfig(storageConfig *StorageConfig,
	logger log.Logger) (storage.NamedStore, error) {

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
			return nil, errors.New("filesystem storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if fsc.RootDirectory == "" {
			return nil, errors.New("rootDirectory key must be non-empty in " +
				"filesystem storage config")
		}
		return storage.NewFileSystemStore(fsc.RootDirectory, addressEncoding)

	case IPFS:
		ipfsc := storageConfig.IPFSConfig
		if ipfsc == nil {
			return nil, errors.New("IPFS storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if ipfsc.Protocol == "" {
			ipfsc.Protocol = "https://"
		}
		if ipfsc.Address == "" {
			return nil, errors.New("http api url must be non-empty in " +
				"ipfs storage config")
		}
		if ipfsc.Port == "" {
			return nil, errors.New("http api port must be non-empty in " +
				"ipfs storage config")
		}
		return storage.NewIPFSStore(ipfsc.Protocol, ipfsc.Address, ipfsc.Port, addressEncoding)

	case S3:
		s3c := storageConfig.S3Config
		if s3c == nil {
			return nil, errors.New("s3 configuration must be supplied to use " +
				"the S3 storage backend")
		}

		creds, err := AWSCredentialsFromChain(s3c.CredentialsProviderChain)
		if err != nil {
			return nil, fmt.Errorf("could not create credentials: %s", err)
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
		return storage.NewS3Store(s3c.S3Bucket, s3c.S3Prefix, addressEncoding,
			awsConfig, logger)

	case GCS:
		gcsc := storageConfig.GCSConfig
		if gcsc == nil {
			return nil, errors.New("gpc configuration must be supplied to use " +
				"the GPC storage backend")
		}
		return storage.NewGCSStore(gcsc.GCSBucket, gcsc.GCSPrefix, addressEncoding, logger)

	default:
		return nil, fmt.Errorf("did not recognise storage type '%s'",
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
