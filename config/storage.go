package config

import (
	"fmt"

	"github.com/monax/hoard/v6/stores"

	"bytes"

	"github.com/BurntSushi/toml"
)

const DefaultAddressEncodingName = stores.Base64EncodingName

func NewDefaultStorage() *Storage {
	return NewStorage(Memory, DefaultAddressEncodingName)
}

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

// Storage identifies the configured back-end
type Storage struct {
	// Acts a string enum
	StorageType StorageType
	// Address encoding name
	AddressEncoding string
	// Embedding a pointer to each type of config struct allows us to access the
	// relevant one, while at the same time those that are left as nil will be
	// omitted from being serialised.
	*FileSystemConfig
	*Cloud
	*IPFSConfig
}

func NewStorage(storageType StorageType, addressEncoding string) *Storage {
	return &Storage{
		StorageType:     storageType,
		AddressEncoding: addressEncoding,
	}
}

func GetStorageTypes() []StorageType {
	return []StorageType{
		Memory,
		Filesystem,
		AWS,
		Azure,
		GCP,
		IPFS,
	}
}

func GetDefaultStorage(storageType StorageType) (*Storage, error) {
	switch storageType {
	case Memory, Unspecified:
		return NewDefaultMemory(), nil
	case Filesystem:
		return NewDefaultFileSystemConfig(), nil
	case IPFS:
		return NewDefaultIPFSConfig(), nil
	case AWS:
		return NewDefaultCloud(storageType), nil
	case Azure:
		return NewDefaultCloud(storageType), nil
	case GCP:
		return NewDefaultCloud(storageType), nil
	default:
		return nil, fmt.Errorf("did not recognise storage type '%s'", storageType)
	}
}

func ConfigFromString(tomlString string) (*Storage, error) {
	storageConfig := new(Storage)
	_, err := toml.Decode(tomlString, storageConfig)
	if err != nil {
		return nil, err
	}
	return storageConfig, nil
}

func (storageConfig *Storage) TOMLString() string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	err := encoder.Encode(storageConfig)
	if err != nil {
		return fmt.Sprintf("<Could not serialise Storage>")
	}
	return buf.String()
}
