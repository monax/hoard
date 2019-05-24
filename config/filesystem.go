package config

import (
	"fmt"
	"path"

	"github.com/monax/hoard/v4/stores"

	"github.com/cep21/xdgbasedir"
)

type FileSystemConfig struct {
	RootDirectory string
}

func NewFileSystemConfig(addressEncoding, rootDirectory string) *Storage {
	return &Storage{
		StorageType:     Filesystem,
		AddressEncoding: addressEncoding,
		FileSystemConfig: &FileSystemConfig{
			RootDirectory: rootDirectory,
		},
	}
}

func DefaultFileSystemConfig() *Storage {
	dataDir, err := xdgbasedir.DataHomeDirectory()
	if err != nil {
		panic(fmt.Errorf("could not get XDG data dir: %s", err))
	}
	// Avoid '/' character for filesystem storage
	return NewFileSystemConfig(stores.Base32EncodingName,
		path.Join(dataDir, "hoard"))
}
