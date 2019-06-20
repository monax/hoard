package stores

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

type fileSystemStore struct {
	rootDirectory string
	encoding      AddressEncoding
}

func NewFileSystemStore(rootDirectory string, encoding AddressEncoding) (*fileSystemStore, error) {
	err := os.MkdirAll(rootDirectory, 0700)
	if err != nil {
		return nil, err
	}
	return &fileSystemStore{
		rootDirectory: rootDirectory,
		encoding:      encoding,
	}, nil
}

func (inv *fileSystemStore) Put(address []byte, data []byte) ([]byte, error) {
	return address, ioutil.WriteFile(inv.Path(address), data, 0644)
}

func (inv *fileSystemStore) Get(address []byte) ([]byte, error) {
	return ioutil.ReadFile(inv.Path(address))
}

func (inv *fileSystemStore) Stat(address []byte) (*StatInfo, error) {
	fileInfo, err := os.Stat(inv.Path(address))
	statInfo := new(StatInfo)
	// Any kind of error means we should set exists false
	statInfo.Exists = err == nil
	if statInfo.Exists {
		statInfo.Size_ = uint64(fileInfo.Size())
	}
	// Don't treat not existing as an error
	if os.IsNotExist(err) {
		return statInfo, nil
	}
	return statInfo, err
}

func (inv *fileSystemStore) Location(address []byte) string {
	filePath := inv.Path(address)
	uri, err := url.Parse(filePath)
	if err != nil {
		return filePath
	}
	return uri.String()
}

func (inv *fileSystemStore) Path(address []byte) string {
	return path.Join(inv.rootDirectory,
		inv.encoding.EncodeToString(address))
}

func (inv *fileSystemStore) Name() string {
	return fmt.Sprintf("fileSystemStore[root=%s]", inv.rootDirectory)
}
