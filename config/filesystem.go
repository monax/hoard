package config

import (
	"fmt"
	"path"

	"github.com/monax/hoard/v7/stores"

	"github.com/cep21/xdgbasedir"
)

type FileSystemConfig struct {
	RootDirectory string
}

func NewFileSystemConfig(addressEncoding, rootDirectory string) *Storage {
	conf := NewDefaultStorage()
	conf.StorageType = Filesystem
	conf.AddressEncoding = addressEncoding
	conf.FileSystemConfig = &FileSystemConfig{
		RootDirectory: rootDirectory,
	}
	return conf
}

func NewDefaultFileSystemConfig() *Storage {
	dataDir, err := xdgbasedir.DataHomeDirectory()
	if err != nil {
		panic(fmt.Errorf("could not get XDG data dir: %s", err))
	}
	// Avoid '/' character for filesystem storage
	return NewFileSystemConfig(stores.Base32EncodingName,
		path.Join(dataDir, "hoard"))
}
