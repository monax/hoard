package storage

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

type fileSystemStore struct {
	rootDirectory   string
	addressEncoding AddressEncoding
}

func NewFileSystemStore(rootDirectory string, addressEncoding AddressEncoding) (*fileSystemStore, error) {
	err := os.MkdirAll(rootDirectory, 0700)
	if err != nil {
		return nil, err
	}
	return &fileSystemStore{
		rootDirectory:   rootDirectory,
		addressEncoding: addressEncoding,
	}, nil
}

func (fss *fileSystemStore) Put(address []byte, data []byte) ([]byte, error) {
	return address, ioutil.WriteFile(fss.Path(address), data, 0644)
}

func (fss *fileSystemStore) Get(address []byte) ([]byte, error) {
	return ioutil.ReadFile(fss.Path(address))
}

func (fss *fileSystemStore) Stat(address []byte) (*StatInfo, error) {
	fileInfo, err := os.Stat(fss.Path(address))
	statInfo := new(StatInfo)
	// Any kind of error means we should set exists false
	statInfo.Exists = err == nil
	if statInfo.Exists {
		statInfo.Size = uint64(fileInfo.Size())
	}
	// Don't treat not existing as an error
	if os.IsNotExist(err) {
		return statInfo, nil
	}
	return statInfo, err
}

func (fss *fileSystemStore) Location(address []byte) string {
	filePath := fss.Path(address)
	uri, err := url.Parse(filePath)
	if err != nil {
		return filePath
	}
	return uri.String()
}

func (fss *fileSystemStore) Path(address []byte) string {
	return path.Join(fss.rootDirectory,
		fss.addressEncoding.EncodeToString(address))
}

func (fss *fileSystemStore) Name() string {
	return fmt.Sprintf("fileSystemStore[root=%s]", fss.rootDirectory)
}
