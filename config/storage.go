package config

import (
	"fmt"

	"github.com/monax/hoard/v5/stores"

	"bytes"

	"github.com/BurntSushi/toml"
)

const DefaultAddressEncodingName = stores.Base64EncodingName

var DefaultStorage = NewMemory(DefaultAddressEncodingName)

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

func GetDefaultStorage(c string) (*Storage, error) {
	switch StorageType(c) {
	case Memory, Unspecified:
		return DefaultMemory(), nil
	case Filesystem:
		return DefaultFileSystemConfig(), nil
	case IPFS:
		return DefaultIPFSConfig(), nil
	case AWS:
		return DefaultCloud(c), nil
	case Azure:
		return DefaultCloud(c), nil
	case GCP:
		return DefaultCloud(c), nil
	default:
		return nil, fmt.Errorf("did not recognise storage type '%s'", c)
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
