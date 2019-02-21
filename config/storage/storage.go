package storage

import (
	"fmt"

	"bytes"

	"errors"

	"github.com/BurntSushi/toml"
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
	AWS         StorageType = "aws"
	Azure       StorageType = "azure"
	GCP         StorageType = "gcp"
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
	*CloudConfig
	*IPFSConfig
}

func NewStorageConfig(storageType StorageType, addressEncoding string) *StorageConfig {
	return &StorageConfig{
		StorageType:     storageType,
		AddressEncoding: addressEncoding,
	}
}

func GetDefaultConfig(c string) (*StorageConfig, error) {

	switch StorageType(c) {
	case Memory, Unspecified:
		return DefaultMemoryConfig(), nil
	case Filesystem:
		return DefaultFileSystemConfig(), nil
	case IPFS:
		return DefaultIPFSConfig(), nil
	case AWS:
		return DefaultCloudConfig(c), nil
	case Azure:
		return DefaultCloudConfig(c), nil
	case GCP:
		return DefaultCloudConfig(c), nil
	default:
		return nil, fmt.Errorf("did not recognise storage type '%s'", c)
	}
}

func StoreFromStorageConfig(storageConfig *StorageConfig, logger log.Logger) (storage.NamedStore, error) {
	addressEncoding, err := storage.GetAddressEncoding(storageConfig.AddressEncoding)
	if err != nil {
		return nil, err
	}

	switch storageConfig.StorageType {
	case Memory, Unspecified:
		return storage.NewMemoryStore(), nil

	case Filesystem:
		fsConf := storageConfig.FileSystemConfig
		if fsConf == nil {
			return nil, errors.New("filesystem storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if fsConf.RootDirectory == "" {
			return nil, errors.New("rootDirectory key must be non-empty in " +
				"filesystem storage config")
		}
		return storage.NewFileSystemStore(fsConf.RootDirectory, addressEncoding)

	case IPFS:
		ipfsConf := storageConfig.IPFSConfig
		if ipfsConf == nil {
			return nil, errors.New("IPFS storage configuration must be " +
				"supplied to use the filesystem storage backend")
		}
		if ipfsConf.RemoteAPI == "" {
			return nil, errors.New("http api url must be non-empty in " +
				"ipfs storage config")
		}
		return storage.NewIPFSStore(ipfsConf.RemoteAPI, addressEncoding)

	case AWS:
		awsConf := storageConfig.CloudConfig
		if awsConf == nil {
			return nil, errors.New("aws configuration must be supplied")
		}
		return storage.NewCloudStore(storage.CloudType(AWS), awsConf.Bucket, awsConf.Prefix, awsConf.Region, addressEncoding, logger)

	case Azure:
		azureConf := storageConfig.CloudConfig
		if azureConf == nil {
			return nil, errors.New("azure configuration must be supplied")
		}
		return storage.NewCloudStore(storage.CloudType(Azure), azureConf.Bucket, azureConf.Prefix, azureConf.Region, addressEncoding, logger)

	case GCP:
		gcpConf := storageConfig.CloudConfig
		if gcpConf == nil {
			return nil, errors.New("gcp configuration must be supplied")
		}
		return storage.NewCloudStore(storage.CloudType(GCP), gcpConf.Bucket, gcpConf.Prefix, gcpConf.Region, addressEncoding, logger)

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
