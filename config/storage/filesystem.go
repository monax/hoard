package storage

import (
	"fmt"
	"path"

	"github.com/cep21/xdgbasedir"
)

type FileSystemConfig struct {
	RootDirectory string
}

func NewFileSystemConfig(addressEncoding, rootDirectory string) *StorageConfig {
	return &StorageConfig{
		StorageType:     Filesystem,
		AddressEncoding: addressEncoding,
		FileSystemConfig: &FileSystemConfig{
			RootDirectory: rootDirectory,
		},
	}
}

func DefaultFileSystemConfig() *StorageConfig {
	dataDir, err := xdgbasedir.DataHomeDirectory()
	if err != nil {
		panic(fmt.Errorf("Could not get XDG data dir: %s", err))
	}
	return NewFileSystemConfig(DefaultAddressEncodingName,
		path.Join(dataDir, "hoard"))
}
