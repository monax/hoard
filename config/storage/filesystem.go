package storage

import (
	"fmt"
	"path"

	"github.com/cep21/xdgbasedir"
	"github.com/monax/hoard/storage"
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
		panic(fmt.Errorf("could not get XDG data dir: %s", err))
	}
	// Avoid '/' character for filesystem storage
	return NewFileSystemConfig(storage.Base32EncodingName,
		path.Join(dataDir, "hoard"))
}
